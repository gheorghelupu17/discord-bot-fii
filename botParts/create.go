package botParts

import (
	"fmt"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/session"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type CategoryCumulativeData struct {
	Name                     string
	HasVoice                 bool
	ReadWritePermissionRoles []discord.Role
	ReadPermissionRoles      []discord.Role
}

func CreateCategoryWithChannels(everyoneRoleID discord.RoleID, s *session.Session, guildID discord.GuildID, categoryCumulativeData CategoryCumulativeData) (*discord.Channel, error) {
	readOverwrites := make([]discord.Overwrite, 0)
	if categoryCumulativeData.ReadPermissionRoles != nil && len(categoryCumulativeData.ReadPermissionRoles) != 0 {
		for _, readPermissionRole := range categoryCumulativeData.ReadPermissionRoles {
			readOverwrites = append(readOverwrites, discord.Overwrite{ID: discord.Snowflake(readPermissionRole.ID),
				Type: discord.OverwriteRole,
				Allow: discord.PermissionAddReactions |
					discord.PermissionViewChannel |
					discord.PermissionReadMessageHistory |
					discord.PermissionConnect,
				Deny: discord.PermissionAll})
		}
	}
	writeOverwrites := make([]discord.Overwrite, 0)
	onlyAdminOverwrites := make([]discord.Overwrite, 0)
	if categoryCumulativeData.ReadWritePermissionRoles != nil && len(categoryCumulativeData.ReadWritePermissionRoles) != 0 {
		for _, writePermissionRole := range categoryCumulativeData.ReadWritePermissionRoles {
			onlyAdminOverwrites = append(onlyAdminOverwrites, discord.Overwrite{ID: discord.Snowflake(writePermissionRole.ID),
				Type: discord.OverwriteRole,
				Allow: discord.PermissionAddReactions |
					discord.PermissionViewChannel |
					discord.PermissionReadMessageHistory |
					discord.PermissionConnect,
				Deny: discord.PermissionAll})
			writeOverwrites = append(writeOverwrites, discord.Overwrite{ID: discord.Snowflake(writePermissionRole.ID), Type: discord.OverwriteRole,
				Allow: discord.PermissionAddReactions |
					discord.PermissionStream |
					discord.PermissionViewChannel |
					discord.PermissionSendMessages |
					discord.PermissionSendTTSMessages |
					discord.PermissionEmbedLinks |
					discord.PermissionAttachFiles |
					discord.PermissionReadMessageHistory |
					discord.PermissionMentionEveryone |
					discord.PermissionConnect |
					discord.PermissionSpeak |
					discord.PermissionUsePublicThreads |
					discord.PermissionRequestToSpeak |
					discord.PermissionUseVAD,
				Deny: discord.PermissionAll})
		}
	}
	everyoneRestriction := discord.Overwrite{ID: discord.Snowflake(everyoneRoleID), Type: discord.OverwriteRole, Deny: discord.PermissionAll}
	onlyAdminOverwrites = append(onlyAdminOverwrites, everyoneRestriction)
	readOverwrites = append(readOverwrites, everyoneRestriction)

	category, err := s.CreateChannel(guildID, api.CreateChannelData{Name: categoryCumulativeData.Name, Type: discord.GuildCategory, Overwrites: append(readOverwrites, writeOverwrites...)})
	if err != nil {
		return nil, err
	}
	_, err2 := s.CreateChannel(guildID, api.CreateChannelData{Name: "announcements", Type: discord.GuildText, CategoryID: category.ID, Overwrites: onlyAdminOverwrites})
	if err2 != nil {
		return nil, err2
	}
	_, err1 := s.CreateChannel(guildID, api.CreateChannelData{Name: "general", Type: discord.GuildText, CategoryID: category.ID})
	if err1 != nil {
		return nil, err1
	}
	if categoryCumulativeData.HasVoice == true {
		_, err2 := s.CreateChannel(guildID, api.CreateChannelData{Name: "voice", Type: discord.GuildVoice, CategoryID: category.ID})
		if err2 != nil {
			return nil, err2
		}
	}
	return category, nil
}

