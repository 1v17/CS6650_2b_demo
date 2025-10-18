# PowerShell script to test all API endpoints locally
# Make sure your Go application is running locally on port 8080

param(
    [string]$BaseUrl = "http://localhost:8080",
    [switch]$Verbose
)

# Colors for output
$Green = "Green"
$Red = "Red"
$Yellow = "Yellow"
$Cyan = "Cyan"

# Function to make HTTP requests and display results
function Test-Endpoint {
    param(
        [string]$Method,
        [string]$Url,
        [string]$Body = $null,
        [string]$Description
    )
    
    Write-Host "`n=== Testing: $Description ===" -ForegroundColor $Cyan
    Write-Host "Method: $Method" -ForegroundColor $Yellow
    Write-Host "URL: $Url" -ForegroundColor $Yellow
    
    if ($Body) {
        Write-Host "Body: $Body" -ForegroundColor $Yellow
    }
    
    try {
        $headers = @{
            "Content-Type" = "application/json"
        }
        
        $params = @{
            Uri = $Url
            Method = $Method
            Headers = $headers
        }
        
        if ($Body) {
            $params.Body = $Body
        }
        
        $response = Invoke-RestMethod @params -ErrorAction Stop
        $statusCode = 200 # Invoke-RestMethod only returns on success
        
        Write-Host "✅ SUCCESS - Status: $statusCode" -ForegroundColor $Green
        
        if ($Verbose) {
            Write-Host "Response:" -ForegroundColor $Yellow
            $response | ConvertTo-Json -Depth 10 | Write-Host
        } else {
            # Show abbreviated response
            if ($response -is [PSCustomObject]) {
                $response | ConvertTo-Json -Depth 2 -Compress | Write-Host
            } else {
                Write-Host $response
            }
        }
        
        return $response
    }
    catch {
        $statusCode = $_.Exception.Response.StatusCode.value__
        Write-Host "❌ FAILED - Status: $statusCode" -ForegroundColor $Red
        
        if ($_.Exception.Response) {
            $reader = [System.IO.StreamReader]::new($_.Exception.Response.GetResponseStream())
            $responseBody = $reader.ReadToEnd()
            Write-Host "Error Response: $responseBody" -ForegroundColor $Red
        } else {
            Write-Host "Error: $($_.Exception.Message)" -ForegroundColor $Red
        }
        
        return $null
    }
}

# Function to check if server is running
function Test-ServerHealth {
    Write-Host "Checking if server is running at $BaseUrl..." -ForegroundColor $Cyan
    
    try {
        $response = Invoke-RestMethod -Uri "$BaseUrl/health" -Method GET -TimeoutSec 5
        Write-Host "✅ Server is running!" -ForegroundColor $Green
        return $true
    }
    catch {
        Write-Host "❌ Server is not running or not accessible at $BaseUrl" -ForegroundColor $Red
        Write-Host "Please start your Go application first using: go run main.go" -ForegroundColor $Yellow
        return $false
    }
}

