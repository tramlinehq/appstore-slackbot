package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
)

func isUserSessionAuthorized(c *gin.Context) bool {
	session := sessions.Default(c)
	email := session.Get(authorizedUserKey)

	return email != nil
}

func isSlackConnected(c *gin.Context) bool {
	email := getEmailFromSession(c)
	var user User
	result := db.Where("email = ? AND slack_access_token NOT NULL", email).First(&user)
	return result.Error == nil
}

func isAppStoreConnected(c *gin.Context) bool {
	email := getEmailFromSession(c)
	var user User
	result := db.Where("email = ? AND app_store_connected = 1", email).First(&user)
	return result.Error == nil
}

func getUserFromDB(c *gin.Context) *User {
	email := getEmailFromSession(c)
	var user User
	db.Where("email = ?", email).First(&user)

	return &user
}

func getEmailFromSession(c *gin.Context) string {
	session := sessions.Default(c)
	email := session.Get(authorizedUserKey)

	if email == nil {
		return ""
	}

	return email.(string)
}

func encrypt(data []byte, key []byte, iv []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return ciphertext
}

func decrypt(data []byte, key []byte, iv []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	plaintext := make([]byte, len(data)-aes.BlockSize)
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(plaintext, data[aes.BlockSize:])

	return plaintext
}

func validateAppStoreCreds(bundleID string, issuerID string, keyID string, p8FileBytes []byte) error {
	appleCredentials := AppleCredentials{
		BundleID: bundleID,
		IssuerID: issuerID,
		KeyID:    keyID,
		P8File:   p8FileBytes,
	}
	appMetadata, err := GetAppMetadata(appleCredentials)
	fmt.Println(appMetadata)
	return err
}
