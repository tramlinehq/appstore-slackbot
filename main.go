package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"github.com/shareed2k/goth_fiber"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	//"fmt"
)

var (
	clientID     string
	clientSecret string
	db           *gorm.DB
	tpl          = html.New("./templates/", ".html")
	app          *fiber
)

// CREATE TABLE IF NOT EXISTS users (
//
//	    id INTEGER PRIMARY KEY,
//	    provider TEXT NOT NULL,
//	    provider_id TEXT NOT NULL UNIQUE,
//	    name TEXT NOT NULL,
//	    email TEXT NOT NULL UNIQUE,
//	    avatar_url TEXT NOT NULL
//	);
type User struct {
	ID         uint `gorm:"primary_key"`
	Provider   string
	ProviderID string
	Email      string
	Name       string
	AvatarURL  string
}

func initEnv() {
	e := godotenv.Load()

	if e != nil {
		log.Fatalf("Error loading .env file: %s", e)
	}

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
}

func initGoth() {
	goth.UseProviders(
		google.New(clientID, clientSecret, "http://localhost:3000/auth/google/callback"),
	)
}

func initFiber() {
	app = fiber.New(fiber.Config{
		Views: tpl,
	})

	// Initialize default config
	app.Use(logger.New())
}

func initSession() {
	store = session.New()
}

func initRoutes() {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("login", fiber.Map{})
	})

	app.Get("/auth/:provider", goth_fiber.BeginAuthHandler)

	app.Get("/auth/google/callback", func(c *fiber.Ctx) error {
		user, err := goth_fiber.CompleteUserAuth(c)
		if err != nil {
			return err
		}
		err = saveUser(user)
		if err != nil {
			return err
		}

		// Store user session
		session, err := store.Get(c)
		if err != nil {
			return err
		}
		session.Set("userID", user.UserID)
		err = session.Save()
		if err != nil {
			return err
		}

		return c.Render("dashboard", fiber.Map{
			"user": user,
		})
	})

	app.Get("/logout", func(c *fiber.Ctx) error {
		// Destroy user session
		session, err := store.Get(c)
		if err != nil {
			return err
		}
		session.Destroy()
		err = session.Save()
		if err != nil {
			return err
		}

		return c.Redirect("/")
	})

	app.Get("/dashboard", func(c *fiber.Ctx) error {
		// Retrieve user session
		session, err := store.Get(c)
		if err != nil {
			return err
		}
		userID := session.Get("userID")
		if userID == nil {
			return c.Redirect("/")
		}

		// Retrieve user from database
		user, err := getUserByProviderID("google", userID.(string))
		if err != nil {
			return err
		}

		return c.Render("dashboard", fiber.Map{
			"user": user,
		})
	})
}

func initDB() {
	// Connect to SQLite database
	var err error
	db, err = gorm.Open(sqlite.Open("db.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// defer db.Close()
	db.AutoMigrate(&User{})
}

func initServer() {
	log.Fatal(app.Listen(":3000"))
}

func main() {
	initEnv()
	initFiber()
	initGoth()
	initSession()
	initRoutes()
	initDB()
	initServer()
}

func saveUser(user goth.User) error {
	var count int64
	result := db.Model(&User{}).Where("provider = ? AND provider_id = ?", user.Provider, user.UserID).Count(&count)
	if result.Error != nil {
		return result.Error
	}

	newUser := User{Name: user.Name, Email: user.Email, Provider: user.Provider, ProviderID: user.UserID, AvatarURL: user.AvatarURL}

	if count == 0 {
		result = db.Create(&newUser)
	} else {
		result = db.Model(&User{}).Where("provider = ? AND provider_id = ?", user.Provider, user.UserID).Updates(map[string]interface{}{
			"name":       user.Name,
			"email":      user.Email,
			"avatar_url": user.AvatarURL,
		})
	}

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func getUserByProviderID(provider string, providerID string) (goth.User, error) {
	var user goth.User
	result := db.Where("provider = ? AND provider_id = ?", provider, providerID).First(&User{})
	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}
