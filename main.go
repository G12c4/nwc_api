package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/untreu2/go-nwc"
)

// makePayment handles Lightning payments between any two wallets
// sender and recipient are keys in the walletURIs map
// amount is in millisatoshis
func makePayment(walletURIs map[string]string, sender string, recipient string, amount int) error {
	// Get URIs for both wallets
	senderURI, ok := walletURIs[sender]
	if !ok {
		return fmt.Errorf("sender wallet '%s' not found", sender)
	}
	
	recipientURI, ok := walletURIs[recipient]
	if !ok {
		return fmt.Errorf("recipient wallet '%s' not found", recipient)
	}
	
	// Initialize wallet clients
	senderClient, err := nwc.NewClient(senderURI)
	if err != nil {
		return fmt.Errorf("failed to initialize sender wallet: %w", err)
	}
	
	recipientClient, err := nwc.NewClient(recipientURI)
	if err != nil {
		return fmt.Errorf("failed to initialize recipient wallet: %w", err)
	}
	
	// Check sender balance
	balance, err := senderClient.GetBalance()
	if err != nil {
		return fmt.Errorf("failed to get sender balance: %w", err)
	}
	
	log.Printf("%s balance: %d msat", sender, balance.Balance)
	
	if balance.Balance < int64(amount) {
		return fmt.Errorf("insufficient funds in sender wallet: %d msat needed, %d msat available", amount, balance.Balance)
	}
	
	// Create invoice from recipient
	invoice, err := recipientClient.MakeInvoice(amount, fmt.Sprintf("Payment from %s to %s", sender, recipient))
	if err != nil {
		return fmt.Errorf("failed to create invoice: %w", err)
	}
	
	log.Printf("Created invoice for %d msat", amount)
	
	// Pay invoice with sender
	result, err := senderClient.PayInvoice(invoice)
	if err != nil {
		return fmt.Errorf("payment failed: %w", err)
	}
	
	log.Printf("Payment successful! Fees: %d msat", result.FeesPaid)
	
	// Check updated balances
	senderBalance, _ := senderClient.GetBalance()
	recipientBalance, _ := recipientClient.GetBalance()
	
	log.Printf("Updated %s balance: %d msat", sender, senderBalance.Balance)
	log.Printf("Updated %s balance: %d msat", recipient, recipientBalance.Balance)
	
	return nil
}

// Wallet URI functions moved to api.go

// euroToMsats converts Euro amount to millisatoshis using current exchange rate
// Returns the equivalent amount in millisatoshis
func euroToMsats(euroAmount float64) (int, error) {
	// CoinGecko API endpoint for BTC price in EUR
	url := "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=eur"
	
	// Create a client with timeout
	client := &http.Client{Timeout: 10 * time.Second}
	
	// Make the request
	resp, err := client.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch exchange rate: %w", err)
	}
	defer resp.Body.Close()
	
	// Check status code
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API returned non-200 status: %d", resp.StatusCode)
	}
	
	// Parse response
	var result map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to parse API response: %w", err)
	}
	
	// Extract BTC price in EUR
	btcPriceInEur, ok := result["bitcoin"]["eur"]
	if !ok {
		return 0, fmt.Errorf("could not find BTC/EUR exchange rate in response")
	}
	
	// Calculate conversions
	// 1 BTC = 100,000,000 satoshis
	// 1 satoshi = 1,000 millisatoshis
	btcAmount := euroAmount / btcPriceInEur
	satoshis := btcAmount * 100000000
	millisatoshis := satoshis * 1000000
	
	// Return as integer (rounded)
	return int(millisatoshis), nil
}

func main() {
	// Initialize our API
	router, err := InitializeAPI()
	if err != nil {
		log.Fatalf("Failed to initialize API: %v", err)
	}

	// Start the server
	port := ":8080"
	log.Printf("Starting server on%s", port)
	log.Printf("NWC Payment API endpoint: http://localhost%s/nwc_payment", port)
	log.Printf("Swagger documentation available at http://localhost%s/swagger/index.html", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}