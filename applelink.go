package main

import (
	"ciderbot/types"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func applelinkRequest(credentials *types.AppleCredentials, requestURL string, httpMethod string) ([]byte, error) {
	var response []byte
	req, err := http.NewRequest(httpMethod, requestURL, nil)
	if err != nil {
		fmt.Printf("applelink: could not create request: %s\n", err)
		return response, err
	}

	storeToken, err := getAppStoreToken(credentials)
	if err != nil {
		fmt.Printf("applelink: could not create store token: %s\n", err)
		return response, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", getApplelinkAuthToken(*applelinkCredentials)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-AppStoreConnect-Key-Id", credentials.KeyID)
	req.Header.Set("X-AppStoreConnect-Issuer-Id", credentials.IssuerID)
	req.Header.Set("X-AppStoreConnect-Token", storeToken)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Printf("applelink: request failed: %s\n", err)
		return response, err
	}

	if resp.StatusCode > 299 {
		return response, fmt.Errorf("applelink: request failed with status - %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("applelink: could not read response body: %s\n", err)
		return body, err
	}

	return body, err

}

func getAppMetadata(credentials *types.AppleCredentials) (types.AppMetadata, error) {
	appMetadata := types.AppMetadata{}

	requestURL := fmt.Sprintf("%s/apple/connect/v1/apps/%s", applelinkHost, credentials.BundleID)
	body, err := applelinkRequest(credentials, requestURL, http.MethodGet)

	err = json.Unmarshal(body, &appMetadata)
	if err != nil {
		fmt.Printf("applelink: could not parse reponse body: %s\n", err)
		return appMetadata, err
	}

	return appMetadata, err
}

func getAppCurrentStatus(credentials *types.AppleCredentials) ([]types.AppCurrentStatus, error) {
	var appCurrentStatuses []types.AppCurrentStatus

	requestURL := fmt.Sprintf("%s/apple/connect/v1/apps/%s/current_status", applelinkHost, credentials.BundleID)

	body, err := applelinkRequest(credentials, requestURL, http.MethodGet)
	err = json.Unmarshal(body, &appCurrentStatuses)
	if err != nil {
		fmt.Printf("applelink: could not parse reponse body: %s\n", err)
		return appCurrentStatuses, err
	}

	return appCurrentStatuses, err
}

func getBetaGroups(credentials *types.AppleCredentials) ([]types.BetaGroup, error) {
	var betaGroups []types.BetaGroup

	requestURL := fmt.Sprintf("%s/apple/connect/v1/apps/%s/groups", applelinkHost, credentials.BundleID)

	body, err := applelinkRequest(credentials, requestURL, http.MethodGet)
	err = json.Unmarshal(body, &betaGroups)
	if err != nil {
		fmt.Printf("applelink: could not parse reponse body: %s\n", err)
		return betaGroups, err
	}

	return betaGroups, err
}

func getInflightRelease(credentials *types.AppleCredentials) (types.Release, error) {
	var inflightRelease types.Release

	requestURL := fmt.Sprintf("%s/apple/connect/v1/apps/%s/release", applelinkHost, credentials.BundleID)

	body, err := applelinkRequest(credentials, requestURL, http.MethodGet)
	err = json.Unmarshal(body, &inflightRelease)
	if err != nil {
		fmt.Printf("applelink: could not parse reponse body: %s\n", err)
		return inflightRelease, err
	}

	return inflightRelease, err
}

func getLiveRelease(credentials *types.AppleCredentials) (types.Release, error) {
	var liveRelease types.Release

	requestURL := fmt.Sprintf("%s/apple/connect/v1/apps/%s/release/live", applelinkHost, credentials.BundleID)

	body, err := applelinkRequest(credentials, requestURL, http.MethodGet)
	err = json.Unmarshal(body, &liveRelease)
	if err != nil {
		fmt.Printf("applelink: could not parse reponse body: %s\n", err)
		return liveRelease, err
	}

	return liveRelease, err
}

func pauseLiveRelease(credentials *types.AppleCredentials) (types.Release, error) {
	var liveRelease types.Release

	requestURL := fmt.Sprintf("%s/apple/connect/v1/apps/%s/release/live/rollout/pause", applelinkHost, credentials.BundleID)

	body, err := applelinkRequest(credentials, requestURL, http.MethodPatch)
	err = json.Unmarshal(body, &liveRelease)
	if err != nil {
		fmt.Printf("applelink: could not parse reponse body: %s\n", err)
		return liveRelease, err
	}

	return liveRelease, err
}

func resumeLiveRelease(credentials *types.AppleCredentials) (types.Release, error) {
	var liveRelease types.Release

	requestURL := fmt.Sprintf("%s/apple/connect/v1/apps/%s/release/live/rollout/resume", applelinkHost, credentials.BundleID)

	body, err := applelinkRequest(credentials, requestURL, http.MethodPatch)
	err = json.Unmarshal(body, &liveRelease)
	if err != nil {
		fmt.Printf("applelink: could not parse reponse body: %s\n", err)
		return liveRelease, err
	}

	return liveRelease, err
}

func releaseToAll(credentials *types.AppleCredentials) (types.Release, error) {
	var liveRelease types.Release

	requestURL := fmt.Sprintf("%s/apple/connect/v1/apps/%s/release/live/rollout/complete", applelinkHost, credentials.BundleID)

	body, err := applelinkRequest(credentials, requestURL, http.MethodPatch)
	err = json.Unmarshal(body, &liveRelease)
	if err != nil {
		fmt.Printf("applelink: could not parse reponse body: %s\n", err)
		return liveRelease, err
	}

	return liveRelease, err
}

func getApplelinkAuthToken(credentials types.ApplelinkCredentials) string {
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
	return signedToken
}

func getAppStoreToken(credentials *types.AppleCredentials) (string, error) {
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
