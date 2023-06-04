package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var ValidSlackCommands = map[string]string{
	"app_info":            "app_info",
	"overall_status":      "overall_status",
	"beta_groups":         "beta_groups",
	"inflight_release":    "inflight_release",
	"live_release":        "live_release",
	"pause_live_release":  "pause_live_release",
	"resume_live_release": "resume_live_release",
	"release_to_all":      "release_to_all",
}

func handleSlackCommand(form SlackFormData, user *User) SlackResponse {
	command := strings.Split(form.Text, " ")[0]
	if command, ok := ValidSlackCommands[command]; ok == true {
		go handleValidSlackCommand(command, form.ResponseUrl, user)
		return createEphemeralSlackResponse(fmt.Sprintf("Got %s command, working on it.", command))
	}

	return createEphemeralSlackResponse("Please input a valid command")
}

func handleValidSlackCommand(command string, responseURL string, user *User) {
	slackResponse := processValidSlackCommand(command, user)
	sendResponseToSlack(responseURL, slackResponse)
}

func processValidSlackCommand(command string, user *User) SlackResponse {
	switch command {
	case "app_info":
		return handleInfoCommand(user)
	case "overall_status":
		return handleCurrentStatusCommand(user)
	case "beta_groups":
		return handleBetaGroupsCommand(user)
	case "inflight_release":
		return handleInflightReleaseCommand(user)
	case "live_release":
		return handleLiveReleaseCommand(user)
	case "pause_live_release":
		return handlePauseReleaseCommand(user)
	case "resume_live_release":
		return handleResumeReleaseCommand(user)
	case "release_to_all":
		return handleReleaseToAllCommand(user)
	default:
		return createEphemeralSlackResponse("Please input a valid command")

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
		return createEphemeralSlackResponse("Could not find an app.")
	}

	return createSlackResponse("App Info", []string{fmt.Sprintf(`Name: %s
SKU: %s
Bundle ID: %s
ID: %s`,
		appMetadata.Name, appMetadata.Sku, appMetadata.BundleId, appMetadata.Id)})
}

func handleCurrentStatusCommand(user *User) SlackResponse {
	appCurrentStatuses, err := getAppCurrentStatus(userAppleCredentials(user))
	if err != nil {
		return createEphemeralSlackResponse("Could not find an app.")
	}

	var slackMessages []string
	for _, appCurrentStatus := range appCurrentStatuses {
		channelMessage := fmt.Sprintf("Channel: %s\n", appCurrentStatus.Name)
		for _, build := range appCurrentStatus.Builds {
			channelMessage += fmt.Sprintf("Build: %s (%s) - %s - %s\n", build.VersionString, build.BuildNumber, build.Status, build.ReleaseDate)
		}
		slackMessages = append(slackMessages, channelMessage)
	}

	return createSlackResponse("Current Store Status for your app", slackMessages)
}

func handleBetaGroupsCommand(user *User) SlackResponse {
	betaGroups, err := getBetaGroups(userAppleCredentials(user))
	if err != nil {
		return createEphemeralSlackResponse("Could not find an app.")
	}

	var slackMessages []string
	for _, betaGroup := range betaGroups {
		groupType := "external"

		if betaGroup.Internal == true {
			groupType = "internal"
		}
		groupMessage := fmt.Sprintf("*%s* -- *%s* group with *%d* testers.", betaGroup.Name, groupType, len(betaGroup.Testers))
		slackMessages = append(slackMessages, groupMessage)
	}

	return createSlackResponse("Test Groups", slackMessages)
}

func handleInflightReleaseCommand(user *User) SlackResponse {
	appInfo, err := getAppMetadata(userAppleCredentials(user))
	if err != nil {
		return createEphemeralSlackResponse("Could not find your app.")
	}
	inflightRelease, err := getInflightRelease(userAppleCredentials(user))
	if err != nil {
		return createEphemeralSlackResponse("Could not find an inflight release for your app.")
	}

	phasedReleaseEnabled := "on"

	if inflightRelease.PhasedRelease.Id == "" {
		phasedReleaseEnabled = "off"
	}

	return createSlackResponse("Inflight Release",
		[]string{
			fmt.Sprintf("Next up for release is *%s (%s)* with current status `%s`.",
				inflightRelease.VersionName,
				inflightRelease.BuildNumber,
				inflightRelease.AppStoreState),
			fmt.Sprintf("The release type is `%s` and phased release is turned *%s*.", inflightRelease.ReleaseType, phasedReleaseEnabled),
			fmt.Sprintf("<https://appstoreconnect.apple.com/apps/%s/appstore/ios/version/inflight|App Store Connect>", appInfo.Id)})
}

