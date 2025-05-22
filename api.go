package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "nwc_app/docs"
	"nwc_app/middleware"
	"nwc_app/wallet"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/untreu2/go-nwc"
)

// @title           NWC Lightning Wallet API
// @version         1.0
// @description     API for making Lightning Network payments between wallets
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key

var walletURIs map[string]string

// NwcPaymentRequest represents the data needed to make an NWC payment
type NwcPaymentRequest struct {
	Sender     string  `json:"sender" binding:"required" example:"WALLET_JOSIP"`
	Recipient  string  `json:"recipient" binding:"required" example:"WALLET_VRATA_KRKE"`
	EuroAmount float64 `json:"euro_amount" binding:"required" example:"0.000001"`
}

// NwcPaymentResponse is the structure returned after making a payment
type NwcPaymentResponse struct {
	Success          bool    `json:"success"`
	Message          string  `json:"message"`
	EuroAmount       float64 `json:"euro_amount"`
	AmountMsats      int     `json:"amount_msats"`
	SenderBalance    int64   `json:"sender_balance"`
	RecipientBalance int64   `json:"recipient_balance"`
	FeesPaid         int64   `json:"fees_paid,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// ConversionResponse represents a currency conversion response
type ConversionResponse struct {
	EuroAmount float64 `json:"euro_amount"`
	MsatAmount int     `json:"msat_amount"`
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Status string            `json:"status"`
	Wallets map[string]bool  `json:"wallets"`
}



// @Summary      Make an NWC payment
// @Description  Transfer funds from one wallet to another using EUR amount
// @Tags         payments
// @Accept       json
// @Produce      json
// @Param        api_key   query   string  true  "API Key for authentication"
// @Param        payment   body    NwcPaymentRequest  true  "Payment Information"
// @Success      200      {object}  NwcPaymentResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      401      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /nwc_payment [post]
func nwcPaymentHandler(c *gin.Context) {
	// Validate API key
	requestAPIKey := c.Query("api_key")
	if requestAPIKey == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "API key is required. Please provide it in the api_key query parameter",
		})
		return
	}

	// Get the API key from environment variable
	apiKey := os.Getenv("NWC_API_KEY")

	// Check if API key is valid
	if requestAPIKey != apiKey {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Invalid API key",
		})
		return
	}
	
	// Process the request
	var req NwcPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: fmt.Sprintf("invalid request: %v", err),
		})
		return
	}

	// Convert Euro to msats
	msatAmount, err := euroToMsats(req.EuroAmount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: fmt.Sprintf("failed to convert EUR to msats: %v", err),
		})
		return
	}
	log.Printf("Converted %f EUR to %d msats", req.EuroAmount, msatAmount)

	if msatAmount <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "converted amount must be greater than 0",
		})
		return
	}

	// Make the payment
	err = makePayment(walletURIs, req.Sender, req.Recipient, msatAmount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Get updated balances
	senderClient, _ := nwc.NewClient(walletURIs[req.Sender])
	recipientClient, _ := nwc.NewClient(walletURIs[req.Recipient])
	
	senderBalance, _ := senderClient.GetBalance()
	recipientBalance, _ := recipientClient.GetBalance()

	c.JSON(http.StatusOK, NwcPaymentResponse{
		Success:          true,
		Message:          fmt.Sprintf("Successfully transferred %d msats (%.8f EUR) from %s to %s", msatAmount, req.EuroAmount, req.Sender, req.Recipient),
		EuroAmount:       req.EuroAmount,
		AmountMsats:      msatAmount,
		SenderBalance:    senderBalance.Balance,
		RecipientBalance: recipientBalance.Balance,
	})
}

// @Summary      Convert EUR to millisatoshis
// @Description  Converts a Euro amount to millisatoshis using current exchange rate
// @Tags         conversion
// @Produce      json
// @Param        amount    query  number  true  "Amount in EUR"
// @Param        api_key   query  string  true  "API Key for authentication"
// @Success      200  {object}  ConversionResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /convert/eur-to-msats [get]
func euroToMsatsHandler(c *gin.Context) {
	// Validate API key
	requestAPIKey := c.Query("api_key")
	if requestAPIKey == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "API key is required. Please provide it in the api_key query parameter",
		})
		return
	}

	// Get the API key from environment variable
	apiKey, err := wallet.LoadAPIKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to load API key",
		})
		return
	}

	// Check if API key is valid
	if requestAPIKey != apiKey {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Invalid API key",
		})
		return
	}

	// Process the request
	amountStr := c.Query("amount")
	if amountStr == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "amount parameter is required",
		})
		return
	}

	euroAmount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid amount format",
		})
		return
	}

	msatAmount, err := euroToMsats(euroAmount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: fmt.Sprintf("conversion failed: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, ConversionResponse{
		EuroAmount: euroAmount,
		MsatAmount: msatAmount,
	})
}

// @Summary      Check health of wallet
// @Description  Verifies connectivity to a specified wallet or all wallets if none specified
// @Tags         health
// @Produce      json
// @Param        wallet_id   query   string  false  "Wallet ID to check. If not provided, checks all wallets."
// @Success      200  {object}  HealthResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /health [get]
func healthCheckHandler(c *gin.Context) {
	// Create a map to store wallet statuses
	walletStatus := make(map[string]bool)
	
	// Check if a specific wallet ID was provided
	walletID := strings.ToUpper(c.Query("wallet_id"))
	if walletID != "" {
		// Check only the specified wallet
		walletURI, exists := walletURIs[walletID]
		if !exists || walletID == "" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: fmt.Sprintf("Wallet with ID '%s' not found", walletID),
			})
			return
		}
		
		// Initialize wallet client
		walletClient, err := nwc.NewClient(walletURI)
		if err != nil {
			walletStatus[walletID] = false
		} else {
			// Try to get balance to verify connection
			_, err = walletClient.GetBalance()
			walletStatus[walletID] = (err == nil)
		}
	}
	
	// Determine overall status
	overallStatus := "healthy"
	for _, status := range walletStatus {
		if !status {
			overallStatus = "degraded"
			break
		}
	}
	
	c.JSON(http.StatusOK, HealthResponse{
		Status: overallStatus,
		Wallets: walletStatus,
	})
}

// InitializeAPI sets up the Gin router with all routes and middleware
func InitializeAPI() (*gin.Engine, error) {
	// Load wallet URIs using our wallet package
	var err error
	walletURIs, err = wallet.LoadWalletURIs()
	if err != nil {
		return nil, fmt.Errorf("failed to load wallet URIs: %w", err)
	}

	// Set Gin to release mode in production
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create a new Gin router
	router := gin.New() // Use New() instead of Default() to avoid using the default logger

	// Add essential middleware
	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.CORSMiddleware()) // Add CORS support
	
	// All routes in a single group
	routes := router.Group("/")
	{
		// Health check endpoint - publicly accessible
		routes.GET("/health", healthCheckHandler)

		// Swagger UI endpoint
		routes.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, 
			ginSwagger.DefaultModelsExpandDepth(-1),
			ginSwagger.DocExpansion("list"),
			ginSwagger.PersistAuthorization(true)))
		
		// Payment endpoint - authentication handled in handler
		routes.POST("/nwc_payment", nwcPaymentHandler)
		
		// EUR to msat conversion endpoint - authentication handled in handler
		routes.GET("/convert/eur-to-msats", euroToMsatsHandler)
	}

	// Print registered routes for debugging
	log.Println("Registered routes:")
	for _, route := range router.Routes() {
		log.Printf("%s %s", route.Method, route.Path)
	}

	// Log API security information
	// Get the API key from environment variable
	apiKey, _ := wallet.LoadAPIKey()
	if apiKey == "" {
		log.Println("WARNING: Using default API key. Set NWC_API_KEY environment variable for production use.")
	} else {
		log.Println("API Key authentication is enabled.")
	}

	return router, nil
}
