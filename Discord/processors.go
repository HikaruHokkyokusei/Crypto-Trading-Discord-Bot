package Discord

import (
	dgo "github.com/bwmarrin/discordgo"
	"log"
)

type BotCommand struct {
	info    *dgo.ApplicationCommand
	handler func(bot *Bot, s *dgo.Session, i *dgo.InteractionCreate)
}

func (botCommand BotCommand) Info() dgo.ApplicationCommand {
	return *botCommand.info
}

var BotCommands = map[string]BotCommand{
	"echo": {
		info: &dgo.ApplicationCommand{
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
		handler: func(bot *Bot, s *dgo.Session, i *dgo.InteractionCreate) {
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
				log.Println("DiscordProcessors BotCommandHandlers: Error when interacting with input.", err)
			}
		},
	},
	"like": {
		info: &dgo.ApplicationCommand{
			Name:        "like",
			Description: "Check whether you like me?",
		},
		handler: func(bot *Bot, s *dgo.Session, i *dgo.InteractionCreate) {
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
				log.Println("DiscordProcessors BotCommandHandlers: Error when responding to application command", err)
			}
		},
	},
}

var componentHandlers = map[string]func(bot *Bot, s *dgo.Session, i *dgo.InteractionCreate){
	"like-yes": func(bot *Bot, s *dgo.Session, i *dgo.InteractionCreate) {
		if err := s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
			Type: dgo.InteractionResponseChannelMessageWithSource,
			Data: &dgo.InteractionResponseData{
				Content: "🙂",
				Flags:   dgo.MessageFlagsEphemeral,
			},
		}); err != nil {
			log.Println("DiscordProcessors BotCommandHandlers: Error when responding to component interaction", err)
		}
	},
	"like-no": func(bot *Bot, s *dgo.Session, i *dgo.InteractionCreate) {
		if err := s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
			Type: dgo.InteractionResponseChannelMessageWithSource,
			Data: &dgo.InteractionResponseData{
				Content: "😭",
				Flags:   dgo.MessageFlagsEphemeral,
			},
		}); err != nil {
			log.Println("DiscordProcessors BotCommandHandlers: Error when responding to component interaction", err)
		}
	},
}

func GetBotHandlers(bot Bot) *[]interface{} {
	return &[]interface{}{
		func(s *dgo.Session, r *dgo.Ready) {
			log.Println("DiscordInit StartSession: Session Started. Logged in as: `" + s.State.User.Username + "#" + s.State.User.Discriminator + "`")
		},
		func(s *dgo.Session, m *dgo.MessageCreate) {
			if m.Author.ID != s.State.User.ID {
				if _, err := s.ChannelMessageSend(m.ChannelID, m.Author.ID+" : "+m.Content); err != nil {
					log.Println("DiscordProcessors BotHandlers: Error when sending message:", err)
				}
			}
		},
		func(s *dgo.Session, i *dgo.InteractionCreate) {
			var handler func(bot *Bot, s *dgo.Session, i *dgo.InteractionCreate)

			switch i.Type {
			case dgo.InteractionApplicationCommand:
				botCommand, ok := BotCommands[i.ApplicationCommandData().Name]
				if ok {
					handler = botCommand.handler
				} else {
					log.Println("DiscordProcessors BotHandlers: No handler for interaction command", i.ApplicationCommandData().Name)
					return
				}
			case dgo.InteractionMessageComponent:
				fun, ok := componentHandlers[i.MessageComponentData().CustomID]
				if ok {
					handler = fun
				} else {
					log.Println("DiscordProcessors BotHandlers: No handler for message command", i.MessageComponentData().CustomID)
					return
				}
			}

			handler(&bot, s, i)
		},
	}
}
