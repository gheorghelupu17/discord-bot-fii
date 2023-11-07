package main

import (
	"context"
	"fiibot/botParts"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/session"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	options "go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	discordBotName := os.Getenv("BOT_NAME")
	if discordBotName == "" {
		log.Fatalln("No $BOT_NAME provided.")
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatalln(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
		}
	}(client, ctx)
	botDb := client.Database("bot")
	serversCollection := botDb.Collection("servers")
	privateChannels := botDb.Collection("dms")
	discordIdsToStudIdsCollection := botDb.Collection("discord_ids_to_stud_ids")
	// studIdsToRolesCollection := botDb.Collection("stud_ids_to_roles")
	studIdsToRolesCollection := botDb.Collection("roles")

	members := botDb.Collection("members")
	specialRolesCollection := botDb.Collection("special_roles")
	exportRolesCollection := botDb.Collection("roles_export")
	if err != nil {
		log.Fatalln(err)
	}
	s, err := botParts.GetSession()
	if err != nil {
		log.Fatalln(err)
	}
	// guilds, err := s.Guilds(200)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// for _, g := range guilds {
	// 	log.Println(g.ID.String())
	// 	if g.ID.String() != "" {
	// 		fmt.Println("Leaving", g.Name)
	// 		err := s.LeaveGuild(g.ID)
	// 		if err != nil {
	// 			log.Fatalln(err)
	// 		}
	// 	} else {
	// 		fmt.Println("Skipping since it is the main server")
	// 	}
	// }
	s.AddHandler(func(e *gateway.ReadyEvent) {
		log.Println("The bot is ready")
	})
	s.AddHandler(func(e *gateway.GuildMemberAddEvent) {
		upsert := true
		if _, err := members.UpdateOne(nil, bson.M{"_id": fmt.Sprintf("%s%s", e.GuildID, e.User.ID)}, bson.D{{"$set", bson.M{
			"present": true,
		}}}, &options.UpdateOptions{Upsert: &upsert}); err != nil {
			log.Println(err)
			return
		}
		if err := botParts.AddRolesToMember(s, botDb, e.GuildID, e.User.ID, false); err != nil {
			log.Println(err)
			return
		}
		
	})
	s.AddHandler(func(e *gateway.MessageCreateEvent) {
		if e.Author.Bot {
			return
		}
		if botParts.PrivateDmHandler(e, privateChannels, discordIdsToStudIdsCollection, s, studIdsToRolesCollection, serversCollection, botDb) {
		  log.Println("Private DM")
		}
		log.Println(e.Message.Content)
		if botParts.SpecialCommands(e, botDb, specialRolesCollection, s, serversCollection, exportRolesCollection, discordBotName) {
			return
		}
	})
	s.AddIntents(gateway.IntentGuilds)
	s.AddIntents(gateway.IntentGuildMessages)
	s.AddIntents(gateway.IntentGuildMembers)
	s.AddIntents(gateway.IntentDirectMessages)
	if err := s.Open(context.Background()); err != nil {
		log.Fatalln("Failed to open:", err)
	}
	defer func(s *session.Session) {
		err := s.Close()
		if err != nil {
			log.Fatalln("Error closing session: ", err)
		}
	}(s)
	// Block forever
	select {}
}
