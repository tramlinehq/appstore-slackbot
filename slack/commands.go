package slack

import (
	"ciderbot/types"
	"fmt"
	"strings"
	"time"
)

const appStoreUrl string = "<https://appstoreconnect.apple.com/apps/%s/appstore|App Store Connect>"
const appStoreIcon string = "https://storage.googleapis.com/tramline-public-assets/app-store.png"

type SlackCommand interface {
	Render() types.SlackResponse
}

type AppInfo struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	BundleId string `json:"bundle_id"`
	Sku      string `json:"sku"`
}

type EphemeralMessage struct {
	Msg string `json:"msg"`
}

type HelpText struct {
	Commands map[string][2]string `json:"commands"`
}

type LiveRelease struct {
	AppId               string `json:"app_id"`
	Version             string `json:"version"`
	BuildNumber         string `json:"build_number"`
	PhasedReleaseStatus int    `json:"phased_release_status"`
	ReleaseStatus       string `json:"release_status"`
}

type CurrentStoreStatus struct {
	Channels []struct {
		Name   string `json:"name"`
		Builds []struct {
			Id            string    `json:"id"`
			BuildNumber   string    `json:"build_number"`
			Status        string    `json:"status"`
			VersionString string    `json:"version_string"`
			ReleaseDate   time.Time `json:"release_date"`
		} `json:"builds"`
	}
}

type BetaGroups struct {
	Groups []struct {
		Name        string `json:"name"`
		Internal    bool   `json:"internal"`
		TesterCount int    `json:"testers"`
	}
}

type InflightRelease struct {
	VersionString string `json:"version_string"`
	BuildNumber   string `json:"build_number"`
	StoreStatus   string `json:"store_status"`
	ReleaseType   string `json:"release_type"`
	PhasedRelease bool   `json:"phased_release"`
	AppId         string `json:"app_id"`
}

func (data InflightRelease) Render() types.SlackResponse {
	phasedRelease := ""

	if data.PhasedRelease {
		phasedRelease = "on"
	} else {
		phasedRelease = "off"
	}

	line1 := fmt.Sprintf("The upcoming release in progress is *%s (%s)* with the current status of `%s`.", data.VersionString, data.BuildNumber, data.StoreStatus)
	line2 := fmt.Sprintf("The release type is `%s` and phased release is turned *%s*.", data.ReleaseType, phasedRelease)

	slackResponse := types.SlackResponse{
		ResponseType: "in_channel",
		Blocks: []types.Block{
			{
				Type: "header",
				Text: &types.Text{
					Type:  "plain_text",
					Text:  ":airplane_departure: Inflight Release",
					Emoji: true,
				},
			},
			{
				Type: "divider",
			},
			{
				Type: "section",
				Text: &types.Text{
					Type: "mrkdwn",
					Text: line1,
				},
			},
			{
				Type: "section",
				Text: &types.Text{
					Type: "mrkdwn",
					Text: line2,
				},
			},
			{
				Type: "divider",
			},
			{
				Type: "context",
				Elements: []types.Element{
					{
						Type:     "image",
						ImageURL: appStoreIcon,
						AltText:  "app store connect",
					},
					{
						Type: "mrkdwn",
						Text: fmt.Sprintf(appStoreUrl, data.AppId),
					},
				},
			},
			{
				Type: "divider",
			},
		},
	}

	return slackResponse
}

func (data BetaGroups) Render() types.SlackResponse {
	var slackBlocks []types.Block

	slackBlocks = append(slackBlocks, types.Block{
		Type: "header",
		Text: &types.Text{
			Type:  "plain_text",
			Text:  ":test_tube: Beta Groups",
			Emoji: true,
		},
	})

	slackBlocks = append(slackBlocks, types.Block{
		Type: "divider",
	})

	for _, group := range data.Groups {
		var groupType string
		if group.Internal {
			groupType = "*internal* group"
		} else {
			groupType = "*external* group"
		}

		text := fmt.Sprintf("❖ *%s* is an %s with %d tester", group.Name, groupType, group.TesterCount)
		if group.TesterCount > 1 || group.TesterCount == 0 {
			text += "s"
		}

		text += "."

		slackBlocks = append(slackBlocks, types.Block{
			Type: "section",
			Text: &types.Text{
				Type: "mrkdwn",
				Text: text,
			},
		})
	}

	slackBlocks = append(slackBlocks, types.Block{
		Type: "divider",
	})

	return types.SlackResponse{Blocks: slackBlocks, ResponseType: "in_channel"}
}