func handleLiveReleaseCommand(user *User) SlackResponse {
	appInfo, err := getAppMetadata(userAppleCredentials(user))
	if err != nil {
		return createEphemeralSlackResponse("Could not find your app.")
	}
	liveRelease, err := getLiveRelease(userAppleCredentials(user))
	if err != nil {
		return createEphemeralSlackResponse("Could not find a live release for your app.")
	}

	return createSlackResponse("Live Release",
		[]string{
			fmt.Sprintf("*%s (%s)* is on day *%d* of phased release with status `%s`.",
				liveRelease.VersionName,
				liveRelease.BuildNumber,
				liveRelease.PhasedRelease.CurrentDayNumber,
				liveRelease.PhasedRelease.PhasedReleaseState),
			fmt.Sprintf("<https://appstoreconnect.apple.com/apps/%s/appstore/ios/version/deliverable|App Store Connect>", appInfo.Id)})
}

func handlePauseReleaseCommand(user *User) SlackResponse {
	appInfo, err := getAppMetadata(userAppleCredentials(user))
	if err != nil {
		return createEphemeralSlackResponse("Could not find your app.")
	}
	liveRelease, err := pauseLiveRelease(userAppleCredentials(user))
	if err != nil {
		return createEphemeralSlackResponse("Could not find an live release to pause.")
	}

	return createSlackResponse("Live Release",
		[]string{
			fmt.Sprintf("*%s (%s)* is on day *%d* of phased release with status `%s`.",
				liveRelease.VersionName,
				liveRelease.BuildNumber,
				liveRelease.PhasedRelease.CurrentDayNumber,
				liveRelease.PhasedRelease.PhasedReleaseState),
			fmt.Sprintf("<https://appstoreconnect.apple.com/apps/%s/appstore/ios/version/deliverable|App Store Connect>", appInfo.Id)})
}

func handleResumeReleaseCommand(user *User) SlackResponse {
	appInfo, err := getAppMetadata(userAppleCredentials(user))
	if err != nil {
		return createEphemeralSlackResponse("Could not find your app.")
	}
	liveRelease, err := resumeLiveRelease(userAppleCredentials(user))
	if err != nil {
		return createEphemeralSlackResponse("Could not find an paused release to resume for your app.")
	}

	return createSlackResponse("Live Release",
		[]string{
			fmt.Sprintf("*%s (%s)* is on day *%d* of phased release with status `%s`.",
				liveRelease.VersionName,
				liveRelease.BuildNumber,
				liveRelease.PhasedRelease.CurrentDayNumber,
				liveRelease.PhasedRelease.PhasedReleaseState),
			fmt.Sprintf("<https://appstoreconnect.apple.com/apps/%s/appstore/ios/version/deliverable|App Store Connect>", appInfo.Id)})
}

func handleReleaseToAllCommand(user *User) SlackResponse {
	appInfo, err := getAppMetadata(userAppleCredentials(user))
	if err != nil {
		return createEphemeralSlackResponse("Could not find your app.")
	}
	liveRelease, err := releaseToAll(userAppleCredentials(user))
	if err != nil {
		return createEphemeralSlackResponse("Could not find an live release to release to all for your app.")
	}

	return createSlackResponse("Live Release",
		[]string{
			fmt.Sprintf("*%s (%s)* is on day *%d* of phased release with status `%s`.",
				liveRelease.VersionName,
				liveRelease.BuildNumber,
				liveRelease.PhasedRelease.CurrentDayNumber,
				liveRelease.PhasedRelease.PhasedReleaseState),
			fmt.Sprintf("<https://appstoreconnect.apple.com/apps/%s/appstore/ios/version/deliverable|App Store Connect>", appInfo.Id)})
}

func createSlackResponse(header string, messages []string) SlackResponse {
	blocks := make([]SlackResponseText, len(messages)+1)
	blocks[0] = SlackResponseText{
		Type: "header",
		Text: SlackResponseInsideText{Type: "plain_text", Text: header},
	}
	for i, message := range messages {
		blocks[i+1] = SlackResponseText{
			Type: "section",
			Text: SlackResponseInsideText{Type: "mrkdwn", Text: message},
		}
	}

	return SlackResponse{Blocks: blocks, ResponseType: "in_channel"}
}

func createEphemeralSlackResponse(message string) SlackResponse {
	blocks := make([]SlackResponseText, 1)
	blocks[0] = SlackResponseText{
		Type: "section",
		Text: SlackResponseInsideText{Type: "plain_text", Text: message},
	}

	return SlackResponse{Blocks: blocks, ResponseType: "ephemeral"}
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
