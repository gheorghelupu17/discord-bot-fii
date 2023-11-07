package botParts

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/session"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SpecialCommands(e *gateway.MessageCreateEvent, database *mongo.Database, specialRoles *mongo.Collection, s *session.Session, serversCollection *mongo.Collection, exportRolesCollection *mongo.Collection, discordBotName string) bool {


	log.Println(e.Author.ID)
	result := specialRoles.FindOne(nil, bson.M{"_id": fmt.Sprintf("%s", e.Author.ID)})
	
	if result.Err() == nil {
		if strings.HasPrefix(e.Content, "give") {
			members, err := s.Members(e.GuildID, 0)
			if err != nil {
				log.Println(err)
				return true
			}
			i := 0
			for _, member := range members {
				log.Println(member.User.Username)
				rolesToGive := make([]discord.RoleID, 0)
				if err := s.ModifyMember(e.GuildID, member.User.ID, api.ModifyMemberData{Roles: &rolesToGive}); err != nil {
					fmt.Println(err)
				} else {
					i += 1
					fmt.Println("Successfully cleaned", member.Nick, "member", i)
				}
				err := AddRolesToMember(s, database, e.GuildID, member.User.ID, false)
				if err != nil {
					fmt.Println("Error when giving role to", member.User.Username, err)
				}
			}
			fmt.Println("Done give roles")
			return false
		}
		if strings.HasPrefix(e.Content, "give_x_series") {
			members, err := s.Members(e.GuildID, 0)
			if err != nil {
				log.Println(err)
				return true
			}
			for _, member := range members {
				log.Println(member.User.Username)
				err := AddRolesToMember(s, database, e.GuildID, member.User.ID, true)
				if err != nil {
					fmt.Println("Error when giving role to", member.User.Username, err)
				}
			}
			fmt.Println("Done give roles to X series")
			return false
		}
		if strings.HasPrefix(e.Content, "verified") {
			members, err := s.Members(e.GuildID, 0)
			if err != nil {
				log.Println(err)
				return true
			}
			guild, err := s.Guild(e.GuildID)
			if err != nil {
				log.Println(err)
				return true
			}
			var verifiedRole discord.Role
			for _, role := range guild.Roles {
				if role.Name == "verified" {
					verifiedRole = role
				}
			}
			for _, member := range members {
				if member.User.ID == guild.OwnerID {
					continue
				}
				fmt.Printf("Member %s\n", member.User.Username)
				toNotGive := false
				for _, roleId := range member.RoleIDs {
					for _, guildRole := range guild.Roles {
						if roleId == guildRole.ID && (guildRole.Name == "not_assigned" || guildRole.Name == "verified") {
							toNotGive = true
						}
					}
				}
				if toNotGive == false {
					x := append(member.RoleIDs, verifiedRole.ID)
					if err := s.ModifyMember(guild.ID, member.User.ID, api.ModifyMemberData{Roles: &x}); err != nil {
						log.Println(err)
						return false
					}
				}
			}
			fmt.Println("finish")
			return false
		}
		if strings.HasPrefix(e.Content, "unite") {
			guild, err := s.Guild(e.GuildID)
			if err != nil {
				log.Println(err)
				return true
			}
			roles := make(map[string]string)
			for _, guildRole := range guild.Roles {
				roles[guildRole.Name] = fmt.Sprintf("%d", guildRole.ID)
			}
			upsert := true
			if _, err := serversCollection.UpdateOne(nil,
				bson.M{"_id": fmt.Sprintf("%d", guild.ID)},
				bson.D{{"$set", bson.D{{"name", guild.Name}, {"roles", roles}}}},
				&options.UpdateOptions{Upsert: &upsert},
			); err != nil {
				log.Println(err)
				return true
			}
			return false
		}
		if strings.HasPrefix(e.Content, "export") {
			members, err := s.Members(e.GuildID, 0)
			if err != nil {
				log.Println(err)
				return true
			}
			fmt.Println("Got all members")
			type User struct {
				DiscordId string            `json:"_id" bson:"_id"`
				Roles     map[string]string `json:"roles" bson:"roles"`
				Nickname  string            `json:"nickname" bson:"nickname"`
			}
			guild, err := s.Guild(e.GuildID)
			if err != nil {
				log.Println(err)
				return true
			}
			membersDetails := make([]User, 0)
			for _, member := range members {
				fmt.Printf("Member %s\n", member.User.Username)
				roles := make(map[string]string)
				for _, roleId := range member.RoleIDs {
					for _, guildRole := range guild.Roles {
						if roleId == guildRole.ID {
							roles[guildRole.Name] = fmt.Sprintf("%d", roleId)
						}
					}
				}
				nickname := member.Nick
				if len(nickname) == 0 {
					nickname = member.User.Username
				}
				membersDetails = append(membersDetails, User{
					DiscordId: fmt.Sprintf("%s", member.User.ID),
					Roles:     roles,
					Nickname:  nickname,
				})

			}
			upsert := true
			if _, err := exportRolesCollection.UpdateOne(nil,
				bson.M{"_id": fmt.Sprintf("%d", guild.ID)},
				bson.D{{"$set", bson.D{{"name", guild.Name}, {"members", membersDetails}}}},
				&options.UpdateOptions{Upsert: &upsert},
			); err != nil {
				log.Println(err)
				return true
			}
			log.Println("Exporting successful")
			return false
		}
		if strings.HasPrefix(e.Content, "licenta,") {
			args := strings.Split(strings.ReplaceAll(e.Content, "licenta,", ""), ",")
			if len(args) != 7 {
				//if _, err := s.SendMessage(e.ChannelID, fmt.Sprintf("Wrong no of args: must use `licenta,<server name>,<year>,<no a>,<no b>,<no e>,<no x>,<x inside>")); err != nil {
				//	log.Println(err)
				//	return true
				//}
				return true
			}
			year, err0 := strconv.Atoi(args[1])
			if err0 != nil {
				log.Println(err0)
				return true
			}
			noa, err1 := strconv.Atoi(args[2])
			if err1 != nil {
				log.Println(err1)
				return true
			}
			nob, err2 := strconv.Atoi(args[3])
			if err2 != nil {
				log.Println(err2)
				return true
			}
			noe, err3 := strconv.Atoi(args[4])
			if err3 != nil {
				log.Println(err3)
				return true
			}
			nox, err4 := strconv.Atoi(args[5])
			if err4 != nil {
				log.Println(err4)
				return true
			}
			config := BachelorConfig{ServerName: args[0],
				Year:          year,
				NoGroupsA:     noa,
				NoGroupsB:     nob,
				NoGroupsE:     noe,
				NoGroupsX:     nox,
				CustomXMiddle: args[6],
			}
			if err := CreateBachelorServer(serversCollection, s, e.GuildID, config, discordBotName); err != nil {
				log.Println(err)
				return true
			}
			return false
		}
		if strings.HasPrefix(e.Content, "masopt,") {
			args := strings.Split(strings.ReplaceAll(e.Content, "masopt,", ""), ",")
			if len(args) != 4 {
				//if _, err := s.SendMessage(e.ChannelID, fmt.Sprintf("Wrong no of args: must use `masopt,<server name>,<year>,<no groups>,<middle>")); err != nil {
				//	log.Println(err)
				//	return true
				//}
				return true
			}
			year, err0 := strconv.Atoi(args[1])
			if err0 != nil {
				log.Println(err0)
				return true
			}
			noGroups, err1 := strconv.Atoi(args[2])
			if err1 != nil {
				log.Println(err1)
				return true
			}
			config := MastersAndOptionalsConfig{ServerName: args[0],
				Year:     year,
				NoGroups: noGroups,
				Middle:   args[3],
			}
			if err := CreateMastersServer(serversCollection, s, e.GuildID, config, discordBotName); err != nil {
				log.Println(err)
				return true
			}
		}
		return false
	}
	return false
}

