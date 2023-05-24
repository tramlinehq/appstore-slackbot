package main

import (
	"database/sql"
	"time"
)

type User struct {
	Email             string `gorm:"primary_key"`
	ProviderID        string `gorm:"index:idx_name,unique"`
	Provider          string
	Name              sql.NullString
	AvatarURL         sql.NullString
	SlackAccessToken  sql.NullString
	SlackRefreshToken sql.NullString
	SlackTeamID       sql.NullString
	SlackTeamName     sql.NullString
	AppStoreBundleID  sql.NullString
	AppStoreIssuerID  sql.NullString
	AppStoreKeyID     sql.NullString
	AppStoreConnected bool `gorm:"default:false"`
	AppStoreP8File    []byte
	AppStoreP8FileIV  []byte
}

type AppleCredentials struct {
	BundleID string
	IssuerID string
	KeyID    string
	P8File   []byte
}

type ApplelinkCredentials struct {
	Aud    string
	Issuer string
	Secret string
}

type AppMetadata struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	BundleId string `json:"bundle_id"`
	Sku      string `json:"sku"`
}

type AppCurrentStatus struct {
	Name   string `json:"name"`
	Builds []struct {
		Id            string    `json:"id"`
		BuildNumber   string    `json:"build_number"`
		Status        string    `json:"status"`
		VersionString string    `json:"version_string"`
		ReleaseDate   time.Time `json:"release_date"`
	} `json:"builds"`
}

type SlackFormData struct {
	Token          string `form:"token"`
	TeamId         string `form:"team_id"`
	TeamDomain     string `form:"team_domain"`
	EnterpriseId   string `form:"enterprise_id"`
	EnterpriseName string `form:"enterprise_name"`
	ChannelId      string `form:"channel_id"`
	ChannelName    string `form:"channel_name"`
	UserId         string `form:"user_id"`
	UserName       string `form:"user_name"`
	Command        string `form:"command"`
	Text           string `form:"text"`
	ResponseUrl    string `form:"response_url"`
	TriggerId      string `form:"trigger_id"`
	ApiAppId       string `form:"api_app_id"`
}

type SlackResponseInsideText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type SlackResponseText struct {
	Type string                  `json:"type"`
	Text SlackResponseInsideText `json:"text"`
}

type SlackResponse struct {
	Blocks       []SlackResponseText `json:"blocks"`
	ResponseType string              `json:"response_type"`
}
