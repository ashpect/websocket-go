package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var secretKey = []byte("secret")

// TokenClaims represents the structure of your custom claims (without role)
type TokenClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

// validateToken parses and validates the JWT token
func validateToken(tokenString string) (*jwt.Token, error) {
	parsedToken, err := jwt.Parse(tokenString, getSecretKey)
	if err != nil {
		return nil, err
	}
	return parsedToken, nil
}

// getSecretKey checks if the signing method is HMAC and returns the secret key
func getSecretKey(token *jwt.Token) (interface{}, error) {
	// Check if the signing method is HMAC
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
	}
	return secretKey, nil
}

// createJWTToken creates a new JWT token for a user
func createJWTToken(userID string) (string, error) {
	// Set expiration time for the token (e.g., 1 hour)
	expirationTime := time.Now().Add(1 * time.Hour)

	// Create the claims, which includes standard claims and custom claims (user_id)
	claims := &TokenClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), // Set expiration time
			IssuedAt:  time.Now().Unix(),     // Set issued time
			Issuer:    "websocket",           // Define your app's name or ID
		},
	}

	// Create a new token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return signedToken, nil
}
