package domain

const (
	RoleAdmin     = "admin"
	RoleSVP       = "svp"
	RoleAsmen     = "asmen"
	RoleInspector = "inspector"
)

func IsAdminLike(role string) bool {
	return role == RoleAdmin || role == RoleSVP
}
