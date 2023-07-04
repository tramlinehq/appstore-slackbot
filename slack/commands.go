package slack

import (
	"ciderbot/types"
	"fmt"
)

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

func (data EphemeralMessage) Render() types.SlackResponse {
	slackResponse := types.SlackResponse{
		ResponseType: "ephemeral",
		Blocks: []types.Block{
			{
				Type: "section",
				Text: types.Text{
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
				Type: "section",
				Text: types.Text{
					Type: "mrkdwn",
					Text: "*Usage Guide*",
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
			Text: types.Text{
				Type: "mrkdwn",
				Text: fmt.Sprintf("%s *%s*\n%s", desc[0], name, desc[1]),
			},
		}

		slackResponse.Blocks = append(slackResponse.Blocks, section)
	}

	return slackResponse
}

func (data AppInfo) Render() types.SlackResponse {
	slackResponse := types.SlackResponse{
		ResponseType: "in_channel",
		Blocks: []types.Block{
			{
				Type: "section",
				Text: types.Text{
					Type: "mrkdwn",
					Text: ":information_source: *App Info*",
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
		},
	}

	return slackResponse
}
