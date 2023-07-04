package main

import (
	"bytes"
	slack "ciderbot/slack"
	"ciderbot/types"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
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
	case "live_release":
		return handleLiveReleaseCommand(user)
	case "overall_status":
		return handleCurrentStatusCommand(user)
	case "beta_groups":
		return handleBetaGroupsCommand(user)
	case "inflight_release":
		return handleInflightReleaseCommand(user)
	case "pause_live_release":
		return handlePauseReleaseCommand(user)
	case "resume_live_release":
		return handleResumeReleaseCommand(user)
	case "release_to_all":
		return handleReleaseToAllCommand(user)
	default:
		return slack.EphemeralMessage{Msg: "Please input a valid command!"}.Render()
	}
}

func handleHelpCommand(_user *types.User) types.SlackResponse {
	return slack.HelpText{Commands: ValidSlackCommands}.Render()
}

func handleInfoCommand(user *types.User) types.SlackResponse {
	appMetadata, err := getAppMetadata(userAppleCredentials(user))
	if err != nil {
		return slack.EphemeralMessage{Msg: "Could not find an app."}.Render()
	}

	return slack.AppInfo{
		Name:     appMetadata.Name,
		Sku:      appMetadata.Sku,
		BundleId: appMetadata.BundleId,
		Id:       appMetadata.Id,
	}.Render()
}

func handleLiveReleaseCommand(user *types.User) types.SlackResponse {
	appInfo, err := getAppMetadata(userAppleCredentials(user))
	if err != nil {
		return slack.EphemeralMessage{Msg: "Could not find your app."}.Render()
	}

	liveRelease, err := getLiveRelease(userAppleCredentials(user))

	if err != nil {
		return slack.EphemeralMessage{Msg: "Could not find a live release for your app."}.Render()
	}

	return slack.LiveRelease{
		AppId:               appInfo.Id,
		Version:             liveRelease.VersionName,
		BuildNumber:         liveRelease.BuildNumber,
		PhasedReleaseStatus: liveRelease.PhasedRelease.CurrentDayNumber,
		ReleaseStatus:       liveRelease.PhasedRelease.PhasedReleaseState,
	}.Render()
}

func handleCurrentStatusCommand(user *types.User) types.SlackResponse {
	appCurrentStatuses, err := getAppCurrentStatus(userAppleCredentials(user))
	if err != nil {
		return slack.EphemeralMessage{Msg: "Could not find an app."}.Render()
	}

	storeStatus := slack.CurrentStoreStatus{}

	for _, channelStatus := range appCurrentStatuses {
		channel := struct {
			Name   string `json:"name"`
			Builds []struct {
				Id            string    `json:"id"`
				BuildNumber   string    `json:"build_number"`
				Status        string    `json:"status"`
				VersionString string    `json:"version_string"`
				ReleaseDate   time.Time `json:"release_date"`
			} `json:"builds"`
		}{
			Name:   channelStatus.Name,
			Builds: channelStatus.Builds,
		}

		storeStatus.Channels = append(storeStatus.Channels, channel)
	}

	return storeStatus.Render()
}

func handleBetaGroupsCommand(user *types.User) types.SlackResponse {
	betaGroups, err := getBetaGroups(userAppleCredentials(user))
	if err != nil {
		return slack.EphemeralMessage{Msg: "Could not find an app."}.Render()
	}

	groups := slack.BetaGroups{}

	for _, betaGroup := range betaGroups {
		group := struct {
			Name        string `json:"name"`
			Internal    bool   `json:"internal"`
			TesterCount int    `json:"testers"`
		}{
			Name:        betaGroup.Name,
			Internal:    betaGroup.Internal,
			TesterCount: len(betaGroup.Testers),
		}

		groups.Groups = append(groups.Groups, group)
	}

	return groups.Render()
}

func handleInflightReleaseCommand(user *types.User) types.SlackResponse {
	appInfo, err := getAppMetadata(userAppleCredentials(user))
	if err != nil {
		return slack.EphemeralMessage{Msg: "Could not find your app."}.Render()
	}

	inflightRelease, err := getInflightRelease(userAppleCredentials(user))
	if err != nil {
		return slack.EphemeralMessage{Msg: "Could not find an inflight release for your app."}.Render()
	}

	phasedReleaseEnabled := true

	if inflightRelease.PhasedRelease.Id == "" {
		phasedReleaseEnabled = false
	}

	return slack.InflightRelease{
		VersionString: inflightRelease.VersionName,
		BuildNumber:   inflightRelease.BuildNumber,
		StoreStatus:   inflightRelease.AppStoreState,
		ReleaseType:   inflightRelease.ReleaseType,
		PhasedRelease: phasedReleaseEnabled,
		AppId:         appInfo.Id,
	}.Render()
}

func handlePauseReleaseCommand(user *types.User) types.SlackResponse {
	appInfo, err := getAppMetadata(userAppleCredentials(user))
	if err != nil {
		return slack.EphemeralMessage{Msg: "Could not find your app."}.Render()
	}
	liveRelease, err := pauseLiveRelease(userAppleCredentials(user))
	if err != nil {
		return slack.EphemeralMessage{Msg: "Could not find a live release to pause."}.Render()
	}

	return slack.LiveRelease{
		AppId:               appInfo.Id,
		Version:             liveRelease.VersionName,
		BuildNumber:         liveRelease.BuildNumber,
		PhasedReleaseStatus: liveRelease.PhasedRelease.CurrentDayNumber,
		ReleaseStatus:       liveRelease.PhasedRelease.PhasedReleaseState,
	}.Render()
}

func handleResumeReleaseCommand(user *types.User) types.SlackResponse {
	appInfo, err := getAppMetadata(userAppleCredentials(user))
	if err != nil {
		return slack.EphemeralMessage{Msg: "Could not find your app."}.Render()
	}
	liveRelease, err := resumeLiveRelease(userAppleCredentials(user))
	if err != nil {
		return slack.EphemeralMessage{Msg: "Could not find a paused release to resume."}.Render()
	}

	return slack.LiveRelease{
		AppId:               appInfo.Id,
		Version:             liveRelease.VersionName,
		BuildNumber:         liveRelease.BuildNumber,
		PhasedReleaseStatus: liveRelease.PhasedRelease.CurrentDayNumber,
		ReleaseStatus:       liveRelease.PhasedRelease.PhasedReleaseState,
	}.Render()
}

func handleReleaseToAllCommand(user *types.User) types.SlackResponse {
	appInfo, err := getAppMetadata(userAppleCredentials(user))
	if err != nil {
		return slack.EphemeralMessage{Msg: "Could not find your app."}.Render()
	}
	liveRelease, err := releaseToAll(userAppleCredentials(user))
	if err != nil {
		return slack.EphemeralMessage{Msg: "Could not find a live release."}.Render()
	}

	return slack.LiveRelease{
		AppId:               appInfo.Id,
		Version:             liveRelease.VersionName,
		BuildNumber:         liveRelease.BuildNumber,
		PhasedReleaseStatus: liveRelease.PhasedRelease.CurrentDayNumber,
		ReleaseStatus:       liveRelease.PhasedRelease.PhasedReleaseState,
	}.Render()
}

func sendResponseToSlack(requestURL string, slackResponse types.SlackResponse) error {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(slackResponse)
	if err != nil {
		fmt.Printf("slack: could not encode the response: %s\n", err)
		return err
	}

	fmt.Println(body.String())

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

	fmt.Printf("slack: response was delivered to!")
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
