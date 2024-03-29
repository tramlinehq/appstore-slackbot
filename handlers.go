package main

import (
	"ciderbot/types"
	"crypto/aes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const LANDING_PAGE_URL = "https://appstoreslackbot.com"
const APP_STORE_CONNECT_FAILURE = "Failed to connect to the App Store! Please try again."
const SLACK_CONNECT_FAILURE = "Failed to authorize Slack! Please try again."

func handleAppStoreCreds() gin.HandlerFunc {
	return func(c *gin.Context) {
		bundleID := c.PostForm("bundle-id")
		issuerID := c.PostForm("issuer-id")
		keyID := c.PostForm("key-id")
		file, err := c.FormFile("p8-file")
		if err != nil {
			redirectWithFlash(c, APP_STORE_CONNECT_FAILURE)
			return
		}

		// Open the uploaded file
		fileData, err := file.Open()
		if err != nil {
			redirectWithFlash(c, APP_STORE_CONNECT_FAILURE)
			return
		}
		defer fileData.Close()

		// Read the file data
		p8FileBytes, err := io.ReadAll(fileData)
		if err != nil {
			redirectWithFlash(c, APP_STORE_CONNECT_FAILURE)
			return
		}

		// Validate app store creds
		err = validateAppStoreCreds(bundleID, issuerID, keyID, p8FileBytes)
		if err != nil {
			fmt.Printf("Validation of apple credentials failed with : %s\n", err)
			redirectWithFlash(c, APP_STORE_CONNECT_FAILURE)
			return
		}

		// Generate a random IV (Initialization Vector)
		iv := make([]byte, aes.BlockSize)
		if _, err := rand.Read(iv); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// Encrypt the file data using AES encryption
		encryptedP8File := encrypt(p8FileBytes, []byte(encryptionKey), iv)

		userValue, _ := c.Get("user")
		user, _ := userValue.(*types.User)
		user.AppStoreKeyID = sql.NullString{String: keyID, Valid: true}
		user.AppStoreBundleID = sql.NullString{String: bundleID, Valid: true}
		user.AppStoreIssuerID = sql.NullString{String: issuerID, Valid: true}
		user.AppStoreP8File = encryptedP8File
		user.AppStoreP8FileIV = iv
		user.AppStoreConnected = true

		result := db.Save(&user)
		if result.Error != nil {
			c.AbortWithError(http.StatusInternalServerError, result.Error)
			return
		}

		c.Redirect(http.StatusFound, "/")
	}
}

func handleSlackAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authURL := slackOAuthConf.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
		c.Redirect(http.StatusTemporaryRedirect, authURL)
	}
}

func handleSlackAuthCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Query("code")
		userValue, _ := c.Get("user")
		user, _ := userValue.(*types.User)

		token, err := slackOAuthConf.Exchange(c, code)
		if err != nil {
			redirectWithFlash(c, SLACK_CONNECT_FAILURE)
			return
		}

		team := token.Extra("team").(map[string]interface{})
		slackTeamID := team["id"].(string)
		slackTeamName := team["name"].(string)
		user.SlackAccessToken = sql.NullString{String: token.AccessToken, Valid: true}
		user.SlackRefreshToken = sql.NullString{String: token.RefreshToken, Valid: true}
		user.SlackTeamID = sql.NullString{String: slackTeamID, Valid: true}
		user.SlackTeamName = sql.NullString{String: slackTeamName, Valid: true}
		result := db.Save(&user)
		if result.Error != nil {
			c.AbortWithError(http.StatusInternalServerError, result.Error)
			return
		}

		c.Redirect(http.StatusFound, "/")
	}
}

func handleHome() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := getUserFromSession(c)
		session := sessions.Default(c)
		flashMessages := session.Flashes()
		session.Save()

		if err != nil {
			c.HTML(http.StatusOK, "login.html", gin.H{"user": user})
		} else {
			c.HTML(http.StatusOK, "index.html", gin.H{"user": user, "flashErrors": flashMessages})
		}
	}
}

func handleLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := getUserFromSession(c)
		if err != nil {
			authURL := googleOAuthConf.AuthCodeURL("state")
			c.Redirect(http.StatusFound, authURL)
			return
		} else {
			c.Redirect(http.StatusFound, "/")
		}
	}
}

func handleGoogleCallback(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the authorization code from the query parameters
		code := c.Query("code")

		// Exchange the authorization code for a token
		token, err := googleOAuthConf.Exchange(c, code)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		// Get the user's profile from the Google API
		client := googleOAuthConf.Client(c, token)
		resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		defer resp.Body.Close()

		var profile struct {
			ID        string `json:"id"`
			Email     string `json:"email"`
			Name      string `json:"name"`
			AvatarURL string `json:"picture"`
		}
		err = json.NewDecoder(resp.Body).Decode(&profile)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		// Save the authorized user to the database
		user := types.User{
			Provider:   "google",
			ProviderID: profile.ID,
			Email:      profile.Email,
			Name:       sql.NullString{String: profile.Name, Valid: true},
			AvatarURL:  sql.NullString{String: profile.AvatarURL, Valid: true},
		}

		result := db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "provider_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"email",
				"name",
				"avatar_url",
			}),
		}).Create(&user)
		if result.Error != nil {
			c.AbortWithError(http.StatusInternalServerError, result.Error)
			return
		}

		// Set the authorized user in the session
		session := sessions.Default(c)
		session.Set(authorizedUserKey, profile.Email)
		session.Save()

		// Redirect back to the home page
		c.Redirect(http.StatusFound, "/")
	}
}

