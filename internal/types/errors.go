// Package types содержит внутренние типы данных сервиса.
package types

import "errors"

// Бизнес-ошибки сервиса.
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidEmail      = errors.New("invalid email format")
	ErrInvalidPassword   = errors.New("password does not meet requirements")
	ErrUserBlocked       = errors.New("user is blocked")
)

// IsNotFound проверяет, является ли ошибка "не найдено".
func IsNotFound(err error) bool {
	return errors.Is(err, ErrUserNotFound)
}