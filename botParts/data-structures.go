package botParts

type DiscordIDStudIDDoc struct {
	Id     string `json:"_id" bson:"_id"`
	StudId string `json:"stud_id" bson:"stud_id"`
	ServerId string `json:"server_id" bson:"server_id"`
}

type StudIDRolesNickDoc struct {
	StudId    string   `json:"_id" bson:"_id"`
	DiscordId string   `json:"discord_id" bson:"discord_id"`
	NickName  string   `json:"nickname" bson:"nickname"`
	Roles     []string `json:"roles" bson:"roles"`
}

type ServerDoc struct {
	GuildID string            `json:"_id" bson:"_id"`
	Name    string            `json:"name" bson:"name"`
	Roles   map[string]string `json:"roles" bson:"roles"`
}

type BachelorConfig struct {
	ServerName    string
	Year          int
	NoGroupsA     int
	NoGroupsB     int
	NoGroupsE     int
	NoGroupsX     int
	CustomXMiddle string
}

type MastersAndOptionalsConfig struct {
	ServerName string
	Year       int
	NoGroups   int
	Middle     string
}
