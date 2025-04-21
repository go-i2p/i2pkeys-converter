package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-i2p/i2pkeys-converter/i2pkeys"
)

func main() {
	// Command line arguments
	inputFile := flag.String("in", "", "Path to the I2P key file (required)")
	outputFile := flag.String("out", "", "Path to save the formatted key (optional)")
	verbose := flag.Bool("v", false, "Verbose output with key details")
	checkFormat := flag.Bool("check", false, "Check if a file is already in the correct format")

	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "I2P Keys Converter - Format I2P keys for Go I2P libraries\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s -in keyfile [-out outputfile] [-v] [-check]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  Convert binary key file:   %s -in keys.dat -out keys.dat.formatted\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Check key file format:     %s -in keys.dat -check\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Format with verbose info:  %s -in keys.dat -v\n", os.Args[0])
	}

	flag.Parse()

	// Validate input file parameter
	if *inputFile == "" {
		fmt.Println("Error: Input file (-in) is required")
		flag.Usage()
		os.Exit(1)
	}

	// Check if input file exists
	if _, err := os.Stat(*inputFile); os.IsNotExist(err) {
		fmt.Printf("Error: Input file '%s' does not exist\n", *inputFile)
		os.Exit(1)
	}

	// If check mode is enabled, just check the format
	if *checkFormat {
		data, err := os.ReadFile(*inputFile)
		if err != nil {
			fmt.Printf("Error reading file: %s\n", err)
			os.Exit(1)
		}

		if i2pkeys.IsCorrectFormat(string(data)) {
			fmt.Println("File IS in the correct two-line format")
			os.Exit(0)
		} else {
			fmt.Println("File is NOT in the correct two-line format")
			os.Exit(1)
		}
	}

	// Set default output file if not specified
	if *outputFile == "" {
		baseName := filepath.Base(*inputFile)
		dir := filepath.Dir(*inputFile)
		*outputFile = filepath.Join(dir, baseName+".formatted")
	}

	// Print operation info
	fmt.Printf("Formatting I2P key file: %s\n", *inputFile)
	fmt.Printf("Output file: %s\n", *outputFile)

	// Convert the key file
	err := i2pkeys.ConvertKeyFile(*inputFile, *outputFile)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	// Verify the result
	resultData, err := os.ReadFile(*outputFile)
	if err != nil {
		fmt.Printf("Error reading result file: %s\n", err)
		os.Exit(1)
	}

	if i2pkeys.IsCorrectFormat(string(resultData)) {
		fmt.Println("Conversion successful - key is now in the correct format")

		// Display additional information if verbose mode is enabled
		if *verbose {
			lines := strings.Split(string(resultData), "\n")
			if len(lines) >= 2 {
				publicKeyPreview := truncateString(lines[0], 40)
				fullKeyPreview := truncateString(lines[1], 40)

				fmt.Println("\nKey Information:")
				fmt.Printf("- Destination (public key): %s...\n", publicKeyPreview)
				fmt.Printf("- Full key length: %d characters\n", len(lines[1]))
				fmt.Printf("- Full key preview: %s...\n", fullKeyPreview)
				fmt.Println("\nFormat: Two lines")
				fmt.Println("- Line 1: Base64-encoded destination (public key)")
				fmt.Println("- Line 2: Base64-encoded full keypair (public + private)")
			}
		}
	} else {
		fmt.Println("Warning: Output file is not in the correct format")
		os.Exit(1)
	}
}

// truncateString truncates a string and adds ellipsis if needed
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
