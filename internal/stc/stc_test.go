package stc

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCreateFmeStc(t *testing.T) {
	// Create a request body for testing
	requestBody := []byte(`{
		"name": "John Doe",
		"address": "44, Sanusi street Somolu",
		"state": "Lagos",
		"email": "john@example.com",
		"phoneNumber": "07088547089"
	}`)

	// Create a new HTTP request with the request body
	req, err := http.NewRequest("POST", "/create-fme-stc", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Create a mock Gin context
	ctx, _ := gin.CreateTestContext(rr)
	ctx.Request = req

	// Call the handler function
	CreateFmeStc(ctx)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	expected := `{"message":"Student created successfully"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}