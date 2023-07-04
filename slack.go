package main

import (
	"bytes"
	slack "ciderbot/slack"
	"ciderbot/types"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var ValidSlackCommands = map[string][2]string{
	"app_info":            {":information_source:", "Get some basic information about your app to verify you are working with the correct app"},
	"overall_status":      {":convenience_store:", "Get an overall store status for your app, what builds are distributed to which channels (TestFlight and AppStore)"},
	"inflight_release":    {":airplane_departure:", "Get the current inflight release in the App Store"},
	"help":                {":eyes:", "Get the usage guide for App Store SlackBot"},
	"beta_groups":         {":test_tube:", "List all the beta groups present in TestFlight"},
	"live_release":        {":iphone:", "Get the current live release in the App Store"},
	"pause_live_release":  {":double_vertical_bar:", "Pause the phased release of the current live release in the App Store"},
	"resume_live_release": {":arrow_forward:", "Resume the phased release of the current live release in the App Store"},
	"release_to_all":      {":roller_coaster:", "Release the current live release in the App Store to all users"},
}

func handleSlackCommand(form types.SlackFormData, user *types.User) types.SlackResponse {
	if !user.AppStoreBundleID.Valid {
		return slack.EphemeralMessage{Msg: "No iOS app registered. Please add ASC details to use appstoreslackbot."}.Render()
	}
	command := strings.Split(form.Text, " ")[0]
	if _, ok := ValidSlackCommands[command]; ok {
		user.CommandCount += 1
		db.Save(&user)
		go handleValidSlackCommand(command, form.ResponseUrl, user)
		return slack.EphemeralMessage{Msg: fmt.Sprintf("Got the `%s` command, working on it.", command)}.Render()
	}

	return slack.EphemeralMessage{Msg: "Please input a valid command. Use the `help` command to see all the valid commands."}.Render()
}

func handleValidSlackCommand(command string, responseURL string, user *types.User) {
	slackResponse := processValidSlackCommand(command, user)
	sendResponseToSlack(responseURL, slackResponse)
}

func processValidSlackCommand(command string, user *types.User) types.SlackResponse {
	switch command {
	case "help":
		return handleHelpCommand(user)
	case "app_info":
		return handleInfoCommand(user)
	// case "overall_status":
	// 	return handleCurrentStatusCommand(user)
	// case "beta_groups":
	// 	return handleBetaGroupsCommand(user)
	// case "inflight_release":
	// 	return handleInflightReleaseCommand(user)
	// case "live_release":
	// 	return handleLiveReleaseCommand(user)
	// case "pause_live_release":
	// 	return handlePauseReleaseCommand(user)
	// case "resume_live_release":
	// 	return handleResumeReleaseCommand(user)
	// case "release_to_all":
	// 	return handleReleaseToAllCommand(user)
	default:
		return slack.EphemeralMessage{Msg: "Please input a valid command"}.Render()
	}
}

func handleHelpCommand(_user *types.User) types.SlackResponse {
	return slack.HelpText{Commands: ValidSlackCommands}.Render()
}

func handleInfoCommand(user *types.User) types.SlackResponse {
	appMetadata, err := getAppMetadata(userAppleCredentials(user))
	if err != nil {
		return slack.EphemeralMessage{Msg: "Could not find an app"}.Render()
	}

	return slack.AppInfo{
		Name:     appMetadata.Name,
		Sku:      appMetadata.Sku,
		BundleId: appMetadata.BundleId,
		Id:       appMetadata.Id,
	}.Render()
}

// func handleCurrentStatusCommand(user *types.User) types.SlackResponse {
// 	appCurrentStatuses, err := getAppCurrentStatus(userAppleCredentials(user))
// 	if err != nil {
// 		return slack.EphemeralMessage{Msg: "Could not find an app"}.Render()
// 	}

// 	var slackMessages []string
// 	for _, appCurrentStatus := range appCurrentStatuses {
// 		channelMessage := fmt.Sprintf("Channel: %s\n", appCurrentStatus.Name)
// 		for _, build := range appCurrentStatus.Builds {
// 			channelMessage += fmt.Sprintf("Build: %s (%s) - %s - %s\n", build.VersionString, build.BuildNumber, build.Status, build.ReleaseDate)
// 		}
// 		slackMessages = append(slackMessages, channelMessage)
// 	}

// 	return createSlackResponse("Current Store Status for your app", slackMessages)
// }

// func handleBetaGroupsCommand(user *types.User) types.SlackResponse {
// 	betaGroups, err := getBetaGroups(userAppleCredentials(user))
// 	if err != nil {
// 		return slack.EphemeralMessage{Msg: "Could not find an app"}.Render()
// 	}

// 	var slackMessages []string
// 	for _, betaGroup := range betaGroups {
// 		groupType := "external"

// 		if betaGroup.Internal == true {
// 			groupType = "internal"
// 		}
// 		groupMessage := fmt.Sprintf("*%s* -- *%s* group with *%d* testers.", betaGroup.Name, groupType, len(betaGroup.Testers))
// 		slackMessages = append(slackMessages, groupMessage)
// 	}

// 	return createSlackResponse("Test Groups", slackMessages)
// }

// func handleInflightReleaseCommand(user *types.User) types.SlackResponse {
// 	appInfo, err := getAppMetadata(userAppleCredentials(user))
// 	if err != nil {
// 		return slack.EphemeralMessage{Msg: "Could not find your app"}.Render()
// 	}
// 	inflightRelease, err := getInflightRelease(userAppleCredentials(user))
// 	if err != nil {
// 		return slack.EphemeralMessage{Msg: "Could not find an inflight release for your app."}.Render()
// 	}

// 	phasedReleaseEnabled := "on"

// 	if inflightRelease.PhasedRelease.Id == "" {
// 		phasedReleaseEnabled = "off"
// 	}

// 	return createSlackResponse("Inflight Release",
// 		[]string{
// 			fmt.Sprintf("Next up for release is *%s (%s)* with current status `%s`.",
// 				inflightRelease.VersionName,
// 				inflightRelease.BuildNumber,
// 				inflightRelease.AppStoreState),
// 			fmt.Sprintf("The release type is `%s` and phased release is turned *%s*.", inflightRelease.ReleaseType, phasedReleaseEnabled),
// 			fmt.Sprintf("<https://appstoreconnect.apple.com/apps/%s/appstore/ios/version/inflight|App Store Connect>", appInfo.Id)})
// }

// func handleLiveReleaseCommand(user *types.User) types.SlackResponse {
// 	appInfo, err := getAppMetadata(userAppleCredentials(user))
// 	if err != nil {
// 		return slack.EphemeralMessage{Msg: "Could not find your app"}.Render()
// 	}
// 	liveRelease, err := getLiveRelease(userAppleCredentials(user))
// 	if err != nil {
// 		return slack.EphemeralMessage{Msg: "Could not find a live release for your app."}.Render()
// 	}

// 	return createSlackResponse("Live Release",
// 		[]string{
// 			fmt.Sprintf("*%s (%s)* is on day *%d* of phased release with status `%s`.",
// 				liveRelease.VersionName,
// 				liveRelease.BuildNumber,
// 				liveRelease.PhasedRelease.CurrentDayNumber,
// 				liveRelease.PhasedRelease.PhasedReleaseState),
// 			fmt.Sprintf("<https://appstoreconnect.apple.com/apps/%s/appstore/ios/version/deliverable|App Store Connect>", appInfo.Id)})
// }

// func handlePauseReleaseCommand(user *types.User) types.SlackResponse {
// 	appInfo, err := getAppMetadata(userAppleCredentials(user))
// 	if err != nil {
// 		return slack.EphemeralMessage{Msg: "Could not find your app"}.Render()
// 	}
// 	liveRelease, err := pauseLiveRelease(userAppleCredentials(user))
// 	if err != nil {
// 		return slack.EphemeralMessage{Msg: "Could not find an live release to pause."}.Render()
// 	}

// 	return createSlackResponse("Live Release",
// 		[]string{
// 			fmt.Sprintf("*%s (%s)* is on day *%d* of phased release with status `%s`.",
// 				liveRelease.VersionName,
// 				liveRelease.BuildNumber,
// 				liveRelease.PhasedRelease.CurrentDayNumber,
// 				liveRelease.PhasedRelease.PhasedReleaseState),
// 			fmt.Sprintf("<https://appstoreconnect.apple.com/apps/%s/appstore/ios/version/deliverable|App Store Connect>", appInfo.Id)})
// }

// func handleResumeReleaseCommand(user *types.User) types.SlackResponse {
// 	appInfo, err := getAppMetadata(userAppleCredentials(user))
// 	if err != nil {
// 		return slack.EphemeralMessage{Msg: "Could not find your app"}.Render()
// 	}
// 	liveRelease, err := resumeLiveRelease(userAppleCredentials(user))
// 	if err != nil {
// 		return slack.EphemeralMessage{Msg: "Could not find an paused release to resume for your app."}.Render()
// 	}

// 	return createSlackResponse("Live Release",
// 		[]string{
// 			fmt.Sprintf("*%s (%s)* is on day *%d* of phased release with status `%s`.",
// 				liveRelease.VersionName,
// 				liveRelease.BuildNumber,
// 				liveRelease.PhasedRelease.CurrentDayNumber,
// 				liveRelease.PhasedRelease.PhasedReleaseState),
// 			fmt.Sprintf("<https://appstoreconnect.apple.com/apps/%s/appstore/ios/version/deliverable|App Store Connect>", appInfo.Id)})
// }

// func handleReleaseToAllCommand(user *types.User) types.SlackResponse {
// 	appInfo, err := getAppMetadata(userAppleCredentials(user))
// 	if err != nil {
// 		return slack.EphemeralMessage{Msg: "Could not find your app"}.Render()
// 	}
// 	liveRelease, err := releaseToAll(userAppleCredentials(user))
// 	if err != nil {
// 		return slack.EphemeralMessage{Msg: "Could not find an live release to release to all for your app."}.Render()
// 	}

// 	return createSlackResponse("Live Release",
// 		[]string{
// 			fmt.Sprintf("*%s (%s)* is on day *%d* of phased release with status `%s`.",
// 				liveRelease.VersionName,
// 				liveRelease.BuildNumber,
// 				liveRelease.PhasedRelease.CurrentDayNumber,
// 				liveRelease.PhasedRelease.PhasedReleaseState),
// 			fmt.Sprintf("<https://appstoreconnect.apple.com/apps/%s/appstore/ios/version/deliverable|App Store Connect>", appInfo.Id)})
// }

// func createSlackResponse(header string, messages []string) types.SlackResponse {
// 	blocks := make([]types.SlackResponseText, len(messages)+1)
// 	blocks[0] = types.SlackResponseText{
// 		Type: "header",
// 		Text: types.SlackResponseInsideText{Type: "plain_text", Text: header},
// 	}
// 	for i, message := range messages {
// 		blocks[i+1] = types.SlackResponseText{
// 			Type: "section",
// 			Text: types.SlackResponseInsideText{Type: "mrkdwn", Text: message},
// 		}
// 	}

// 	return types.SlackResponse{Blocks: blocks, ResponseType: "in_channel"}
// }

func sendResponseToSlack(requestURL string, slackResponse types.SlackResponse) error {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(slackResponse)
	if err != nil {
		fmt.Printf("slack: could not encode the response: %s\n", err)
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

	fmt.Printf("slack: response was delivered!")
	return nil
}

func userAppleCredentials(user *types.User) *types.AppleCredentials {
	appleCredentials := types.AppleCredentials{
		BundleID: user.AppStoreBundleID.String,
		IssuerID: user.AppStoreIssuerID.String,
		KeyID:    user.AppStoreKeyID.String,
		P8File:   decrypt(user.AppStoreP8File, []byte(encryptionKey), user.AppStoreP8FileIV),
	}

	return &appleCredentials
}