func (data CurrentStoreStatus) Render() types.SlackResponse {
	slackBlocks := []types.Block{
		{
			Type: "header",
			Text: &types.Text{
				Type:  "plain_text",
				Text:  ":convenience_store: Current Store Status",
				Emoji: true,
			},
		},
		{
			Type: "divider",
		},
	}

	for _, channel := range data.Channels {
		slackBlocks = append(slackBlocks,
			types.Block{
				Type: "section",
				Text: &types.Text{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*%s*", strings.Title(channel.Name)),
				},
			},
		)
		for _, build := range channel.Builds {
			slackBlocks = append(slackBlocks,
				types.Block{
					Type: "context",
					Elements: []types.Element{
						{
							Type: "mrkdwn",
							Text: fmt.Sprintf("*%s (%s)* was `%s` on *%s*", build.VersionString, build.BuildNumber, build.Status, build.ReleaseDate.Format("Monday, Jan 15th 15:04:05, 2006")),
						},
					},
				},
			)
		}
	}

	slackBlocks = append(slackBlocks, types.Block{
		Type: "divider",
	})

	return types.SlackResponse{Blocks: slackBlocks, ResponseType: "in_channel"}
}

func (data LiveRelease) Render() types.SlackResponse {
	line1 := fmt.Sprintf("We're on *day %d* of *phased release* with status `%s`.", data.PhasedReleaseStatus, data.ReleaseStatus)

	if data.ReleaseStatus == "COMPLETE" {
		line1 = fmt.Sprintf("The release was fully rolled out to *all users* after *day %d* of the *phased rollout*.", data.PhasedReleaseStatus)
	} else if data.ReleaseStatus == "" {
		line1 = "The release was fully rolled out to *all users* without phased release."
	}

	slackResponse := types.SlackResponse{
		ResponseType: "in_channel",
		Blocks: []types.Block{
			{
				Type: "header",
				Text: &types.Text{
					Type:  "plain_text",
					Text:  ":iphone: Live Release",
					Emoji: true,
				},
			},
			{
				Type: "divider",
			},
			{
				Type: "section",
				Fields: []types.Text{
					{
						Type: "mrkdwn",
						Text: fmt.Sprintf("*Version:* %s :package:", data.Version),
					},
					{
						Type: "mrkdwn",
						Text: fmt.Sprintf("*Build Number:* %s :1234:", data.BuildNumber),
					},
				},
			},
			{
				Type: "section",
				Text: &types.Text{
					Type: "mrkdwn",
					Text: line1,
				},
			},
			{
				Type: "divider",
			},
			{
				Type: "context",
				Elements: []types.Element{
					{
						Type:     "image",
						ImageURL: appStoreIcon,
						AltText:  "app store connect",
					},
					{
						Type: "mrkdwn",
						Text: fmt.Sprintf(appStoreUrl, data.AppId),
					},
				},
			},
			{
				Type: "divider",
			},
		},
	}

	return slackResponse
}

func (data EphemeralMessage) Render() types.SlackResponse {
	slackResponse := types.SlackResponse{
		ResponseType: "ephemeral",
		Blocks: []types.Block{
			{
				Type: "section",
				Text: &types.Text{
					Type: "mrkdwn",
					Text: data.Msg,
				},
			},
		},
	}

	return slackResponse
}

func (data HelpText) Render() types.SlackResponse {
	slackResponse := types.SlackResponse{
		ResponseType: "in_channel",
		Blocks: []types.Block{
			{
				Type: "header",
				Text: &types.Text{
					Type:  "plain_text",
					Text:  ":books: Usage Guide",
					Emoji: true,
				},
			},
			{
				Type: "divider",
			},
		},
	}

	for name, desc := range data.Commands {
		section := types.Block{
			Type: "section",
			Text: &types.Text{
				Type: "mrkdwn",
				Text: fmt.Sprintf("%s `%s`\n%s", desc[0], name, desc[1]),
			},
		}

		slackResponse.Blocks = append(slackResponse.Blocks, section)
	}

	slackResponse.Blocks = append(slackResponse.Blocks, types.Block{
		Type: "divider",
	})

	return slackResponse
}

func (data AppInfo) Render() types.SlackResponse {
	slackResponse := types.SlackResponse{
		ResponseType: "in_channel",
		Blocks: []types.Block{
			{
				Type: "header",
				Text: &types.Text{
					Type:  "plain_text",
					Text:  ":information_source: App Info",
					Emoji: true,
				},
			},
			{
				Type: "divider",
			},
			{
				Type: "section",
				Fields: []types.Text{
					{
						Type: "mrkdwn",
						Text: fmt.Sprintf("*Name:* %s :phone:", data.Name),
					},
					{
						Type: "mrkdwn",
						Text: fmt.Sprintf("*SKU:* %s :package:", data.Sku),
					},
					{
						Type: "mrkdwn",
						Text: fmt.Sprintf("*Bundle ID:* %s :file_folder:", data.BundleId),
					},
					{
						Type: "mrkdwn",
						Text: fmt.Sprintf("*ID:* %s :id:", data.Id),
					},
				},
			},
			{
				Type: "divider",
			},
		},
	}

	return slackResponse
}
