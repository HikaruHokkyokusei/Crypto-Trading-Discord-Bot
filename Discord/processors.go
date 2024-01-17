package Discord

import (
	dgo "github.com/bwmarrin/discordgo"
	"log"
)

type BotCommand struct {
	Info    *dgo.ApplicationCommand
	Handler func(bot *Bot, s *dgo.Session, i *dgo.InteractionCreate)
	GuildId string
}

var (
	botCommands       = map[string]BotCommand{}
	componentHandlers = map[string]func(bot *Bot, s *dgo.Session, i *dgo.InteractionCreate){}
)

func gatherCommandsAndHandlers(bot *Bot) {
	for name, botCommand := range *getGeneralHandlers() {
		botCommands[name] = botCommand
	}
	for name, handler := range *getGeneralComponentHandlers() {
		componentHandlers[name] = handler
	}
	for name, botCommand := range *getOwnerHandlers() {
		botCommands[name] = botCommand
	}
	for name, handler := range *getOwnerComponentHandlers() {
		componentHandlers[name] = handler
	}
	for name, botCommand := range *getRegistrationHandlers() {
		botCommands[name] = botCommand
	}
	for name, handler := range *getRegistrationComponentHandlers() {
		componentHandlers[name] = handler
	}
}

func getBotCommandsAndHandlers(bot *Bot) (*map[string]BotCommand, *[]interface{}) {
	gatherCommandsAndHandlers(bot)
	return &botCommands, &[]interface{}{
		func(s *dgo.Session, r *dgo.Ready) {
			log.Println("DiscordInit StartSession: Session Started. Logged in as: `" + s.State.User.Username + "#" + s.State.User.Discriminator + "`")
		},
		func(s *dgo.Session, m *dgo.MessageCreate) {
			if m.Author.ID != s.State.User.ID {
				if _, err := s.ChannelMessageSend(m.ChannelID, m.Author.ID+" : "+m.Content); err != nil {
					log.Println("DiscordProcessors getBotCommandsAndHandlers: Error when sending message:", err)
				}
			}
		},
		func(s *dgo.Session, i *dgo.InteractionCreate) {
			var handler func(bot *Bot, s *dgo.Session, i *dgo.InteractionCreate)

			switch i.Type {
			case dgo.InteractionApplicationCommand:
				botCommand, ok := botCommands[i.ApplicationCommandData().Name]
				if ok {
					handler = botCommand.Handler
				} else {
					log.Println("DiscordProcessors getBotCommandsAndHandlers: No handler for interaction command", i.ApplicationCommandData().Name)
					return
				}
			case dgo.InteractionMessageComponent:
				fun, ok := componentHandlers[i.MessageComponentData().CustomID]
				if ok {
					handler = fun
				} else {
					log.Println("DiscordProcessors getBotCommandsAndHandlers: No handler for message command", i.MessageComponentData().CustomID)
					return
				}
			}

			handler(bot, s, i)
		},
	}
}
