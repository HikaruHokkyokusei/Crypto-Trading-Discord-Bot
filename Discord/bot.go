package Discord

import (
	"Crypto-Trading-Discord-Bot/Mongo"
	"context"
	dgo "github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

type Bot struct {
	ownerId            string
	session            *dgo.Session
	db                 *Mongo.Db
	botCommands        *map[string]BotCommand
	registeredCommands []*dgo.ApplicationCommand
	registeredUsers    map[string]*Mongo.RegisteredUser
}

func (bot Bot) OwnerId() string {
	return bot.ownerId
}

func (bot Bot) Db() *Mongo.Db {
	return bot.db
}

func (bot Bot) loadDatabase() {
	var rootDocument Mongo.RootDocument
	if err := bot.db.GetCollection("_ROOT").C.FindOne(context.TODO(), Mongo.RootDocument{Id: "0"}).Decode(&rootDocument); err != nil {
		log.Fatal("DiscordInit BuildSession: Unable to obtain root document", err)
	}
	bot.ownerId = rootDocument.OwnerId

	col := bot.Db().GetCollection("RegisteredUsers").C
	if cur, err := col.Find(context.TODO(), bson.D{}); err == nil {
		var results []Mongo.RegisteredUser
		if err := cur.All(context.TODO(), &results); err == nil {
			bot.registeredUsers = make(map[string]*Mongo.RegisteredUser)
			for _, user := range results {
				bot.registeredUsers[user.UId] = &user
			}
		} else {
			log.Fatal("DiscordInit loadDatabase: Error when decoding all RegisteredUsers ", err)
		}
	} else {
		log.Fatal("DiscordInit loadDatabase: Error when finding all RegisteredUser documents from DB ", err)
	}
}
