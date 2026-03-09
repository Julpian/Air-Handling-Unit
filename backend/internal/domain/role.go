package domain

const (
	RoleAdmin      = "admin"
	RoleSupervisor = "Supervisor"
	RoleAsmen      = "AssistantManager"
	RoleInspector  = "inspector"
)

func IsAdminLike(role string) bool {
	return role == RoleAdmin || role == RoleSupervisor
}