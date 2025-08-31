package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "mysecretkey"
	expiresIn := time.Minute * 15

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Error making JWT: %v", err)
	}

	returnedUserID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("Error validating JWT: %v", err)
	}

	if returnedUserID != userID {
		t.Fatalf("Expected userID %v, got %v", userID, returnedUserID)
	}

	// Test with invalid token
	_, err = ValidateJWT(token+"invalid", tokenSecret)
	if err == nil {
		t.Fatalf("Expected error validating invalid JWT, got none")
	}

	// Test with wrong secret
	_, err = ValidateJWT(token, "wrongsecret")
	if err == nil {
		t.Fatalf("Expected error validating JWT with wrong secret, got none")
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "mysecretkey"
	expiresIn := time.Second * 1

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Error making JWT: %v", err)
	}

	// Wait for the token to expire
	time.Sleep(2 * time.Second)

	_, err = ValidateJWT(token, tokenSecret)
	if err == nil {
		t.Fatalf("Expected error validating expired JWT, got none")
	}
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	tokenSecret := "mysecretkey"
	invalidToken := "this.is.not.a.valid.token"

	_, err := ValidateJWT(invalidToken, tokenSecret)

	if err == nil {
		t.Fatalf("Expected error validating invalid JWT, got none")
	}
}
