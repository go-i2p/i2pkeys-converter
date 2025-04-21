// Package i2pkeys provides utilities for working with I2P key formats
package i2pkeys

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// I2P uses a custom Base64 encoding with '-' and '~' instead of '+' and '/'
var i2pB64Encoding = base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-~")

// KeyPair represents an I2P key pair with both public and private components
type KeyPair struct {
	PublicKey  []byte // The destination (public key)
	PrivateKey []byte // The private key
	FullData   []byte // The complete key data
}

// ConvertKeyFile converts an I2P binary key file to the two-line format required by Go I2P
func ConvertKeyFile(inputPath, outputPath string) error {
	// Read the key file as binary data
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read key file: %w", err)
	}

	// Check if input is already in the expected format
	if IsCorrectFormat(string(data)) {
		// Create output directory if needed
		outputDir := filepath.Dir(outputPath)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		// Just copy the file as is
		if err := os.WriteFile(outputPath, data, 0600); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}

		return nil
	}

	// Try to extract public key information if it's in I2P Base64 format
	keyData := string(data)
	var formattedOutput string

	// If data is in I2P Base64 format, try to extract the public key portion
	if isI2PBase64Format(keyData) {
		// Split by newlines in case there are multiple keys
		lines := strings.Split(keyData, "\n")
		completeKey := lines[0]

		// For I2P tunnel keys, the public key is the first 516 characters
		// This is a heuristic based on the standard format of I2P keys
		if len(completeKey) >= 516 {
			publicPart := completeKey[:516]
			formattedOutput = publicPart + "\n" + completeKey
		} else {
			// If we can't extract, convert the entire binary file
			completeKey = toI2PBase64(data)

			// Public key is typically the first 516 characters
			if len(completeKey) >= 516 {
				publicPart := completeKey[:516]
				formattedOutput = publicPart + "\n" + completeKey
			} else {
				return errors.New("key data too short to extract public key portion")
			}
		}
	} else {
		// Not in Base64 format, treat as binary and convert
		completeKey := toI2PBase64(data)

		// Public key is typically the first 516 characters
		if len(completeKey) >= 516 {
			publicPart := completeKey[:516]
			formattedOutput = publicPart + "\n" + completeKey
		} else {
			return errors.New("key data too short to extract public key portion")
		}
	}

	// Create output directory if needed
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write formatted output to file
	if err := os.WriteFile(outputPath, []byte(formattedOutput), 0600); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

// IsCorrectFormat checks if the data is already in the correct two-line format
func IsCorrectFormat(data string) bool {
	lines := strings.Split(strings.TrimSpace(data), "\n")
	if len(lines) != 2 {
		return false
	}

	// Check if both lines appear to be valid I2P Base64
	return isI2PBase64Format(lines[0]) && isI2PBase64Format(lines[1])
}

// isI2PBase64Format checks if a string appears to be in I2P Base64 format
func isI2PBase64Format(data string) bool {
	// Remove whitespace
	data = strings.TrimSpace(data)
	if data == "" {
		return false
	}

	// Check for I2P Base64 character set
	for _, r := range data {
		if !((r >= 'A' && r <= 'Z') ||
			(r >= 'a' && r <= 'z') ||
			(r >= '0' && r <= '9') ||
			r == '-' || r == '~' || r == '=') {
			return false
		}
	}

	// Try to decode
	_, err := fromI2PBase64(data)
	return err == nil
}

// toI2PBase64 converts binary data to I2P's Base64 variant
func toI2PBase64(data []byte) string {
	return i2pB64Encoding.EncodeToString(data)
}

// fromI2PBase64 converts I2P Base64 format back to binary
func fromI2PBase64(i2pBase64 string) ([]byte, error) {
	return i2pB64Encoding.DecodeString(i2pBase64)
}

// FormatKeysFile formats an existing I2P Base64 key into the proper two-line format
func FormatKeysFile(inputPath, outputPath string) error {
	// Read the key file
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read key file: %w", err)
	}

	// Check if it's already in the correct format
	if IsCorrectFormat(string(data)) {
		// Already in the correct format, just copy
		if inputPath != outputPath {
			if err := os.WriteFile(outputPath, data, 0600); err != nil {
				return fmt.Errorf("failed to write output file: %w", err)
			}
		}
		return nil
	}

	// Clean the input
	cleanedInput := cleanI2PBase64(string(data))

	// Split by lines (there might be multiple keys)
	lines := strings.Split(cleanedInput, "\n")

	// Process the first non-empty line
	var completeKey string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			completeKey = line
			break
		}
	}

	// Ensure we have enough data
	if len(completeKey) < 516 {
		return errors.New("key data too short to format correctly")
	}

	// Extract public key (first 516 characters)
	publicPart := completeKey[:516]

	// Create the proper two-line format
	formattedOutput := publicPart + "\n" + completeKey

	// Create output directory if needed
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write to output file
	if err := os.WriteFile(outputPath, []byte(formattedOutput), 0600); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

// cleanI2PBase64 cleans a string to ensure it only contains valid I2P Base64 characters
func cleanI2PBase64(data string) string {
	// Remove whitespace
	data = strings.TrimSpace(data)

	// Clean the line of any invalid characters
	var cleaned strings.Builder
	for _, r := range data {
		if (r >= 'A' && r <= 'Z') ||
			(r >= 'a' && r <= 'z') ||
			(r >= '0' && r <= '9') ||
			r == '-' || r == '~' || r == '=' ||
			r == '\n' {
			cleaned.WriteRune(r)
		}
	}

	return cleaned.String()
}