func PrivateDmHandler(e *gateway.MessageCreateEvent, privateChannels *mongo.Collection, discordIdsToStudIdsCollection *mongo.Collection, s *session.Session, studIdsToRolesCollection *mongo.Collection, serversCollection *mongo.Collection, botDb *mongo.Database) bool {
	channel, err := s.Channel(e.ChannelID)

	if err != nil {
		return false
	}
	if channel.Type != discord.DirectMessage {
		// return false
	}
	result1 := discordIdsToStudIdsCollection.FindOne(nil,bson.M{"_id": fmt.Sprintf("%d", e.Author.ID),"server_id":fmt.Sprintf("%d",e.GuildID)})
	log.Println("Result1", result1.Err(),fmt.Sprintf("%d", e.Author.ID))
	if result1.Err() == nil {
		if _, err := s.SendMessage(e.ChannelID, fmt.Sprintf("You already got your roles. If that's not the case contactAdmin")); err != nil {
			log.Println(err)
			return true
		}
		return false
	}

	// we are in a private message
	if strings.HasPrefix(e.Content, "#") && !strings.HasSuffix(e.Content, "#") {
		if _, err := s.SendMessage(e.ChannelID, fmt.Sprintf("Your input is incorrect")); err != nil {
			log.Println(err)
			return true
		}
		
		return true
	}

	if !strings.HasPrefix(e.Content, "#") && strings.HasSuffix(e.Content, "#") {
		if _, err := s.SendMessage(e.ChannelID, fmt.Sprintf("Your input is incorrect")); err != nil {
			log.Println(err)
			return true
		}
		
		return true
	}
	if strings.HasPrefix(e.Content, "#") && strings.HasSuffix(e.Content, "#") {

	possibleStudId := e.Content[1 : len(e.Content)-1]
	log.Println("Possible stud id", possibleStudId)
	result0 := studIdsToRolesCollection.FindOne(nil, bson.M{"_id": possibleStudId})
	if result0.Err() != nil {
		if _, err := s.SendMessage(e.ChannelID, fmt.Sprintf("Your student id was not found in the database. ContactAdmin")); err != nil {
			log.Println(err)
			return true
		}
	}
	var studEntry StudIDRolesNickDoc
	log.Println("Found stud id", result0.Decode(&studEntry))
	if err := result0.Decode(&studEntry); err != nil {
		log.Println(err)
		return true
	}
	if studEntry.DiscordId != "" && studEntry.DiscordId != fmt.Sprintf("%d", e.Author.ID) {
		//if _, err := s.SendMessage(e.ChannelID, fmt.Sprintf("Someone used your student id number. ContactAdmin")); err != nil {
		//	log.Println(err)
		//	return true
		//}
	}
	authorIdAsString := fmt.Sprintf("%d", e.Author.ID)
	upsert := true
	log.Println("Updating", authorIdAsString, "with", studEntry)
	if _, err := discordIdsToStudIdsCollection.UpdateOne(nil, bson.M{"_id": authorIdAsString},
		bson.D{
			{"$set", bson.D{{"stud_id", studEntry.StudId},{"serve_id",e.GuildID}}},
		}, &options.UpdateOptions{Upsert: &upsert},
	); err != nil {
		log.Println(err)
		return true
	}
	if _, err := studIdsToRolesCollection.UpdateOne(nil,
		bson.M{"_id": studEntry.StudId},
		bson.D{
			{"$set", bson.D{{"discord_id", authorIdAsString}}},
		},
	); err != nil {
		log.Println(err)
		return true
	}
	guildsCursor, _ := serversCollection.Find(nil, bson.M{})
	defer func(guildsCursor *mongo.Cursor, ctx context.Context) {
		err := guildsCursor.Close(ctx)
		if err != nil {
			log.Println(err)
			return
		}
	}(guildsCursor, context.Background())
	for guildsCursor.Next(context.Background()) {
		var serverDoc ServerDoc
		if err := guildsCursor.Decode(&serverDoc); err != nil {
			log.Println("print")

			log.Println(err)
			return true
		}
		var guildIdAsUint uint64
		_, err := fmt.Sscan(serverDoc.GuildID, &guildIdAsUint)
		if err != nil {
			log.Println(err)
			return true
		}
		if err := AddRolesToMember(s, botDb, discord.GuildID(guildIdAsUint), e.Author.ID, false); err != nil {
			log.Println(err)
			return true
		}
	}
	log.Println("Done")
}
	return false
}
