package Discord

import (
	"Crypto-Trading-Discord-Bot/utils"
	"fmt"
	dgo "github.com/bwmarrin/discordgo"
	"log"
)

func getOwnerHandlers() *map[string]BotCommand {
	return &map[string]BotCommand{
		"delete-command": {
			Info: &dgo.ApplicationCommand{
				Name:        "delete-command",
				Description: "Registers you in the system. If already registered, only the nickname is updated.",
				Options: []*dgo.ApplicationCommandOption{
					{
						Type:        dgo.ApplicationCommandOptionString,
						Name:        "command-name",
						Description: "Name of the command to delete.",
						Required:    true,
					},
					{
						Type:        dgo.ApplicationCommandOptionString,
						Name:        "guild-id",
						Description: "Any specific guildId to delete from",
						Required:    false,
					},
				},
			},
			Handler: func(bot *Bot, s *dgo.Session, i *dgo.InteractionCreate) {
				user, _ := utils.GetUserAndOptionMap(i)
				var msg string

				if user.ID == bot.ownerId {
					optionMap := map[string]*dgo.ApplicationCommandInteractionDataOption{}
					for _, option := range i.ApplicationCommandData().Options {
						optionMap[option.Name] = option
					}

					gId, cName := "", optionMap["command-name"].StringValue()

					if v, ok := optionMap["guild-id"]; ok {
						gId = v.StringValue()
					}

					var command *dgo.ApplicationCommand = nil
					if registeredCommands, _ := bot.session.ApplicationCommands(bot.session.State.User.ID, gId); registeredCommands != nil {
						for _, registeredCommand := range registeredCommands {
							if registeredCommand.Name == cName {
								command = registeredCommand
							}
						}
					}
					if command != nil {
						if err := bot.session.ApplicationCommandDelete(bot.session.State.User.ID, gId, command.ID); err == nil {
							msg = "Command successfully deleted."
						} else {
							msg = fmt.Sprint("Error when deleting the command:", err)
						}
					} else {
						msg = "No such command registered."
					}
				} else {
					msg = "You are not allowed to use this command."
				}

				if err := s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
					Type: dgo.InteractionResponseChannelMessageWithSource,
					Data: &dgo.InteractionResponseData{
						Content: msg,
					},
				}); err != nil {
					log.Println("DiscordOwnerHandlers getOwnerHandlers: Error when interacting with input.", err)
				}
			},
		},
	}
}

func getOwnerComponentHandlers() *map[string]func(bot *Bot, s *dgo.Session, i *dgo.InteractionCreate) {
	return &map[string]func(bot *Bot, s *dgo.Session, i *dgo.InteractionCreate){}
}
