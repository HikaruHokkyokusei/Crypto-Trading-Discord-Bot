package Discord

import (
	dgo "github.com/bwmarrin/discordgo"
	"log"
)

func GetGeneralHandlers() *map[string]BotCommand {
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
				var uid, msg string

				if i.User != nil {
					uid = i.User.ID // DM
				} else {
					uid = i.Member.User.ID // Server
				}

				options, size := i.ApplicationCommandData().Options, len(i.ApplicationCommandData().Options)
				if size > 0 {
					msg = options[0].StringValue()
				}

				if err := s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
					Type: dgo.InteractionResponseChannelMessageWithSource,
					Data: &dgo.InteractionResponseData{
						Content: uid + " : " + msg,
					},
				}); err != nil {
					log.Println("DiscordGeneralHandlers GetGeneralHandlers: Error when interacting with input.", err)
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
					log.Println("DiscordGeneralHandlers GetGeneralHandlers: Error when responding to application command", err)
				}
			},
		},
	}
}

func GetGeneralComponentHandlers() *map[string]func(bot *Bot, s *dgo.Session, i *dgo.InteractionCreate) {
	return &map[string]func(bot *Bot, s *dgo.Session, i *dgo.InteractionCreate){
		"like-yes": func(bot *Bot, s *dgo.Session, i *dgo.InteractionCreate) {
			if err := s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
				Type: dgo.InteractionResponseUpdateMessage,
				Data: &dgo.InteractionResponseData{
					Content: "🙂",
					Flags:   dgo.MessageFlagsEphemeral,
				},
			}); err != nil {
				log.Println("DiscordGeneralHandlers GetGeneralComponentHandlers: Error when responding to component interaction", err)
				return
			}

			//if err := s.ChannelMessageDelete(i.ChannelID, i.Message.ID); err != nil {
			//	log.Println("DiscordGeneralHandlers GetGeneralComponentHandlers: Error when deleting message", err)
			//}
		},
		"like-no": func(bot *Bot, s *dgo.Session, i *dgo.InteractionCreate) {
			if err := s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
				Type: dgo.InteractionResponseUpdateMessage,
				Data: &dgo.InteractionResponseData{
					Content: "😭",
					Flags:   dgo.MessageFlagsEphemeral,
				},
			}); err != nil {
				log.Println("DiscordGeneralHandlers GetGeneralComponentHandlers: Error when responding to component interaction", err)
				return
			}

			if err := s.ChannelMessageDelete(i.Interaction.ChannelID, i.Interaction.Message.ID); err != nil {
				log.Println("DiscordGeneralHandlers GetGeneralComponentHandlers: Error when deleting message", err)
			}
		},
	}
}