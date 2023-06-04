package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var ValidSlackCommands = map[string]*regexp.Regexp{
	"app_info":            regexp.MustCompile("app_info"),
	"overall_status":      regexp.MustCompile("overall_status"),
	"beta_groups":         regexp.MustCompile("beta_groups"),
	"live_release":        regexp.MustCompile("live_release"),
	"pause_live_release":  regexp.MustCompile("pause_live_release"),
	"resume_live_release": regexp.MustCompile("resume_live_release"),
	"release_to_all":      regexp.MustCompile("release_to_all"),
}

func handleSlackCommand(form SlackFormData, user *User) SlackResponse {
	fmt.Println("command", form.Command)
	fmt.Println("text", form.Text)

	command := strings.Split(form.Text, " ")[0]
	if commandPattern, ok := ValidSlackCommands[command]; ok == true {
		go handleValidSlackCommand(commandPattern, form.Text, form.ResponseUrl, user)
		return createSlackResponse([]string{"Got it, working on it."}, "ephemeral")
	}

	return createSlackResponse([]string{"Please input a valid command"}, "ephemeral")

}

func handleValidSlackCommand(commandPattern *regexp.Regexp, command string, responseURL string, user *User) {
	slackResponse := processValidSlackCommand(commandPattern, command, user)
	sendResponseToSlack(responseURL, slackResponse)
}

func processValidSlackCommand(commandPattern *regexp.Regexp, command string, user *User) SlackResponse {
	matched := commandPattern.FindAllString(command, -1)
	switch matched[0] {
	case "app_info":
		return handleInfoCommand(user)
	case "overall_status":
		return handleCurrentStatusCommand(user)
	case "beta_groups":
		return handleBetaGroupsCommand(user)
	case "live_release":
		return handleLiveReleaseCommand(user)
	case "pause_live_release":
		return handlePauseReleaseCommand(user)
	case "resume_live_release":
		return handleResumeReleaseCommand(user)
	case "release_to_all":
		return handleReleaseToAllCommand(user)
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

func handleBetaGroupsCommand(user *User) SlackResponse {
	betaGroups, err := getBetaGroups(userAppleCredentials(user))
	if err != nil {
		return createSlackResponse([]string{"Could not find an app."}, "ephemeral")
	}

	var slackMessages []string
	for _, betaGroup := range betaGroups {
		groupMessage := fmt.Sprintf("Group: %s\n", betaGroup.Name)
		groupMessage += fmt.Sprintf("Internal: %t\n", betaGroup.Internal)
		groupMessage += fmt.Sprintf("Testers: %d\n", len(betaGroup.Testers))
		groupMessage += fmt.Sprintf("-----------")
		slackMessages = append(slackMessages, groupMessage)
	}

	return createSlackResponse(slackMessages, "in_channel")
}

func handleLiveReleaseCommand(user *User) SlackResponse {
	liveRelease, err := getLiveRelease(userAppleCredentials(user))
	if err != nil {
		return createSlackResponse([]string{"Could not find an app."}, "ephemeral")
	}

	return createSlackResponse([]string{fmt.Sprintf(`Live Release:
Version: %s
Build Number: %s
Store Status: %s
Phased Release Status: %s
Phased Release Day: %d`,
		liveRelease.VersionName,
		liveRelease.BuildNumber,
		liveRelease.AppStoreState,
		liveRelease.PhasedRelease.PhasedReleaseState,
		liveRelease.PhasedRelease.CurrentDayNumber)}, "in_channel")
}

func handlePauseReleaseCommand(user *User) SlackResponse {
	liveRelease, err := pauseLiveRelease(userAppleCredentials(user))
	if err != nil {
		return createSlackResponse([]string{"Could not find an live release to pause."}, "ephemeral")
	}

	return createSlackResponse([]string{fmt.Sprintf(`Live Release:
Version: %s
Build Number: %s
Store Status: %s
Phased Release Status: %s
Phased Release Day: %d`,
		liveRelease.VersionName,
		liveRelease.BuildNumber,
		liveRelease.AppStoreState,
		liveRelease.PhasedRelease.PhasedReleaseState,
		liveRelease.PhasedRelease.CurrentDayNumber)}, "in_channel")
}

func handleResumeReleaseCommand(user *User) SlackResponse {
	liveRelease, err := resumeLiveRelease(userAppleCredentials(user))
	if err != nil {
		return createSlackResponse([]string{"Could not find an paused release to resume."}, "ephemeral")
	}

	return createSlackResponse([]string{fmt.Sprintf(`Live Release:
Version: %s
Build Number: %s
Store Status: %s
Phased Release Status: %s
Phased Release Day: %d`,
		liveRelease.VersionName,
		liveRelease.BuildNumber,
		liveRelease.AppStoreState,
		liveRelease.PhasedRelease.PhasedReleaseState,
		liveRelease.PhasedRelease.CurrentDayNumber)}, "in_channel")
}

func handleReleaseToAllCommand(user *User) SlackResponse {
	liveRelease, err := releaseToAll(userAppleCredentials(user))
	if err != nil {
		return createSlackResponse([]string{"Could not find an live release to release to all."}, "ephemeral")
	}

	return createSlackResponse([]string{fmt.Sprintf(`Live Release:
Version: %s
Build Number: %s
Store Status: %s
Phased Release Status: %s
Phased Release Day: %d`,
		liveRelease.VersionName,
		liveRelease.BuildNumber,
		liveRelease.AppStoreState,
		liveRelease.PhasedRelease.PhasedReleaseState,
		liveRelease.PhasedRelease.CurrentDayNumber)}, "in_channel")
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