func ChangeNameOfGuild(s *session.Session, guildID discord.GuildID, name string) error {
	log.Println("BEGIN: changing the name of the server")
	_, err := s.ModifyGuild(guildID, api.ModifyGuildData{
		Name: name,
	})
	if err != nil {
		return err
	}
	log.Println("END: changing the name of the server")
	return nil
}
func RestrictChangeNickname(s *session.Session, guildID discord.GuildID) (*discord.Role, error) {
	var everyoneRole discord.Role
	roles, err0 := s.Roles(guildID)
	if err0 != nil {
		return nil, err0
	}
	for _, role := range roles {
		if role.ID == discord.RoleID(guildID) {
			everyoneRole = role
			break
		}
	}
	updated := everyoneRole
	if everyoneRole.Permissions.Has(discord.PermissionChangeNickname) {
		fmt.Println("Restricting everyone role")
		newPerms0 := everyoneRole.Permissions ^ discord.PermissionChangeNickname
		s, err := s.ModifyRole(guildID, discord.RoleID(guildID), api.ModifyRoleData{Permissions: &newPerms0})
		if err != nil {
			return nil, err
		}
		updated = *s
	}
	return &updated, nil
}
func CreateServerStructure(serverDoc *ServerDoc, s *session.Session, guildID discord.GuildID, config Config) error {
	log.Println("BEGIN: creating the server structure")
	if _, err := RestrictChangeNickname(s, guildID); err != nil {
		return err
	}
	mongoRolesMap := make(map[string]string)
	adminRole, err := s.CreateRole(guildID, api.CreateRoleData{Name: "admin", Mentionable: true, Permissions: discord.PermissionAdministrator})
	if err != nil {
		return err
	}
	mongoRolesMap["admin"] = fmt.Sprintf("%d", adminRole.ID)
	verifiedRole, err := s.CreateRole(guildID, api.CreateRoleData{Name: "verified", Mentionable: true})
	if err != nil {
		return err
	}
	mongoRolesMap["verified"] = fmt.Sprintf("%d", verifiedRole.ID)
	categoryWelcome, err := s.CreateChannel(guildID, api.CreateChannelData{Name: "Welcome", Type: discord.GuildCategory, Overwrites: append(make([]discord.Overwrite, 0), discord.Overwrite{ID: discord.Snowflake(verifiedRole.ID), Type: discord.OverwriteRole,
		Deny: discord.PermissionViewChannel})})
	if err != nil {
		return err
	}
	_, err1 := s.CreateChannel(guildID, api.CreateChannelData{Name: "unverified", Type: discord.GuildText, CategoryID: categoryWelcome.ID})
	if err1 != nil {
		return err1
	}
	if _, err := CreateCategoryWithChannels(discord.RoleID(guildID), s, guildID, CategoryCumulativeData{
		Name:                     "All people",
		ReadWritePermissionRoles: append(make([]discord.Role, 0), *verifiedRole),
		HasVoice:                 true,
	}); err != nil {
		return err
	}
	if _, err := CreateCategoryWithChannels(discord.RoleID(guildID), s, guildID, CategoryCumulativeData{
		Name:                     "Admin",
		ReadWritePermissionRoles: append(make([]discord.Role, 0), *adminRole),
		HasVoice:                 true,
	}); err != nil {
		return err
	}
	if config.Groups != nil {
		for _, group := range config.Groups {
			if group.HasGeneralRole == true {
				generalRole, err := s.CreateRole(guildID, api.CreateRoleData{Name: group.GeneraleRolePrefix})
				if err != nil {
					return err
				}
				mongoRolesMap[group.GeneraleRolePrefix] = fmt.Sprintf("%d", generalRole.ID)
				rolesList := append(make([]discord.Role, 0), *generalRole)
				if _, err := CreateCategoryWithChannels(discord.RoleID(guildID), s, guildID, CategoryCumulativeData{
					Name:                     group.CategoryTitle,
					ReadWritePermissionRoles: rolesList,
					HasVoice:                 group.HasVoice,
				}); err != nil {
					return err
				}
			}
			for i := 0; i < group.Count; i++ {
				roleName := fmt.Sprintf("%s%d", group.RolePrefix, i+1)
				categoryName := fmt.Sprintf("Group %s%d", group.RolePrefix, i+1)
				if group.Count == 1 {
					roleName = roleName[:len(roleName)-1]
					categoryName = categoryName[:len(categoryName)-1]
				}
				groupRole, err := s.CreateRole(guildID, api.CreateRoleData{Name: roleName, Mentionable: true})
				mongoRolesMap[roleName] = fmt.Sprintf("%d", groupRole.ID)
				if err != nil {
					return err
				}
				rolesList := append(make([]discord.Role, 0), *groupRole)
				if _, err := CreateCategoryWithChannels(discord.RoleID(guildID), s, guildID, CategoryCumulativeData{
					Name:                     categoryName,
					ReadWritePermissionRoles: rolesList,
					HasVoice:                 group.HasVoice,
				}); err != nil {
					return err
				}
			}
		}
	}
	log.Println(serverDoc.Name)
	log.Println("END: creating the server structure")

	if _, err := RestrictChangeNickname(s, guildID); err != nil {
		return err
	}
	serverDoc.Roles = mongoRolesMap
	return nil
}

func CreateServer(serversCollection *mongo.Collection, s *session.Session, guildID discord.GuildID, config Config, botRole string) error {
	serverDoc := ServerDoc{GuildID: fmt.Sprintf("%d", guildID), Name: config.ServerName}
	if err := DeleteAllFromServer(s, guildID, botRole); err != nil {
		return err
	}
	if err := ChangeNameOfGuild(s, guildID, config.ServerName); err != nil {
		return err
	}
	if err := CreateServerStructure(&serverDoc, s, guildID, config); err != nil {
		return err
	}
	upsert := true
	if _, err := serversCollection.UpdateOne(nil,
		bson.M{"_id": serverDoc.GuildID},
		bson.D{{"$set", bson.D{{"name", serverDoc.Name}, {"roles", serverDoc.Roles}}}},
		&options.UpdateOptions{Upsert: &upsert},
	); err != nil {
		return err
	}
	return nil
}
