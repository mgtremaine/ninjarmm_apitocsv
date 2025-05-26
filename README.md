# ninjarmm_apitocsv

A simple Go client for the NinjaOne API that exports device and client data to CSV for accounting or reporting purposes.

## Features

- Connects to the NinjaOne API using API credentials
- Retrieves device and client information
- Outputs data in CSV format for easy import into spreadsheets or accounting tools


## Usage

1. **Clone the repository:**
    ```sh
    git clone https://github.com/mgtremaine/ninjarmm_apitocsv.git
    cd ninjarmm_apitocsv
    ```

2. **Modify/Build the tool:**
    Decide how to pass the API KEY and SECRET. You can pass them as ENV variables or simply hardcode them in the script.
    Once you decide then build the executable.
    ```sh
    go build -o ninjarmm_apitocsv ninjarmm_tocsv.go
    ```

3. **Run the tool:**
    Skip the export if you have added your key to the code.
    ```sh
    export NINJA_API_KEY="your_api_key"
    export NINJA_API_SECRET="your_api_secret"
    ./ninjarmm_apitocsv > output.csv
    ```


## Requirements

- Go 1.18 or newer
- NinjaOne API access

## License

MIT License
