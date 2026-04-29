package constants

type Role string

const (
	ADMIN   Role = "ADMIN"
	MANAGER Role = "MANAGER"
)

var AllRoles = []Role{ADMIN, MANAGER}
