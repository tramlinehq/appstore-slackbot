package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	dbName         = "db.db"
	clientID       = "946207521855-fgkqq6p31mg5t5ql2tqe48knskhk5a2f.apps.googleusercontent.com"
	clientSecret   = "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"
	redirectURL    = "http://localhost:8080/auth/google/callback"
	sessionName    = "ciderbot"
	sessionSecret  = "SESSION_SECRET_KEY"
	authorizedUser = "AUTHORIZED_USER_EMAIL"
)

var (
	db        *gorm.DB
	oauthConf *oauth2.Config
)

type User struct {
	Email      string `gorm:"primary_key"`
	ProviderID string `gorm:"index:idx_name,unique"`
	Provider   string
	Name       string
	AvatarURL  string
}

func main() {
	// Initialize the database
	db = initDB(dbName)
	// defer db.Close()

	// Initialize the OAuth configuration
	oauthConf = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	// Initialize the Gin router
	r := gin.Default()

	// Set up sessions middleware
	store := cookie.NewStore([]byte(sessionSecret))
	r.Use(sessions.Sessions(sessionName, store))

	// Set up routes
	r.GET("/", func(c *gin.Context) {
		// Check if the user is authorized
		if !isAuthorized(c) {
			// If not, redirect to the login page
			authURL := oauthConf.AuthCodeURL("state")
			c.Redirect(http.StatusFound, authURL)
			return
		}

		// If authorized, display the user's name and avatar
		user := getUser(c)
		c.HTML(http.StatusOK, "dashboard.html", gin.H{"user": user})
	})

	r.GET("/auth/google/callback", func(c *gin.Context) {
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
		session.Set(authorizedUser, profile.Email)
		session.Save()

		// Redirect back to the home page
		c.Redirect(http.StatusFound, "/")
	})

	r.GET("/logout", func(c *gin.Context) {
		// Clear the authorized user from the session
		session := sessions.Default(c)
		session.Delete(authorizedUser)
		session.Save()

		// Redirect back to the home page
		c.Redirect(http.StatusFound, "/")
	})

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

func initDB(name string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(name), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{})

	return db
}

func isAuthorized(c *gin.Context) bool {
	session := sessions.Default(c)
	email := session.Get(authorizedUser)

	return email != nil
}

func getUser(c *gin.Context) *User {
	email := getEmail(c)

	var user User
	db.Where("email = ?", email).First(&user)

	return &user
}

func getEmail(c *gin.Context) string {
	session := sessions.Default(c)
	email := session.Get(authorizedUser)

	if email == nil {
		return ""
	}

	return email.(string)
}
