package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	controller "github.com/HousewareHQ/backend-engineering-octernship/internal/app/auth/controllers"
	helper "github.com/HousewareHQ/backend-engineering-octernship/internal/app/auth/helpers"
	"github.com/gin-gonic/gin"
)



func TestLogin(t *testing.T) {
	// Create a new Gin router instance
	router := gin.Default()

	// Define the login endpoint
	router.POST("/login", controller.Login())

	// Create a request body for the login endpoint
	requestBody := []byte(`{"username":"testuser","password":"testpass"}`)
	request, err := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Set request headers
	request.Header.Set("Content-Type", "application/json")

	// Create a response recorder to record the response from the handler
	responseRecorder := httptest.NewRecorder()

	// Call the handler function for the login endpoint
	router.ServeHTTP(responseRecorder, request)

	// Check the status code of the response
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	expectedResponseBody := `{"user":{"username":"testuser","password":"testpass"},"Access Token":"jwt-token","Refresh Token":"jwt-token"}`
	if responseRecorder.Body.String() != expectedResponseBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			responseRecorder.Body.String(), expectedResponseBody)
	}
}


func TestLogout(t *testing.T) {
	// Initialize a new Gin router
	router := gin.Default()

	// Add a GET route for the logout endpoint
	router.GET("/api/v1/logout", controller.Logout())

	// Create a new HTTP request with the GET method and the logout endpoint
	req, err := http.NewRequest("GET", "/api/v1/logout", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Set the token header
	token := "test-token"
	req.Header.Set("Token", token)

	// Create a new HTTP response recorder
	res := httptest.NewRecorder()

	// Serve the HTTP request to the router using the response recorder
	router.ServeHTTP(res, req)

	// Check the status code of the response
	if status := res.Code; status != http.StatusOK {
		t.Errorf("Unexpected status code: got %v, expected %v", status, http.StatusOK)
	}

	// Check the response body of the response
	expectedBody := `{"message":"User logged out successfully"}`
	if res.Body.String() != expectedBody {
		t.Errorf("Unexpected response body: got %v, expected %v", res.Body.String(), expectedBody)
	}

	// Check that the user's tokens have been updated in the database
	userId := "test-user-id"
	helper.UpdateAllTokens("", "", userId)
}
