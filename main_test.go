package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	// Create an in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect test database")
	}

	// Run migrations
	db.AutoMigrate(&Car{})

	return db
}

func TestCreateCar(t *testing.T) {
	// Setup test database
	db := setupTestDB()
	defer db.Close()

	// Setup Gin router
	router := gin.Default()
	router.POST("/cars", createCarHandler(db))

	// Mock HTTP request
	reqBody := []byte(`{"Brand": "Toyota", "Model": "Camry", "Status": "In Garage"}`)
	req, _ := http.NewRequest("POST", "/cars", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)

	// Validate the response body or any other expectations
	// You may want to parse the response JSON and check specific fields
	var response Car
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Toyota", response.Brand)
	// Add more assertions as needed
}

func TestGetCarList(t *testing.T) {
	// Setup test database
	db := setupTestDB()
	defer db.Close()

	// Insert a sample car record for testing
	db.Create(&Car{Brand: "Honda", Model: "Accord", Status: "In Garage"})

	// Setup Gin router
	router := gin.Default()
	router.GET("/cars", getCarListHandler(db))

	// Mock HTTP request
	req, _ := http.NewRequest("GET", "/cars", nil)
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	// Validate the response body or any other expectations
	// You may want to parse the response JSON and check specific fields
	var response []Car
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 1) // Expecting one car record
	// Add more assertions as needed
}

// Add similar test functions for other endpoints (updateCar, deleteCar, etc.)
