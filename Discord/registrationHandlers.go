package Discord

import (
	"Crypto-Trading-Discord-Bot/Mongo"
	"context"
	dgo "github.com/bwmarrin/discordgo"
	"log"
)

func getRegistrationHandlers() *map[string]BotCommand {
	return &map[string]BotCommand{
		"register-me-as": {
			Info: &dgo.ApplicationCommand{
				Name:        "register-me-as",
				Description: "Registers you in the system. If already registered, only the nickname is updated.",
				Options: []*dgo.ApplicationCommandOption{
					{
						Type:        dgo.ApplicationCommandOptionString,
						Name:        "nick-name",
						Description: "What should I call you? By default, I will use your username to refer to you.",
						Required:    false,
					},
				},
			},
			Handler: func(bot *Bot, s *dgo.Session, i *dgo.InteractionCreate) {
				var user *dgo.User
				if i.User != nil {
					user = i.User // DM
				} else {
					user = i.Member.User // Server
				}

				optionMap := map[string]*dgo.ApplicationCommandInteractionDataOption{}
				for _, option := range i.ApplicationCommandData().Options {
					optionMap[option.Name] = option
				}

				var nickName string
				if option := optionMap["nick-name"]; option != nil {
					nickName = option.StringValue()
				} else {
					nickName = user.Username
				}

				res := bot.Db().GetCollection("RegisteredUsers").UpsertOne(
					context.TODO(),
					Mongo.RegisteredUser{UId: user.ID},
					Mongo.RegisteredUser{UId: user.ID, NickName: nickName},
				)

				var msg string
				if res.Success {
					if res.WasInserted {
						msg = "Welcome " + nickName + ".\nYou are now registered with the system."
					} else if res.WasUpdated {
						msg = "I shall now refer to you as " + nickName + "."
					} else {
						msg = "You are already registered in the system."
					}
				} else {
					msg = "Registration Failed. Please try again later."
					log.Println("DiscordRegistrationHandlers getRegistrationHandlers: Registration Failed:", res.Error)
				}

				if err := s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
					Type: dgo.InteractionResponseChannelMessageWithSource,
					Data: &dgo.InteractionResponseData{
						Content: msg,
					},
				}); err != nil {
					log.Println("DiscordRegistrationHandlers getRegistrationHandlers: Error when interacting with input.", err)
				}
			},
		},
	}
}

func getRegistrationComponentHandlers() *map[string]func(bot *Bot, s *dgo.Session, i *dgo.InteractionCreate) {
	return &map[string]func(bot *Bot, s *dgo.Session, i *dgo.InteractionCreate){}
}
