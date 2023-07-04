package types

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
	CommandCount      int64
}

type Metrics struct {
	ID           int8 `gorm:"primary_key"`
	DeletedUsers int64
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

type BetaGroup struct {
	Name     string `json:"name"`
	Id       string `json:"id"`
	Internal bool   `json:"internal"`
	Testers  []struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"testers"`
}

type Release struct {
	Id                  string      `json:"id"`
	VersionName         string      `json:"version_name"`
	AppStoreState       string      `json:"app_store_state"`
	ReleaseType         string      `json:"release_type"`
	EarliestReleaseDate interface{} `json:"earliest_release_date"`
	Downloadable        bool        `json:"downloadable"`
	CreatedDate         time.Time   `json:"created_date"`
	BuildNumber         string      `json:"build_number"`
	BuildId             string      `json:"build_id"`
	PhasedRelease       struct {
		Id                 string    `json:"id"`
		PhasedReleaseState string    `json:"phased_release_state"`
		StartDate          time.Time `json:"start_date"`
		TotalPauseDuration int       `json:"total_pause_duration"`
		CurrentDayNumber   int       `json:"current_day_number"`
	} `json:"phased_release"`
	Details struct {
		Id              string      `json:"id"`
		Description     string      `json:"description"`
		Locale          string      `json:"locale"`
		Keywords        string      `json:"keywords"`
		MarketingUrl    interface{} `json:"marketing_url"`
		PromotionalText interface{} `json:"promotional_text"`
		SupportUrl      string      `json:"support_url"`
		WhatsNew        string      `json:"whats_new"`
	} `json:"details"`
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

type SlackResponse struct {
	Blocks       []Block `json:"blocks"`
	ResponseType string  `json:"response_type"`
}

type Block struct {
	Type     string    `json:"type"`
	Text     Text      `json:"text,omitempty"`
	Fields   []Text    `json:"fields,omitempty"`
	Elements []Element `json:"elements,omitempty"`
}

type Text struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Element struct {
	Type     string `json:"type"`
	ImageURL string `json:"image_url,omitempty"`
	AltText  string `json:"alt_text,omitempty"`
	Text     string `json:"text,omitempty"`
}
