package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	authorizedUserKey = "AUTHORIZED_USER_EMAIL"
)

var (
	sessionName          string
	dbName               string
	clientID             string
	clientSecret         string
	redirectURL          string
	sessionSecret        string
	db                   *gorm.DB
	googleOAuthConf      *oauth2.Config
	slackOAuthConf       *oauth2.Config
	slackClientID        string
	slackClientSecret    string
	slackRedirectURI     string
	appEnv               string
	encryptionKey        string
	applelinkAuthAud     string
	applelinkAuthIssuer  string
	applelinkAuthSecret  string
	applelinkCredentials *ApplelinkCredentials
)

type User struct {
	Email             string `gorm:"primary_key"`
	ProviderID        string `gorm:"index:idx_name,unique"`
	Provider          string
	Name              sql.NullString
	AvatarURL         sql.NullString
	SlackAccessToken  sql.NullString
	AppStoreBundleID  sql.NullString
	AppStoreIssuerID  sql.NullString
	AppStoreKeyID     sql.NullString
	AppStoreConnected bool `gorm:"default:false"`
	AppStoreP8File    []byte
	AppStoreP8FileIV  []byte
}

func initEnv() {
	e := godotenv.Load()

	if e != nil {
		log.Fatalf("Error loading .env file: %s", e)
	}

	sessionName = os.Getenv("APP_NAME")
	dbName = os.Getenv("DB_NAME")
	clientID = os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL = os.Getenv("GOOGLE_REDIRECT_URL")
	sessionSecret = os.Getenv("SECRET_SESSION_KEY")
	slackClientID = os.Getenv("SLACK_CLIENT_ID")
	slackClientSecret = os.Getenv("SLACK_CLIENT_SECRET")
	slackRedirectURI = os.Getenv("SLACK_REDIRECT_URL")
	appEnv = os.Getenv("ENV")
	encryptionKey = os.Getenv("ENCRYPTION_KEY")
	applelinkAuthAud = os.Getenv("APPLELINK_AUTH_AUD")
	applelinkAuthIssuer = os.Getenv("APPLELINK_AUTH_ISSUER")
	applelinkAuthSecret = os.Getenv("APPLELINK_AUTH_SECRET")
}

// TODO: do we need to close the DB "conn"?
func initDB(name string) *gorm.DB {
	var err error
	db, err = gorm.Open(sqlite.Open(name), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{})

	return db
}

func initGoogleOAuthConf() {
	googleOAuthConf = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func initSlackOAuthConf() {
	slackOAuthConf = &oauth2.Config{
		ClientID:     slackClientID,
		ClientSecret: slackClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://slack.com/oauth/v2/authorize",
			TokenURL: "https://slack.com/api/oauth.v2.access",
		},
		RedirectURL: slackRedirectURI,
		Scopes:      []string{"chat:write", "chat:write.customize", "commands"},
	}
}

func initApplelinkCreds() {
	applelinkCredentials = &ApplelinkCredentials{
		Aud:    applelinkAuthAud,
		Issuer: applelinkAuthIssuer,
		Secret: applelinkAuthSecret,
	}
}

func initServer(db *gorm.DB) {
	r := gin.Default()

	// Set up sessions middleware
	store := cookie.NewStore([]byte(sessionSecret))
	r.Use(sessions.Sessions(sessionName, store))

	// Set up routes
	r.GET("/", handleHome(db))
	r.GET("/logout", handleLogout())
	r.GET("/auth/google/callback", handleGoogleCallback(db))
	r.GET("/auth/slack/start", handleSlackAuth())
	r.GET("/auth/slack/callback", handleSlackAuthCallback())
	r.POST("/auth/apple", handleAppStoreCreds())

	// Serve the static files
	r.Static("/static", "./static")

	// Load the HTML templates
	r.LoadHTMLGlob("templates/*")

	// Start the server
	var err error
	if appEnv == "production" {
		err = r.Run(":8080")
	} else {
		err = r.RunTLS(":8080", "./config/certs/localhost.pem", "./config/certs/localhost-key.pem")
	}

	if err != nil {
		fmt.Println(err)
	}
}

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

		//Validate app store creds
		err = validateAppStoreCreds(bundleID, issuerID, keyID, p8FileBytes)
		if err != nil {
			fmt.Printf("validation of apple credebntaisl failed with : %s\n", err)
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

func main() {
	initEnv()
	initSlackOAuthConf()
	initGoogleOAuthConf()
	initApplelinkCreds()
	initServer(initDB(dbName))
}

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
	appMetadata, err := GetAppMetadata(*applelinkCredentials, appleCredentials)
	fmt.Println(appMetadata)
	return err
}
