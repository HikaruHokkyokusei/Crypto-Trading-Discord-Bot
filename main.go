package main

import (
	"Crypto-Trading-Discord-Bot/Discord"
	"Crypto-Trading-Discord-Bot/Mongo"
	"Crypto-Trading-Discord-Bot/utils"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

//goland:noinspection GoSnakeCaseUsage
var (
	discordBotSecretToken string
	mongoUsername         string
	mongoPassword         string
	mongoClusterHost      string
)

func init() {
	log.Println("こんにちは　世界...")

	discordBotSecretToken = utils.GetEnv("DISCORD_BOT_SECRET_TOKEN", "")
	if discordBotSecretToken == "" {
		log.Fatal("MAIN: No value provided for environment variable `DISCORD_BOT_SECRET_TOKEN`")
	}
	mongoUsername = utils.GetEnv("MONGO_USERNAME", "")
	if mongoUsername == "" {
		log.Fatal("MAIN: No value provided for environment variable `MONGO_USERNAME`")
	}
	mongoPassword = utils.GetEnv("MONGO_PASSWORD", "")
	if mongoPassword == "" {
		log.Fatal("MAIN: No value provided for environment variable `MONGO_PASSWORD`")
	}
	mongoClusterHost = utils.GetEnv("MONGO_CLUSTER_HOST", "")
	if mongoClusterHost == "" {
		log.Fatal("MAIN: No value provided for environment variable `MONGO_CLUSTER_HOST`")
	}
}

func main() {
	ctx := context.TODO()

	mongoClient := Mongo.Connect(mongoUsername, mongoPassword, mongoClusterHost, ctx)
	defer mongoClient.Disconnect(ctx)

	discordBot := Discord.BuildBot("Bot "+discordBotSecretToken, mongoClient.GetDB("Go-Discord-Crypto-Trading-Bot"))
	discordBot.StartSession()
	defer discordBot.EndSession()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	inputSignal := <-sc
	log.Print("Received signal: `", inputSignal, "`. Exiting Program...\n")
}
