package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/HousewareHQ/backend-engineering-octernship/api/server/models"
	"github.com/gin-gonic/gin"
)

func TestLogin(t *testing.T) {
	// create a new Gin router instance
	router := gin.Default()

	// define the login endpoint
	router.POST("api/v1/login", func(c *gin.Context) {
		// parse the request body
		var reqBody struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		// validate the username and password
		if reqBody.Username != "testuser" || reqBody.Password != "testpass" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		// generate a JWT token
		token := "jwt-token"

		// send the response
		c.JSON(http.StatusOK, gin.H{"message": "success", "token": token})
	})

	// create a request body for the login endpoint
	body := []byte(`{"username":"testuser","password":"testpass"}`)
	req, err := http.NewRequest("POST", "http://localhost:8080/api/v1/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// set request headers
	req.Header.Set("Content-Type", "application/json")

	// create a response recorder to record the response from the handler
	rr := httptest.NewRecorder()

	// call the handler function for the login endpoint
	router.ServeHTTP(rr, req)

	// check the status code of the response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// check the response body
	expected := `{"message":"success","token":"jwt-token"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestLogout(t *testing.T) {
	// Initialize a new Gin router
	router := gin.Default()

	// Add a GET route for the logout endpoint
	router.GET("/api/v1/logout", func(c *gin.Context) {
		// Implement your logout functionality here
		// ...
		// Return a success response
		c.JSON(http.StatusOK, gin.H{
			"message": "Logout successful",
		})
	})

	// Create a new HTTP request with the GET method and the logout endpoint
	req, err := http.NewRequest("GET", "http://localhost:8080/api/v1/logout", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a new HTTP response recorder
	res := httptest.NewRecorder()

	// Serve the HTTP request to the router using the response recorder
	router.ServeHTTP(res, req)

	// Check the status code of the response
	if status := res.Code; status != http.StatusOK {
		t.Errorf("Unexpected status code: got %v, expected %v", status, http.StatusOK)
	}

	// Check the response body of the response
	expectedBody := `{"message":"Logout successful"}`
	if res.Body.String() != expectedBody {
		t.Errorf("Unexpected response body: got %v, expected %v", res.Body.String(), expectedBody)
	}
}

func TestGetAllUsers(t *testing.T) {
	// Create a new test router
	router := gin.Default()

	// Define the endpoint handler
	router.GET("api/v1/users", func(c *gin.Context) {
		// Your code to retrieve all users from the database
		users := []models.User{{Username: "abcd", Password: "12345"}, {Username: "xyzw", Password: "56789"}}

		// Return the users in the response
		c.JSON(http.StatusOK, users)
	})

	// Create a new HTTP request to the endpoint
	req, err := http.NewRequest("GET", "http://localhost:8080/api/v1/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new HTTP response recorder
	rr := httptest.NewRecorder()

	// Send the HTTP request to the test router
	router.ServeHTTP(rr, req)

	// Check that the response status code is 200 OK
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, rr.Code)
	}

	// Decode the response body into a slice of User structs
	var users []models.User
	err = json.NewDecoder(rr.Body).Decode(&users)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the returned users match the expected users
	expectedUsers := []models.User{{Username: "abcd", Password: "12345"}, {Username: "xyzw", Password: "56789"}}
	if !reflect.DeepEqual(users, expectedUsers) {
		t.Errorf("Expected users %v but got %v", expectedUsers, users)
	}
}
