package botParts

type GroupData struct {
	HasVoice           bool
	RolePrefix         string
	Count              int
	CategoryTitle      string
	GeneraleRolePrefix string
	HasGeneralRole     bool
}

// Config for the creation of servers
type Config struct {
	ServerName     string
	Groups         []GroupData
	RolesWithAdmin []string
}
