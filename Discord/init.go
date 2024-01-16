package Discord

import (
	"Crypto-Trading-Discord-Bot/Mongo"
	dgo "github.com/bwmarrin/discordgo"
	"log"
)

type Bot struct {
	session *dgo.Session
	db      *Mongo.Db
}

func (bot Bot) Db() *Mongo.Db {
	return bot.db
}

func BuildBot(discordBotSecretToken string, db *Mongo.Db) *Bot {
	if session, err := dgo.New(discordBotSecretToken); err == nil {
		var bot = &Bot{
			session: session,
			db:      db,
		}

		for _, handler := range *GetBotHandlers(*bot) {
			session.AddHandler(handler)
		}

		// Use Bitwise OR for multiple intents...
		session.Identify.Intents = dgo.IntentDirectMessages

		return bot
	} else {
		log.Fatal("DiscordInit BuildSession: Unable to acquire a bot session")
		return nil
	}
}

func (bot Bot) StartSession() {
	log.Println("DiscordInit StartSession: Opening session")
	err := bot.session.Open()
	if err != nil {
		log.Fatal("DiscordInit StartSession: Unable to open session")
	}

	for _, botCommand := range GetBotCommands() {
		if _, err := bot.session.ApplicationCommandCreate(bot.session.State.User.ID, "", botCommand.Info); err != nil {
			log.Println("DiscordInit BuildBot: Error when creating command:", botCommand.Info.Name, "Error:", err)
		}
	}
}

func (bot Bot) EndSession() {
	err := bot.session.Close()
	if err != nil {
		log.Println("DiscordInit EndSession: Unable to close discord bot session", err)
	} else {
		log.Println("DiscordInit EndSession: Discord Bot Stopped")
	}
}
