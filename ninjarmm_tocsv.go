/*
ninjarmm_tocsv.go

Copyright (c) 2022 Mike Tremaine

MIT License

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

Description:
    This tool connects to the NinjaRMM API, retrieves a list of customers and their devices,
    and outputs a CSV summary of device counts by type for each customer.

    - Authenticates using HMAC-SHA1 with provided API keys.
    - Fetches organizations (customers) and their devices.
    - Aggregates device types per customer.
    - Outputs results as CSV to stdout.

Usage:
    go run ninjarmm_tocsv.go > output.csv
    
Author:
    Mike Tremaine <mgt@stellarcore.net>
*/

package main

import (
    "encoding/base64"
    "encoding/json"
    "crypto/hmac"
    "crypto/sha1"
    "fmt"
    "os"
    "time"
    "net/http"
)

//Globals
var ACCESS_KEY_ID string
var SECRET_ACCESS_KEY string

func init() {
    ACCESS_KEY_ID = os.Getenv("NINJA_ACCESS_KEY_ID")
    SECRET_ACCESS_KEY = os.Getenv("NINJA_SECRET_ACCESS_KEY")
    if ACCESS_KEY_ID == "" {
        ACCESS_KEY_ID = "<YOUR_ACCESS_KEY_ID>"
    }
    if SECRET_ACCESS_KEY == "" {
        SECRET_ACCESS_KEY = "<YOUR_SECRET_ACCESS_KEY>"
    }
}

//  URL endpoints
var API_HOST = "https://api.ninjarmm.com"

type Customer struct {
    Name            string `json:"name"`
    Description     string `json:"description"`
    Id              int  `json:"id"`
}

type Device struct {
    Id                      int  `json:"id"`
    OrganizationId          int `json:"organizationId"`
    LocationId              int `json:"locationId"`
    NodeClass               string `json:"nodeClass"`
    NodeRoleId              int `json:"nodeRoleId"`
    RolePolicyId            int `json:"rolePolicyId"`
    ApprovalStatus          string `json:"approvalStatus"`
    Offline                 bool `json:"offine"`
    SystemName              string `json:"systemName"`
    DnsName                 string `json:"dnsName"`
    Created                 float64 `json:"created"`
    LastContact             float64 `json:"lastContact"`
    LastUpdate              float64 `json:"lastUpdate"`
}

// getSignature generates the HMAC-SHA1 signature for API authentication.
func getSignature(secretAccessKey string, stringToSign string) string {
    data := base64.StdEncoding.EncodeToString([]byte(stringToSign))
    _hash := hmac.New(sha1.New, []byte(secretAccessKey))
    _hash.Write([]byte(data))
    sha := base64.StdEncoding.EncodeToString(_hash.Sum(nil))
    return sha
}

// getStringToSign builds the string to sign for the API request.
func getStringToSign(httpMethod string, contentMD5 string, contentType string, requestDateTime string, canonicalPath string) string {
    stringToSign := fmt.Sprintf("%s\n%s\n%s\n%s\n%s", httpMethod, contentMD5, contentType, requestDateTime, canonicalPath);
    return stringToSign;
}

// callApi performs a GET request to the NinjaRMM API and decodes the JSON response into dst.
func callApi(url string, signature string, requestDateTime string, dst interface{}) int {
    client := http.Client{}
    req , err := http.NewRequest("GET", API_HOST + url , nil)
    if err != nil {
        //Handle Error
    }

    req.Header = http.Header{
        "Authorization": []string{"NJ " + ACCESS_KEY_ID + ":" + signature},
        "Date": []string{requestDateTime},
    }
    res , err := client.Do(req)
    if err != nil {
        //Handle Error
    }

    defer res.Body.Close()
    dec := json.NewDecoder(res.Body)
    err = dec.Decode(dst)
    if err != nil {
        fmt.Printf("Error: %v", err)
        return 1
    }
    return 0
}

// main is the entry point for the tool.
func main() {

    var CUSTOMERS_URL string = "/v2/organizations"
    var requestDateTime string
    var stringToSign string
    var signature string
    var output_csv string
    
    //CSV Header
    output_csv = "Customer, DeviceType, Qty\n"

    requestDateTime = time.Now().UTC().Format(time.RFC1123Z)
    // Create signature for request
    stringToSign = getStringToSign("GET", "", "", requestDateTime, CUSTOMERS_URL)
    signature = getSignature(SECRET_ACCESS_KEY, stringToSign)

    var customer []Customer
    //parse the JSON response and display it 
    if ecode := callApi(CUSTOMERS_URL, signature, requestDateTime, &customer); ecode != 0 {
        fmt.Printf("Error: in callApi")
    }
    //We have customers
    var DEVICES_URL string
    //Customer Loop
    for _, values := range customer {
        DEVICES_URL = fmt.Sprintf("/v2/organization/%v/devices", values.Id)
        counters := make(map[string]int)
        requestDateTime = time.Now().UTC().Format(time.RFC1123Z)
        // Create signature for request
        stringToSign = getStringToSign("GET", "", "", requestDateTime, DEVICES_URL)
        signature = getSignature(SECRET_ACCESS_KEY, stringToSign)
        var device []Device
        //parse the JSON response and display it 
        if ecode := callApi(DEVICES_URL, signature, requestDateTime, &device); ecode != 0 {
            fmt.Printf("Error: in callApi")
        }

        for _, custdevice := range device {
            _type := custdevice.NodeClass
            //Munge MAC devices to Windows Workstation - client requested this.
            //Remove this if you want to keep MAC devices separate.
            if _type == "MAC" {
                _type = "WINDOWS_WORKSTATION"
            }
            counters[_type]++;
        }
        for _type, count := range counters {
            output_csv += fmt.Sprintf("%v, %v, %v\n", values.Name, _type, count)
        }
    }

    //Emit csv
    fmt.Printf(output_csv)
}
