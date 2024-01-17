package utils

import dgo "github.com/bwmarrin/discordgo"

func GetUserAndOptionMap(i *dgo.InteractionCreate) (*dgo.User, map[string]*dgo.ApplicationCommandInteractionDataOption) {
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

	return user, optionMap
}
