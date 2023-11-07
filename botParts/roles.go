package botParts

import (
	"fmt"
	"strings"
	"log"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/session"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddRolesToMember(s *session.Session, d *mongo.Database, guildID discord.GuildID, userId discord.UserID, onlyXSeries bool) error {

	member, err := s.Member(guildID, userId)
	if err != nil || member == nil {
		log.Println(err)
		return nil
	}

	membersCollection := d.Collection("members")
	
	result12 := membersCollection.FindOne(nil, bson.D{{"_id", fmt.Sprintf("%s%s", guildID, userId)}})
	if result12.Err() != nil {
		return nil
	}

	serversCollection := d.Collection("servers")
	studIdsToRolesCollection := d.Collection("roles")
	discordIdsToStudIdsCollection := d.Collection("discord_ids_to_stud_ids")
	privateChannels := d.Collection("dms")
	channel, err := s.CreatePrivateChannel(userId)

	if err != nil {
		return err
	}
	result1 := studIdsToRolesCollection.FindOne(nil, bson.M{"discord_id": fmt.Sprintf("%s", userId)})
	if result1.Err() == mongo.ErrNoDocuments {
		err2 := sendFirstMessage(s, privateChannels, userId, channel.ID)
		if err2 != nil {
			return err2
		}
		log.Println("alexa4.1")
		return nil
	}
	var studIdDoc StudIDRolesNickDoc
	if err := result1.Decode(&studIdDoc); err != nil {
		return err
	}
	log.Println(guildID);
	guildResult := serversCollection.FindOne(nil, bson.M{"_id": fmt.Sprintf("%s", guildID)})
	if guildResult.Err() == mongo.ErrNoDocuments {
		return nil
		if _, err := s.SendMessage(channel.ID, "Guild not configured contact Catalin Marincia"); err != nil {
			return err
		}
		return nil
	}
	var guild ServerDoc
	if err := guildResult.Decode(&guild); err != nil {
		return err
	}
	rolesToGive := make([]discord.RoleID, 0)
	for _, roleName := range studIdDoc.Roles {
		if _, found := guild.Roles[roleName]; found {
			var roleIdUint64 int64
			_, err := fmt.Sscan(guild.Roles[roleName], &roleIdUint64)
			if err != nil {
				return err
				
			}
			if strings.Contains(roleName, "X") {
				if _, found1 := guild.Roles["x_series"]; found1 {
					var xSeriesRole int64
					_, err := fmt.Sscan(guild.Roles["x_series"], &xSeriesRole)
					if err != nil {
						log.Println("alexa8.2")
						return err
					}
					rolesToGive = append(rolesToGive, discord.RoleID(xSeriesRole))
				}
			}
			if !onlyXSeries {
				rolesToGive = append(rolesToGive, discord.RoleID(roleIdUint64))
			}
		}
	}
	if len(rolesToGive) == 0 {
		return nil
	}
	if err := s.ModifyMember(guildID, userId, api.ModifyMemberData{Nick: &studIdDoc.NickName, Roles: &rolesToGive}); err != nil {
		return err
	}
	if _, err := s.SendMessage(channel.ID, "You received your roles and nickname"); err != nil {
		return err
	}
	upsert := true
	if _, err := discordIdsToStudIdsCollection.UpdateOne(nil, bson.M{"_id": fmt.Sprintf("%s", userId)},
		bson.D{
			{"$set", bson.D{{"stud_id", studIdDoc.StudId}}},
		}, &options.UpdateOptions{Upsert: &upsert},
	); err != nil {
		return err
	}
	return nil
}

func sendFirstMessage(s *session.Session, privateChannels *mongo.Collection, userId discord.UserID, channelID discord.ChannelID) error {
	upsert := true
	if _, err := privateChannels.UpdateOne(nil, bson.M{"_id": fmt.Sprintf("%d", userId)}, bson.D{
		{"$set", bson.D{{"channel_id", fmt.Sprintf("%d", channelID)}}},
	}, &options.UpdateOptions{Upsert: &upsert}); err != nil {
		return err
	}
	return nil
	//	const botDm string = `Hello!
	//You didn't get your roles automatically assigned!
	//This is maybe because you never joined another server using this procedure.
	//To get your roles, type BELOW -in this conversation and not on #unverified- between 2 '#' symbols your:
	//- STUDENT ID(numar matricol)(if you are in years 2, 3 licenta and 2 master)
	//- the code you received(licenta 1 and master 1)(you do not have your STUDENT ID and we are not authorized to share it with you, only the faculty can).
	//If you didn't receive a code contact Alin Ursu(Alin#0763), Catalin Admin or Andrei Scutelnicu#7790
	//Like so: #1234567890123#. Why? Because this will help you avoid adding extra spaces.
	//Questions:
	//How do I use your student ID?
	//- I use it just for searching your roles in a database I made and assign them to you.
	//I did as you told me to do and I still didn't receive my roles. What should I do now?
	//- Contact me: MessengerAdmin, Discord Catalin Admin
	//Other notes:
	//- try to use your student id and not other's or you will be banned.
	//- this is a dm, you and the bot are the only ones that see your messages, you can delete them when you want.
	//Have a great day!
	//`
	//	if _, err := s.SendMessage(channelID, botDm); err != nil {
	//		return err
	//	}
	//	return nil
}
