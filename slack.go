package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var ValidSlackCommands = map[string]string{
	"info":           "info",
	"current_status": "current_status",
	"live":           "live",
}

func handleSlackCommand(form SlackFormData, user *User) SlackResponse {
	if command, ok := ValidSlackCommands[form.Text]; ok == true {
		go handleValidSlackCommand(command, form.ResponseUrl, user)
		return createSlackResponse([]string{"Got it, working on it."}, "ephemeral")
	}

	return createSlackResponse([]string{"Please input a valid command"}, "ephemeral")

}

func handleValidSlackCommand(command string, responseURL string, user *User) {
	slackResponse := processValidSlackCommand(command, user)
	sendResponseToSlack(responseURL, slackResponse)
}

func processValidSlackCommand(command string, user *User) SlackResponse {
	switch command {
	case "info":
		return handleInfoCommand(user)
	case "current_status":
		return handleCurrentStatusCommand(user)
	case "live":
		return createSlackResponse([]string{"Nothing *live* yet, _bruhahahaha_."}, "ephemeral")
	default:
		return createSlackResponse([]string{"Please input a valid command"}, "ephemeral")

	}
}

func sendResponseToSlack(requestURL string, slackResponse SlackResponse) error {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(slackResponse)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, requestURL, &body)
	if err != nil {
		fmt.Printf("slack: could not create request: %s\n", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = http.DefaultClient.Do(req)

	if err != nil {
		fmt.Printf("slack: request failed: %s\n", err)
		return err
	}
	return nil
}

func handleInfoCommand(user *User) SlackResponse {
	appMetadata, err := getAppMetadata(userAppleCredentials(user))
	if err != nil {
		return createSlackResponse([]string{"Could not find an app."}, "ephemeral")
	}

	return createSlackResponse([]string{fmt.Sprintf(`App Info:
Name: %s
SKU: %s
Bundle ID: %s
ID: %s`,
		appMetadata.Name, appMetadata.Sku, appMetadata.BundleId, appMetadata.Id)}, "in_channel")
}

func handleCurrentStatusCommand(user *User) SlackResponse {
	appCurrentStatuses, err := getAppCurrentStatus(userAppleCredentials(user))
	if err != nil {
		return createSlackResponse([]string{"Could not find an app."}, "ephemeral")
	}

	var slackMessages []string
	for _, appCurrentStatus := range appCurrentStatuses {
		channelMessage := fmt.Sprintf("Channel: %s\n", appCurrentStatus.Name)
		for _, build := range appCurrentStatus.Builds {
			channelMessage += fmt.Sprintf("Build: %s (%s) - %s - %s\n", build.VersionString, build.BuildNumber, build.Status, build.ReleaseDate)
		}
		slackMessages = append(slackMessages, channelMessage)
	}

	return createSlackResponse(slackMessages, "in_channel")
}

func createSlackResponse(messages []string, responseType string) SlackResponse {
	blocks := make([]SlackResponseText, len(messages))
	for i, message := range messages {
		blocks[i] = SlackResponseText{
			Type: "section",
			Text: SlackResponseInsideText{Type: "mrkdwn", Text: message},
		}
	}

	return SlackResponse{Blocks: blocks, ResponseType: responseType}
}

func userAppleCredentials(user *User) *AppleCredentials {
	appleCredentials := AppleCredentials{
		BundleID: user.AppStoreBundleID.String,
		IssuerID: user.AppStoreIssuerID.String,
		KeyID:    user.AppStoreKeyID.String,
		P8File:   decrypt(user.AppStoreP8File, []byte(encryptionKey), user.AppStoreP8FileIV),
	}

	return &appleCredentials
}
