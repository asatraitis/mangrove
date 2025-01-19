package dto

type UserRole string

const (
	USER_ROLE_USER      UserRole = "user"
	USER_ROLE_ADMIN     UserRole = "admin"
	USER_ROLE_SUPERUSER UserRole = "superadmin"
)

type UserStatus string

const (
	USER_STATUS_ACTIVE    UserStatus = "active"
	USER_STATUS_INACTIVE  UserStatus = "inactive"
	USER_STATUS_PENDING   UserStatus = "pending"
	USER_STATUS_SUSPENDED UserStatus = "suspended"
)

type MeResponse struct {
	ID          string     `json:"id"`
	DisplayName string     `json:"displayName"`
	Role        UserRole   `json:"role"`
	Status      UserStatus `json:"status"`
}
