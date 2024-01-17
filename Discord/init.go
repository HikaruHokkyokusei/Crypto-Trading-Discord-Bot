package Discord

import (
	"Crypto-Trading-Discord-Bot/Mongo"
	"context"
	dgo "github.com/bwmarrin/discordgo"
	"log"
)

type Bot struct {
	ownerId            string
	session            *dgo.Session
	db                 *Mongo.Db
	botCommands        *map[string]BotCommand
	registeredCommands []*dgo.ApplicationCommand
}

func (bot Bot) OwnerId() string {
	return bot.ownerId
}

func (bot Bot) Db() *Mongo.Db {
	return bot.db
}

func BuildBot(discordBotSecretToken string, db *Mongo.Db) *Bot {
	if session, err := dgo.New(discordBotSecretToken); err == nil {
		var rootDocument Mongo.RootDocument
		if err := db.GetCollection("_ROOT").C().FindOne(context.TODO(), Mongo.RootDocument{Id: "0"}).Decode(&rootDocument); err != nil {
			log.Fatal("DiscordInit BuildSession: Unable to obtain root document", err)
		}

		var bot = &Bot{
			ownerId: rootDocument.OwnerId,
			session: session,
			db:      db,
		}

		cmdMap, handlers := GetBotCommandsAndHandlers(bot)
		bot.botCommands = cmdMap

		for _, handler := range *handlers {
			session.AddHandler(handler)
		}

		// Use Bitwise OR for multiple intents...
		session.Identify.Intents = dgo.IntentDirectMessages

		return bot
	} else {
		log.Fatal("DiscordInit BuildSession: Unable to acquire a bot session", err)
		return nil
	}
}

func (bot Bot) StartSession() {
	log.Println("DiscordInit StartSession: Opening session")
	err := bot.session.Open()
	if err != nil {
		log.Fatal("DiscordInit StartSession: Unable to open session", err)
	}

	count := 0
	bot.registeredCommands = make([]*dgo.ApplicationCommand, len(*bot.botCommands))
	for _, botCommand := range *bot.botCommands {
		if registeredCommand, err := bot.session.ApplicationCommandCreate(bot.session.State.User.ID, botCommand.GuildId, botCommand.Info); err != nil {
			log.Fatal("DiscordInit BuildBot: Error when creating command:", botCommand.Info.Name, "Error:", err)
		} else {
			bot.registeredCommands[count] = registeredCommand
			count += 1
		}
	}
}

func (bot Bot) EndSession() {
	for _, registeredCommand := range bot.registeredCommands {
		if registeredCommand != nil {
			if err := bot.session.ApplicationCommandDelete(bot.session.State.User.ID, "", registeredCommand.ID); err != nil {
				log.Println("DiscordInit EndSession: Error when deleting registered command", err)
			}
		}
	}

	err := bot.session.Close()
	if err != nil {
		log.Println("DiscordInit EndSession: Unable to close discord bot session", err)
	} else {
		log.Println("DiscordInit EndSession: Discord Bot Stopped")
	}
}
