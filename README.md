# I2P Keys Converter

A utility for converting and formatting I2P key files for use with Go I2P libraries.

## Description

I2P Keys Converter handles the conversion between I2P binary key files and the two-line text format required by Go I2P libraries. It properly formats keys using I2P's custom Base64 encoding (with '-' and '~' instead of '+' and '/').

## Installation

```bash
# Clone the repository
git clone https://github.com/go-i2p/i2pkeys-converter.git
cd i2pkeys-converter

# Build the binary
go build

# Or build and install to ~go/bin
go install github.com/go-i2p/i2pkeys-converter
```

## Usage

```bash
# Convert binary key file to formatted two-line format
i2pkeys-converter -in keys.dat -out keys.dat.formatted

# Check if a file is already in the correct format
i2pkeys-converter -in keys.dat -check

# Format with verbose information about the key
i2pkeys-converter -in keys.dat -v
```

## Features

- Converts between binary I2P key formats and the two-line format
- Validates key format correctness
- Preserves the proper I2P Base64 encoding
- Handles the public/private key extraction and formatting
- Provides verbose output with key details

## License

This project is licensed under the MIT License - see the LICENSE file for details.