package Mongo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RootDocument struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Id      string             `bson:"id,omitempty"`
	OwnerId string             `bson:"ownerId,omitempty"`
}

type RegisteredUser struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	UId      string             `bson:"uId,omitempty"`
	NickName string             `bson:"nickName,omitempty"`
}

type RegisteredBlockchain struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	User             primitive.ObjectID `bson:"user,omitempty"`
	Name             string             `bson:"name,omitempty"`
	Symbol           string             `bson:"symbol,omitempty"`
	ChainId          string             `bson:"chainId,omitempty"`
	Enabled          bool               `bson:"enabled,omitempty"`
	WalletAddress    string             `bson:"walletAddress,omitempty"`
	WalletPrivateKey string             `bson:"walletPrivateKey,omitempty"`
}

type RegisteredRouter struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Blockchain primitive.ObjectID `bson:"blockchain,omitempty"`
	Name       string             `bson:"name,omitempty"`
	Address    string             `bson:"address,omitempty"`
}

type RegisteredPair struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Router     primitive.ObjectID `bson:"router,omitempty"`
	Blockchain primitive.ObjectID `bson:"blockchain,omitempty"`
	Name       string             `bson:"name,omitempty"`
	TokenA     string             `bson:"tokenA,omitempty"`
	TokenB     string             `bson:"tokenB,omitempty"`
}

type RegisteredTrade struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Pair         primitive.ObjectID `bson:"pair,omitempty"`
	Blockchain   primitive.ObjectID `bson:"blockchain,omitempty"`
	IsBuy        bool               `bson:"isBuy,omitempty"`
	TokenAddress string             `bson:"tokenAddress,omitempty"`
	Amount       string             `bson:"amount,omitempty"`
	TriggerPrice string             `bson:"triggerPrice,omitempty"`
}