func handleLogout() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Clear the authorized user from the session
		session := sessions.Default(c)
		session.Delete(authorizedUserKey)
		session.Save()

		// Redirect back to the home page
		c.Redirect(http.StatusFound, LANDING_PAGE_URL)
	}
}

func handleDeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userValue, _ := c.Get("user")
		user, _ := userValue.(*types.User)
		tx := db.Begin()
		result := tx.Delete(&user)
		if result.Error != nil {
			tx.Rollback()
			c.AbortWithError(http.StatusInternalServerError, result.Error)
			return
		}

		var metrics types.Metrics
		result = tx.First(&metrics)
		if result.Error != nil {
			metrics = types.Metrics{
				ID:           1,
				DeletedUsers: 1,
			}
			result = tx.Create(&metrics)
		} else {
			metrics.DeletedUsers += 1
			result = tx.Save(metrics)
		}

		if result.Error != nil {
			tx.Rollback()
			c.AbortWithError(http.StatusInternalServerError, result.Error)
			return
		}

		tx.Commit()

		// Clear the authorized user from the session
		session := sessions.Default(c)
		session.Delete(authorizedUserKey)
		session.Save()

		c.Redirect(http.StatusFound, LANDING_PAGE_URL)
	}
}

func handleSlackCommands() gin.HandlerFunc {
	return func(c *gin.Context) {
		signature := c.GetHeader("X-Slack-Signature")
		timestamp := c.GetHeader("X-Slack-Request-Timestamp")

		defer c.Request.Body.Close()
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.Request.Body = io.NopCloser(strings.NewReader(string(body)))

		var form types.SlackFormData
		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Verify signature
		if !verifyRequestSignature(signature, timestamp, body) {
			c.JSON(http.StatusOK, gin.H{"message": "Could not verify request!"})
			return
		}

		// Verify timestamp drift
		if !verifyRequestRecency(timestamp) {
			c.JSON(http.StatusOK, gin.H{"message": "Could not verify request!"})
			return
		}

		// Validate the token
		if form.Token != slackVerificationToken {
			c.JSON(http.StatusOK, gin.H{"message": "Could not verify request!"})
			return
		}

		// Verify valid team
		user := getUserByTeamID(form.TeamId)
		if user == nil {
			c.JSON(http.StatusOK, gin.H{"message": "who are you?"})
			return
		}

		c.JSON(http.StatusOK, handleSlackCommand(form, user))
	}
}

func getUserFromSession(c *gin.Context) (*types.User, error) {
	session := sessions.Default(c)
	email := session.Get(authorizedUserKey)

	if email == nil {
		return nil, fmt.Errorf("no user found")
	}

	var user types.User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func getUserByTeamID(teamID string) *types.User {
	var user types.User
	db.Where("slack_team_id = ?", teamID).First(&user)

	return &user
}

func getUserFromSessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := getUserFromSession(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("no user found"))
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func verifyRequestRecency(timestampStr string) bool {
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return false
	}

	timeDiff := time.Now().Unix() - timestamp
	absoluteDiff := timeDiff
	if timeDiff < 0 {
		absoluteDiff = -absoluteDiff
	}

	return absoluteDiff < 5*60
}

func verifyRequestSignature(signature, timestamp string, body []byte) bool {
	// Concatenate the timestamp and request body
	baseString := "v0:" + timestamp + ":" + string(body)

	// Create a HMAC-SHA256 hash using the Slack signing secret
	mac := hmac.New(sha256.New, []byte(slackSigningSecret))
	mac.Write([]byte(baseString))
	expectedSignature := "v0=" + hex.EncodeToString(mac.Sum(nil))

	// Compare the expected signature with the received signature
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

func validateAppStoreCreds(bundleID string, issuerID string, keyID string, p8FileBytes []byte) error {
	appleCredentials := types.AppleCredentials{
		BundleID: bundleID,
		IssuerID: issuerID,
		KeyID:    keyID,
		P8File:   p8FileBytes,
	}
	appMetadata, err := getAppMetadata(&appleCredentials)
	fmt.Println(appMetadata)
	return err
}

func redirectWithFlash(c *gin.Context, flashMsg string) {
	session := sessions.Default(c)
	session.AddFlash(flashMsg)
	session.Save()

	c.Redirect(http.StatusFound, "/")
}
