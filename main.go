package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	authorizedUserKey = "AUTHORIZED_USER_EMAIL"
)

var (
	sessionName   string
	dbName        string
	clientID      string
	clientSecret  string
	redirectURL   string
	sessionSecret string
	db            *gorm.DB
	oauthConf     *oauth2.Config
)

type User struct {
	Email      string `gorm:"primary_key"`
	ProviderID string `gorm:"index:idx_name,unique"`
	Provider   string
	Name       string
	AvatarURL  string
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
	redirectURL = os.Getenv("AUTH_REDIRECT_URL")
	sessionSecret = os.Getenv("SECRET_SESSION_KEY")
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

func initOAuthConf() *oauth2.Config {
	return &oauth2.Config{
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

func initServer(db *gorm.DB, oauthConf *oauth2.Config) {
	r := gin.Default()

	// Set up sessions middleware
	store := cookie.NewStore([]byte(sessionSecret))
	r.Use(sessions.Sessions(sessionName, store))

	// Set up routes
	r.GET("/", handleHome(db, oauthConf))
	r.GET("/auth/google/callback", handleGoogleCallback(db, oauthConf))
	r.GET("/logout", handleLogout())

	// Serve the static files
	r.Static("/static", "./static")

	// Load the HTML templates
	r.LoadHTMLGlob("templates/*")

	// Start the server
	err := r.Run(":8080")
	if err != nil {
		fmt.Println(err)
	}
}

func handleHome(db *gorm.DB, oauthConf *oauth2.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the user is authorized
		if !isUserSessionAuthorized(c) {
			// If not, redirect to the login page
			authURL := oauthConf.AuthCodeURL("state")
			c.Redirect(http.StatusFound, authURL)
			return
		}

		// If authorized, display the user's name and avatar
		user := getUserFromDB(c)
		c.HTML(http.StatusOK, "dashboard.html", gin.H{"user": user})
	}
}

func handleGoogleCallback(db *gorm.DB, oauthConf *oauth2.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the authorization code from the query parameters
		code := c.Query("code")

		// Exchange the authorization code for a token
		token, err := oauthConf.Exchange(c, code)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		// Get the user's profile from the Google API
		client := oauthConf.Client(c, token)
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
			Name:       profile.Name,
			AvatarURL:  profile.AvatarURL,
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
	initServer(initDB(dbName), initOAuthConf())
}

func isUserSessionAuthorized(c *gin.Context) bool {
	session := sessions.Default(c)
	email := session.Get(authorizedUserKey)

	return email != nil
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