# Main testing function
function Start-APITests {
    Write-Host "🚀 Starting API Endpoint Tests" -ForegroundColor $Cyan
    Write-Host "Base URL: $BaseUrl" -ForegroundColor $Yellow
    
    # Check if server is running
    if (-not (Test-ServerHealth)) {
        return
    }
    
    # Test 1: Health Check
    Test-Endpoint -Method "GET" -Url "$BaseUrl/health" -Description "Health Check"
    
    # Test 2: Create a Product
    $productBody = @{
        product_id = 1
        sku = "TEST-SKU-001"
        manufacturer = "Test Manufacturer"
        category_id = 1
        weight = 100
        some_other_id = 123
    } | ConvertTo-Json
    
    $createdProduct = Test-Endpoint -Method "POST" -Url "$BaseUrl/products" -Body $productBody -Description "Create Product"
    
    # Test 3: Get Product by ID (valid ID)
    Test-Endpoint -Method "GET" -Url "$BaseUrl/products/1" -Description "Get Product by ID (Valid)"
    
    # Test 4: Get Product by ID (invalid ID - should return 404)
    Test-Endpoint -Method "GET" -Url "$BaseUrl/products/999" -Description "Get Product by ID (Not Found)"
    
    # Test 5: Get Product by ID (invalid format - should return 400)
    Test-Endpoint -Method "GET" -Url "$BaseUrl/products/invalid" -Description "Get Product by ID (Invalid Format)"
    
    # Test 6: Create a Shopping Cart
    $cartBody = @{
        customer_id = 1
    } | ConvertTo-Json
    
    $createdCart = Test-Endpoint -Method "POST" -Url "$BaseUrl/shopping-carts" -Body $cartBody -Description "Create Shopping Cart"
    
    # Extract cart ID for subsequent tests
    $cartId = if ($createdCart -and $createdCart.shopping_cart_id) { 
        $createdCart.shopping_cart_id 
    } else { 
        1  # Default fallback
    }
    
    # Test 7: Get Shopping Cart by ID (valid ID)
    Test-Endpoint -Method "GET" -Url "$BaseUrl/shopping-carts/$cartId" -Description "Get Shopping Cart by ID (Valid)"
    
    # Test 8: Add Item to Shopping Cart
    $addItemBody = @{
        product_id = 1
        quantity = 2
    } | ConvertTo-Json
    
    Test-Endpoint -Method "POST" -Url "$BaseUrl/shopping-carts/$cartId/items" -Body $addItemBody -Description "Add Item to Shopping Cart"
    
    # Test 9: Get Shopping Cart after adding items
    Test-Endpoint -Method "GET" -Url "$BaseUrl/shopping-carts/$cartId" -Description "Get Shopping Cart with Items"
    
    # Test 10: Get Shopping Cart by ID (invalid ID - should return 404)
    Test-Endpoint -Method "GET" -Url "$BaseUrl/shopping-carts/999" -Description "Get Shopping Cart by ID (Not Found)"
    
    # Test 11: Add Item to Non-existent Cart (should return 404)
    Test-Endpoint -Method "POST" -Url "$BaseUrl/shopping-carts/999/items" -Body $addItemBody -Description "Add Item to Non-existent Cart"
    
    # Test 12: Create Product with Invalid Data (should return 400)
    $invalidProductBody = @{
        sku = "INVALID-PRODUCT"
        # Missing required fields
    } | ConvertTo-Json
    
    Test-Endpoint -Method "POST" -Url "$BaseUrl/products" -Body $invalidProductBody -Description "Create Product with Invalid Data"
    
    # Test 13: Create Cart with Invalid Data (should return 400)
    $invalidCartBody = @{
        customer_id = -1  # Invalid customer ID
    } | ConvertTo-Json
    
    Test-Endpoint -Method "POST" -Url "$BaseUrl/shopping-carts" -Body $invalidCartBody -Description "Create Cart with Invalid Customer ID"
    
    # Test 14: Add Invalid Item to Cart (should return 400)
    $invalidItemBody = @{
        product_id = -1  # Invalid product ID
        quantity = 0     # Invalid quantity
    } | ConvertTo-Json
    
    Test-Endpoint -Method "POST" -Url "$BaseUrl/shopping-carts/$cartId/items" -Body $invalidItemBody -Description "Add Invalid Item to Cart"
    
    Write-Host "`n🎉 API Testing Complete!" -ForegroundColor $Green
    Write-Host "Check the results above for any failing tests." -ForegroundColor $Yellow
}

# Additional test scenarios
function Start-LoadTests {
    Write-Host "`n🔄 Running Load Tests (Creating multiple products and carts)..." -ForegroundColor $Cyan
    
    # Create multiple products
    for ($i = 10; $i -le 15; $i++) {
        $productBody = @{
            product_id = $i
            sku = "LOAD-TEST-SKU-$i"
            manufacturer = "Load Test Manufacturer"
            category_id = $i
            weight = $i * 10
            some_other_id = $i * 100
        } | ConvertTo-Json
        
        Test-Endpoint -Method "POST" -Url "$BaseUrl/products" -Body $productBody -Description "Load Test - Create Product $i"
    }
    
    # Create multiple carts
    for ($i = 10; $i -le 15; $i++) {
        $cartBody = @{
            customer_id = $i
        } | ConvertTo-Json
        
        Test-Endpoint -Method "POST" -Url "$BaseUrl/shopping-carts" -Body $cartBody -Description "Load Test - Create Cart for Customer $i"
    }
    
    Write-Host "✅ Load Tests Complete!" -ForegroundColor $Green
}

# Help function
function Show-Help {
    Write-Host @"
Usage: .\test-api-endpoints.ps1 [OPTIONS]

OPTIONS:
    -BaseUrl <url>    Base URL for the API (default: http://localhost:8080)
    -Verbose          Show detailed response bodies
    -Help             Show this help message

EXAMPLES:
    .\test-api-endpoints.ps1
    .\test-api-endpoints.ps1 -BaseUrl "http://localhost:3000" -Verbose
    .\test-api-endpoints.ps1 -Help

PREREQUISITES:
    1. Start your Go application: cd src && go run main.go
    2. Ensure the application is running on the specified port (default: 8080)

API ENDPOINTS TESTED:
    GET    /health                        - Health check
    GET    /products/:productId           - Get product by ID
    POST   /products                      - Create product
    POST   /shopping-carts               - Create shopping cart  
    GET    /shopping-carts/:id           - Get shopping cart by ID
    POST   /shopping-carts/:id/items     - Add item to shopping cart
"@
}

# Main execution
if ($args -contains "-Help" -or $args -contains "--help" -or $args -contains "-h") {
    Show-Help
    return
}

# Run the tests
Start-APITests

# Ask if user wants to run load tests
$runLoadTests = Read-Host "`nDo you want to run additional load tests? (y/N)"
if ($runLoadTests -eq "y" -or $runLoadTests -eq "Y") {
    Start-LoadTests
}

Write-Host "`n📋 Test Summary:" -ForegroundColor $Cyan
Write-Host "- All main API endpoints have been tested" -ForegroundColor $Yellow
Write-Host "- Both success and error scenarios were covered" -ForegroundColor $Yellow
Write-Host "- Check the output above for any failures that need attention" -ForegroundColor $Yellow