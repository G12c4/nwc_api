// Package wallet provides loading wallet URIs from .env file
package wallet

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadWalletURIs loads wallet URIs from the .env file
// Returns a map with keys being wallet identifiers and values being wallet URIs
func LoadWalletURIs() (map[string]string, error) {
	// Create a map to store wallet URIs
	walletURIs := make(map[string]string)

	// Read the .env file
	env, err := godotenv.Read(".env")
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v. Using empty wallet map.", err)
	}
	
	// Copy all entries from the env map to the walletURIs map
	for k, v := range env {
		walletURIs[k] = v
	}
	
	return walletURIs, nil
}

func LoadAPIKey() (string, error) {
	apiKey := ""
	// Read the .env file
	env, err := godotenv.Read(".env")
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v. Using empty wallet map.", err)
	}
	
	// Copy all entries from the env map to the walletURIs map
	for _, v := range env {
		apiKey = v
	}
	
	return apiKey, nil
}
