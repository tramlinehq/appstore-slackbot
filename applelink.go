package main

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"net/http"
	"time"
)

type AppleCredentials struct {
	BundleID string
	IssuerID string
	KeyID    string
	P8File   []byte
}

type ApplelinkCredentials struct {
	Aud    string
	Issuer string
	Secret string
}

type AppMetadata struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	BundleId string `json:"bundle_id"`
	Sku      string `json:"sku"`
}

func GetAppMetadata(applelinkCreds ApplelinkCredentials, credentials AppleCredentials) (AppMetadata, error) {
	appMetadata := AppMetadata{}

	requestURL := fmt.Sprintf("%s/apple/connect/v1/apps/%s", applelinkHost, credentials.BundleID)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		fmt.Printf("applelink: could not create request: %s\n", err)
		return appMetadata, err
	}

	storeToken, err := getStoreToken(credentials)
	if err != nil {
		fmt.Printf("applelink: could not create store token: %s\n", err)
		return appMetadata, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", getAuthToken(applelinkCreds)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-AppStoreConnect-Key-Id", credentials.KeyID)
	req.Header.Set("X-AppStoreConnect-Issuer-Id", credentials.IssuerID)
	req.Header.Set("X-AppStoreConnect-Token", storeToken)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Printf("applelink: request failed: %s\n", err)
		return appMetadata, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return appMetadata, fmt.Errorf("applelink: request failed with status - %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("applelink: could not read response body: %s\n", err)
		return appMetadata, err
	}
	err = json.Unmarshal(body, &appMetadata)
	if err != nil {
		fmt.Printf("applelink: could not parse reponse body: %s\n", err)
		return appMetadata, err
	}

	return appMetadata, err
}

func getAuthToken(credentials ApplelinkCredentials) string {
	expiry := time.Now().Add(10 * time.Minute)
	claims := &jwt.RegisteredClaims{
		Audience:  jwt.ClaimStrings{credentials.Aud},
		Issuer:    credentials.Issuer,
		ExpiresAt: jwt.NewNumericDate(expiry),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signingKey := []byte(credentials.Secret)
	signedToken, _ := token.SignedString(signingKey)
	fmt.Println("Signed token", signedToken)
	return signedToken
}

func getStoreToken(credentials AppleCredentials) (string, error) {
	expiry := time.Now().Add(10 * time.Minute)

	claims := &jwt.RegisteredClaims{
		Audience:  jwt.ClaimStrings{"appstoreconnect-v1"},
		Issuer:    credentials.IssuerID,
		ExpiresAt: jwt.NewNumericDate(expiry),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	t.Header["kid"] = credentials.KeyID
	privateKey, err := parsePrivateKey(credentials.P8File)
	if err != nil {
		return "", err
	}
	token, err := t.SignedString(privateKey)

	return token, err
}

// ErrMissingPEM happens when the bytes cannot be decoded as a PEM block.
var ErrMissingPEM = errors.New("no PEM blob found")

// ErrInvalidPrivateKey happens when a key cannot be parsed as a ECDSA PKCS8 private key.
var ErrInvalidPrivateKey = errors.New("key could not be parsed as a valid ecdsa.PrivateKey")

func parsePrivateKey(blob []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(blob)
	if block == nil {
		return nil, ErrMissingPEM
	}

	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	if key, ok := parsedKey.(*ecdsa.PrivateKey); ok {
		return key, nil
	}

	return nil, ErrInvalidPrivateKey
}
