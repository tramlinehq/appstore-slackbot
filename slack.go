package main

import (
	"fmt"
)

func handleSlackCommand(form SlackFormData, user *User) SlackResponse {
	switch form.Text {
	case "info":
		return handleInfoCommand(user)
	case "live":
		return createSlackResponse([]string{"Nothing *live* yet, _bruhahahaha_."})
	default:
		return createSlackResponse([]string{"Please input a valid command"})

	}
}

func handleInfoCommand(user *User) SlackResponse {
	appleCredentials := AppleCredentials{
		BundleID: user.AppStoreBundleID.String,
		IssuerID: user.AppStoreIssuerID.String,
		KeyID:    user.AppStoreKeyID.String,
		P8File:   decrypt(user.AppStoreP8File, []byte(encryptionKey), user.AppStoreP8FileIV),
	}
	appMetadata, err := GetAppMetadata(appleCredentials)
	if err != nil {
		return createSlackResponse([]string{"Could not find an app."})

	}

	return createSlackResponse([]string{fmt.Sprintf("App Info:\nName: %s\nSKU: %s\nBundle ID: %s\nID: %s", appMetadata.Name, appMetadata.Sku, appMetadata.BundleId, appMetadata.Id)})
}

func createSlackResponse(messages []string) SlackResponse {
	blocks := make([]SlackResponseText, len(messages))
	for i, message := range messages {
		blocks[i] = SlackResponseText{
			Type: "section",
			Text: SlackResponseInsideText{Type: "mrkdwn", Text: message},
		}
	}

	return SlackResponse{Blocks: blocks}
}
