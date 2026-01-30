// Package utils содержит вспомогательные функции.
package utils

// Ptr возвращает указатель на значение.
func Ptr[T any](v T) *T {
	return &v
}