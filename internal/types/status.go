package types

// UserStatus - статус пользователя.
type UserStatus int

const (
	UserStatusUnspecified UserStatus = iota
	UserStatusActive
	UserStatusInactive
	UserStatusBlocked
)

// String возвращает строковое представление статуса.
func (s UserStatus) String() string {
	switch s {
	case UserStatusActive:
		return "active"
	case UserStatusInactive:
		return "inactive"
	case UserStatusBlocked:
		return "blocked"
	default:
		return "unspecified"
	}
}

// IsValid проверяет валидность статуса.
func (s UserStatus) IsValid() bool {
	return s >= UserStatusActive && s <= UserStatusBlocked
}