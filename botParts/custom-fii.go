package botParts

import (
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/session"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateBachelorServer(serversCollection *mongo.Collection, s *session.Session,
	guildID discord.GuildID, bachelorConfig BachelorConfig,
	botRoleName string) error {
	config := Config{
		ServerName: bachelorConfig.ServerName,
		Groups: append(make([]GroupData, 0),
			GroupData{RolePrefix: "teacher",
				HasVoice: true,
				Count:    1},
			GroupData{RolePrefix: fmt.Sprintf("%dA", bachelorConfig.Year),
				HasVoice:           true,
				Count:              bachelorConfig.NoGroupsA,
				CategoryTitle:      "Series A",
				GeneraleRolePrefix: "a_series",
				HasGeneralRole:     true,
			},
			GroupData{RolePrefix: fmt.Sprintf("%dB", bachelorConfig.Year),
				HasVoice:           true,
				Count:              bachelorConfig.NoGroupsB,
				CategoryTitle:      "Series B",
				GeneraleRolePrefix: "b_series",
				HasGeneralRole:     true,
			},
			GroupData{RolePrefix: fmt.Sprintf("%dE", bachelorConfig.Year),
				HasVoice:           true,
				Count:              bachelorConfig.NoGroupsE,
				CategoryTitle:      "Series E",
				GeneraleRolePrefix: "e_series",
				HasGeneralRole:     true,
			},
			GroupData{RolePrefix: fmt.Sprintf("%d%s", bachelorConfig.Year, bachelorConfig.CustomXMiddle),
				HasVoice:           true,
				Count:              bachelorConfig.NoGroupsX,
				CategoryTitle:      "Series X",
				GeneraleRolePrefix: "x_series",
				HasGeneralRole:     true,
			},
		),
	}
	err := CreateServer(serversCollection, s, guildID, config, botRoleName)
	if err != nil {
		return err
	}
	return nil
}

func CreateMastersServer(serversCollection *mongo.Collection, s *session.Session,
	guildID discord.GuildID, mastersAndOptionalConfig MastersAndOptionalsConfig,
	botRoleName string) error {
	config := Config{
		ServerName: mastersAndOptionalConfig.ServerName,
		Groups: append(make([]GroupData, 0),
			GroupData{RolePrefix: "teacher",
				HasVoice: true,
				Count:    1},
			GroupData{RolePrefix: fmt.Sprintf("%d%s", mastersAndOptionalConfig.Year, mastersAndOptionalConfig.Middle),
				HasVoice: true,
				Count:    mastersAndOptionalConfig.NoGroups,
			},
		),
	}
	err := CreateServer(serversCollection, s, guildID, config, botRoleName)
	if err != nil {
		return err
	}
	return nil
}
