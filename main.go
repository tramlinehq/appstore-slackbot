package main

import (
	"ciderbot/types"
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
)

const (
	authorizedUserKey = "AUTHORIZED_USER_EMAIL"
	certFilePath      = "./config/certs/localhost.pem"
	certKeyFilePath   = "./config/certs/localhost-key.pem"
	appPort           = ":8080"
)

var (
	sessionName            string
	dbName                 string
	clientID               string
	clientSecret           string
	redirectURL            string
	sessionSecret          string
	db                     *gorm.DB
	googleOAuthConf        *oauth2.Config
	slackOAuthConf         *oauth2.Config
	slackClientID          string
	slackClientSecret      string
	slackSigningSecret     string
	slackVerificationToken string
	slackRedirectURI       string
	appEnv                 string
	encryptionKey          string
	applelinkAuthAud       string
	applelinkAuthIssuer    string
	applelinkAuthSecret    string
	applelinkCredentials   *types.ApplelinkCredentials
	applelinkHost          string
)

func initEnv() {
	e := godotenv.Load()

	if e != nil {
		log.Fatalf("Error loading .env file: %s", e)
	}

	appEnv = os.Getenv("ENV")
	sessionName = os.Getenv("APP_NAME")
	dbName = os.Getenv("DB_NAME")
	clientID = os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL = os.Getenv("GOOGLE_REDIRECT_URL")
	sessionSecret = os.Getenv("SECRET_SESSION_KEY")
	slackClientID = os.Getenv("SLACK_CLIENT_ID")
	slackClientSecret = os.Getenv("SLACK_CLIENT_SECRET")
	slackSigningSecret = os.Getenv("SLACK_SIGNING_SECRET")
	slackVerificationToken = os.Getenv("SLACK_VERIFICATION_TOKEN")
	slackRedirectURI = os.Getenv("SLACK_REDIRECT_URL")
	encryptionKey = os.Getenv("ENCRYPTION_KEY")
	applelinkAuthAud = os.Getenv("APPLELINK_AUTH_AUD")
	applelinkAuthIssuer = os.Getenv("APPLELINK_AUTH_ISSUER")
	applelinkAuthSecret = os.Getenv("APPLELINK_AUTH_SECRET")
	applelinkHost = os.Getenv("APPLELINK_HOST")
}

// TODO: do we need to close the DB "conn"?
func initDB(name string) *gorm.DB {
	var err error
	db, err = gorm.Open(sqlite.Open(name), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&types.User{})
	db.AutoMigrate(&types.Metrics{})

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
	applelinkCredentials = &types.ApplelinkCredentials{
		Aud:    applelinkAuthAud,
		Issuer: applelinkAuthIssuer,
		Secret: applelinkAuthSecret,
	}
}

func handlePing() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	}
}

func initServer(db *gorm.DB) {
	r := gin.Default()

	store := cookie.NewStore([]byte(sessionSecret))
	r.Use(sessions.Sessions(sessionName, store))

	r.GET("/", handleHome())
	r.GET("/logout", handleLogout())
	r.GET("/auth/google/callback", handleGoogleCallback(db))
	r.GET("/auth/slack/start", handleSlackAuth())
	r.GET("/auth/slack/callback", getUserFromSessionMiddleware(), handleSlackAuthCallback())
	r.POST("/auth/apple", getUserFromSessionMiddleware(), handleAppStoreCreds())
	r.POST("/user/delete", getUserFromSessionMiddleware(), handleDeleteUser())
	r.GET("/ping", handlePing())
	r.POST("/slack/listen", handleSlackCommands())

	r.Static("/assets", "./assets")
	r.LoadHTMLGlob("views/*")

	var err error
	if appEnv == "production" {
		err = r.Run(appPort)
	} else {
		err = r.RunTLS(appPort, certFilePath, certKeyFilePath)
	}

	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	initEnv()
	initSlackOAuthConf()
	initGoogleOAuthConf()
	initApplelinkCreds()
	initServer(initDB(dbName))
}
