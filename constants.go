package main

import (
	"golang.org/x/oauth2"
	"gorm.io/gorm"
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
	applelinkHost        string
)
