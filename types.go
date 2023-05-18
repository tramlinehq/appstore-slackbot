package main

import (
	"database/sql"
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
