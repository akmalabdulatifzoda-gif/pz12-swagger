package core

import (
	"errors"
	"time"
)

// Note — сущность «Заметка».
type Note struct {
	ID        int64      `json:"id" example:"1"`
	Title     string     `json:"title" example:"Идея для лекции"`
	Content   string     `json:"content" example:"Добавить пример структуры API"`
	CreatedAt time.Time  `json:"createdAt" example:"2026-11-05T12:00:00Z"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" example:"2026-11-05T12:30:00Z"`
}

// NoteCreate — DTO для создания заметки.
type NoteCreate struct {
	Title   string `json:"title" example:"Новая заметка" minLength:"1" maxLength:"200"`
	Content string `json:"content" example:"Текст заметки"`
}

// NoteUpdate — DTO для частичного обновления: nil означает «поле не меняется».
type NoteUpdate struct {
	Title   *string `json:"title,omitempty" example:"Обновлено" maxLength:"200"`
	Content *string `json:"content,omitempty" example:"Новый текст"`
}

// ErrorResponse — единый формат ошибки API.
type ErrorResponse struct {
	Error string `json:"error" example:"описание ошибки"`
}

// ListParams — параметры списка: пагинация и поиск по заголовку.
type ListParams struct {
	Page  int
	Limit int
	Query string
}

var (
	ErrNotFound   = errors.New("note not found")
	ErrValidation = errors.New("validation error")
)
