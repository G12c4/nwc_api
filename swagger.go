package main

import (
	"html/template"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Custom Swagger UI HTML template that includes input field for API key
const swaggerUITemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>NWC Lightning Payment API</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui.css" />
    <style>
        .api-key-input {
            display: flex;
            align-items: center;
            margin: 15px 0;
            padding: 10px;
            background-color: #f0f0f0;
            border-radius: 4px;
        }
        .api-key-input input {
            flex-grow: 1;
            padding: 8px;
            border: 1px solid #ddd;
            border-radius: 4px;
            margin-left: 10px;
        }
        .api-key-input button {
            margin-left: 10px;
            padding: 8px 15px;
            background-color: #4990e2;
            color: white;
            border: none;
            border-radius: 4px;
            cursor: pointer;
        }
        .api-key-input button:hover {
            background-color: #3a7bc8;
        }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui-bundle.js"></script>
    <script>
        window.onload = function() {
            // Get default API key from local storage or use empty string
            const defaultApiKey = localStorage.getItem('nwc_api_key') || '';
            
            // Create API key input UI
            const apiKeyInput = document.createElement('div');
            apiKeyInput.className = 'api-key-input';
            apiKeyInput.innerHTML = '<strong>API Key:</strong>' +
                '<input type="text" id="api-key-input" placeholder="Enter your API key" value="' + defaultApiKey + '">' +
                '<button onclick="setApiKey()">Apply</button>';
            
            // Build Swagger UI
            const ui = SwaggerUIBundle({
                url: "/swagger/doc.json",
                dom_id: "#swagger-ui",
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIBundle.SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "BaseLayout",
                requestInterceptor: (req) => {
                    const apiKey = document.getElementById('api-key-input').value;
                    if (apiKey) {
                        req.headers['X-API-Key'] = apiKey;
                    }
                    return req;
                }
            });

            // Insert API key input before Swagger UI
            const swaggerUI = document.getElementById('swagger-ui');
            swaggerUI.parentNode.insertBefore(apiKeyInput, swaggerUI);
            
            // Function to set API key
            window.setApiKey = function() {
                const apiKey = document.getElementById('api-key-input').value;
                localStorage.setItem('nwc_api_key', apiKey);
                // Reload to apply changes to all operations
                window.location.reload();
            }
        }
    </script>
</body>
</html>
`

// CustomSwaggerHandler returns a custom Swagger UI handler with API key input
func CustomSwaggerHandler() gin.HandlerFunc {
	defaultHandler := ginSwagger.WrapHandler(swaggerFiles.Handler)
	
	return func(c *gin.Context) {
		// Get the last part of the URL path
		path := c.Param("any")

		// If requesting the JSON spec, use the default handler
		if path == "/doc.json" || path == "doc.json" {
			defaultHandler(c)
			return
		}
		
		// For the UI, use our custom template
		tmpl, err := template.New("swagger").Parse(swaggerUITemplate)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error parsing Swagger UI template")
			return
		}
		
		c.Status(http.StatusOK)
		c.Header("Content-Type", "text/html; charset=utf-8")
		
		// Set default API key from environment (for development only)
		apiKey := os.Getenv("NWC_API_KEY")
		if apiKey == "" {
			apiKey = "default-dev-api-key-change-in-production"
		}
		
		// Execute template
		tmpl.Execute(c.Writer, nil)
	}
}
