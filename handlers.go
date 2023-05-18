package main

import (
	"crypto/aes"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io"
	"log"
	"net/http"
)

func handleAppStoreCreds() gin.HandlerFunc {
	return func(c *gin.Context) {
		bundleID := c.PostForm("bundle-id")
		issuerID := c.PostForm("issuer-id")
		keyID := c.PostForm("key-id")
		file, err := c.FormFile("p8-file")
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		// Open the uploaded file
		fileData, err := file.Open()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		defer fileData.Close()

		// Read the file data
		p8FileBytes, err := io.ReadAll(fileData)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// Validate app store creds
		err = validateAppStoreCreds(bundleID, issuerID, keyID, p8FileBytes)
		if err != nil {
			fmt.Printf("validation of apple credentials failed with : %s\n", err)
			c.AbortWithError(http.StatusUnprocessableEntity, err)
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

		email := getEmailFromSession(c)
		user := User{}
		result := db.Where("email = ?", email).First(&user)
		if result.Error != nil {
			c.AbortWithError(http.StatusInternalServerError, result.Error)
			return
		}

		user.AppStoreKeyID = sql.NullString{String: keyID, Valid: true}
		user.AppStoreBundleID = sql.NullString{String: bundleID, Valid: true}
		user.AppStoreIssuerID = sql.NullString{String: issuerID, Valid: true}
		user.AppStoreP8File = encryptedP8File
		user.AppStoreP8FileIV = iv
		user.AppStoreConnected = true

		result = db.Save(&user)
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

		token, err := slackOAuthConf.Exchange(c, code)
		if err != nil {
			log.Println(err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		email := getEmailFromSession(c)

		// Update the user's Slack access token in the database
		user := User{}
		result := db.Where("email = ?", email).First(&user)
		if result.Error != nil {
			c.AbortWithError(http.StatusInternalServerError, result.Error)
			return
		}

		user.SlackAccessToken = sql.NullString{String: token.AccessToken, Valid: true}
		result = db.Save(&user)
		if result.Error != nil {
			c.AbortWithError(http.StatusInternalServerError, result.Error)
			return
		}

		c.Redirect(http.StatusFound, "/")
	}
}

func handleHome(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the user is authorized
		if !isUserSessionAuthorized(c) {
			// If not, redirect to the login page
			authURL := googleOAuthConf.AuthCodeURL("state")
			c.Redirect(http.StatusFound, authURL)
			return
		}

		// If authorized, display the user's name and avatar
		user := getUserFromDB(c)
		c.HTML(http.StatusOK, "dashboard.html", gin.H{"user": user, "isSlackConnected": isSlackConnected(c), "isAppStoreConnected": isAppStoreConnected(c)})
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
		user := User{
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
		c.Redirect(http.StatusFound, "/")
	}
}
