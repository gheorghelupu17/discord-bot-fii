package botParts

import (
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/session"
	"log"
)

const AuditReason string = "Command Bot"

// DeleteAllRoles is a helper for removing all roles except @everyone and the one used by the bot
func DeleteAllRoles(s *session.Session, guildID discord.GuildID, botRole string) error {
	roles, err := s.Roles(guildID)
	if err != nil {
		return err
	}
	for _, role := range roles {
		if role.Name == "@everyone" || role.Name == botRole {
			continue
		}
		if err := s.DeleteRole(guildID, role.ID, api.AuditLogReason(AuditReason)); err != nil {
			return err
		}
	}
	return nil
}

// DeleteAllCategoriesAndChannels is a helper for removing all categories and channels
func DeleteAllCategoriesAndChannels(s *session.Session, guildID discord.GuildID) error {
	channels, err := s.Channels(guildID)
	if err != nil {
		return err
	}
	for _, channel := range channels {
		if err := s.DeleteChannel(channel.ID, api.AuditLogReason(AuditReason)); err != nil {
			return err
		}
	}
	return nil
}

// DeleteAllFromServer is a helper for removing all roles, categories and channels
func DeleteAllFromServer(s *session.Session, guildID discord.GuildID, botRole string) error {
	log.Println("BEGIN: deleting all channels and roles")
	if err := DeleteAllRoles(s, guildID, botRole); err != nil {
		return err
	}
	if err := DeleteAllCategoriesAndChannels(s, guildID); err != nil {
		return err
	}
	log.Println("END: deleting all channels and roles")
	return nil
}
