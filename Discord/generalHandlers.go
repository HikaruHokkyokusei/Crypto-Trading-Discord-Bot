package Discord

import (
	"Crypto-Trading-Discord-Bot/utils"
	dgo "github.com/bwmarrin/discordgo"
	"log"
)

func getGeneralHandlers() *map[string]BotCommand {
	return &map[string]BotCommand{
		"echo": {
			Info: &dgo.ApplicationCommand{
				Name:        "echo",
				Description: "Echos back your message along with your UID.",
				Options: []*dgo.ApplicationCommandOption{
					{
						Type:        dgo.ApplicationCommandOptionString,
						Name:        "message",
						Description: "Message to echo back",
						Required:    true,
					},
				},
			},
			Handler: func(bot *Bot, s *dgo.Session, i *dgo.InteractionCreate) {
				user, optionMap := utils.GetUserAndOptionMap(i)

				var msg string
				if option := optionMap["message"]; option != nil {
					msg = option.StringValue()
				}

				if err := s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
					Type: dgo.InteractionResponseChannelMessageWithSource,
					Data: &dgo.InteractionResponseData{
						Content: user.ID + " : " + msg,
					},
				}); err != nil {
					log.Println("DiscordGeneralHandlers getGeneralHandlers: Error when interacting with input.", err)
				}
			},
		},
		"like": {
			Info: &dgo.ApplicationCommand{
				Name:        "like",
				Description: "Check whether you like me?",
			},
			Handler: func(bot *Bot, s *dgo.Session, i *dgo.InteractionCreate) {
				if err := s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
					Type: dgo.InteractionResponseChannelMessageWithSource,
					Data: &dgo.InteractionResponseData{
						Content: "Do you like me?",
						Flags:   dgo.MessageFlagsEphemeral,
						Components: []dgo.MessageComponent{
							dgo.ActionsRow{
								Components: []dgo.MessageComponent{
									dgo.Button{
										Emoji: dgo.ComponentEmoji{
											Name: "✔️",
										},
										Label:    "Yep",
										CustomID: "like-yes",
									},
									dgo.Button{
										Emoji: dgo.ComponentEmoji{
											Name: "❌",
										},
										Label:    "Nope",
										CustomID: "like-no",
									},
								},
							},
						},
					},
				}); err != nil {
					log.Println("DiscordGeneralHandlers getGeneralHandlers: Error when responding to application command", err)
				}
			},
		},
	}
}

func getGeneralComponentHandlers() *map[string]func(bot *Bot, s *dgo.Session, i *dgo.InteractionCreate) {
	return &map[string]func(bot *Bot, s *dgo.Session, i *dgo.InteractionCreate){
		"like-yes": func(bot *Bot, s *dgo.Session, i *dgo.InteractionCreate) {
			if err := s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
				Type: dgo.InteractionResponseUpdateMessage,
				Data: &dgo.InteractionResponseData{
					Content: "🙂",
					Flags:   dgo.MessageFlagsEphemeral,
				},
			}); err != nil {
				log.Println("DiscordGeneralHandlers getGeneralComponentHandlers: Error when responding to component interaction", err)
				return
			}
		},
		"like-no": func(bot *Bot, s *dgo.Session, i *dgo.InteractionCreate) {
			if err := s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
				Type: dgo.InteractionResponseUpdateMessage,
				Data: &dgo.InteractionResponseData{
					Content: "😭",
					Flags:   dgo.MessageFlagsEphemeral,
				},
			}); err != nil {
				log.Println("DiscordGeneralHandlers getGeneralComponentHandlers: Error when responding to component interaction", err)
				return
			}

			if err := s.ChannelMessageDelete(i.Interaction.ChannelID, i.Interaction.Message.ID); err != nil {
				log.Println("DiscordGeneralHandlers getGeneralComponentHandlers: Error when deleting message", err)
			}
		},
	}
}
