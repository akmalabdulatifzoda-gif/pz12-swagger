package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/akmalabdulatifzoda-gif/pz12-swagger/internal/core"
)

// NoteRepository — интерфейс слоя хранения. Сервис не знает,
// где лежат данные (память, PostgreSQL и т.д.).
type NoteRepository interface {
	Create(n core.Note) (core.Note, error)
	GetAll() ([]core.Note, error)
	GetByID(id int64) (core.Note, error)
	Update(n core.Note) (core.Note, error)
	Delete(id int64) error
}

// NoteService — бизнес-логика: валидация, временные метки.
type NoteService struct {
	repo NoteRepository
}

func NewNoteService(r NoteRepository) *NoteService {
	return &NoteService{repo: r}
}

func validateTitle(t string) error {
	if strings.TrimSpace(t) == "" {
		return fmt.Errorf("%w: title is required", core.ErrValidation)
	}
	if len([]rune(t)) > 200 {
		return fmt.Errorf("%w: title must be at most 200 characters", core.ErrValidation)
	}
	return nil
}

func (s *NoteService) Create(in core.NoteCreate) (core.Note, error) {
	if err := validateTitle(in.Title); err != nil {
		return core.Note{}, err
	}
	n := core.Note{
		Title:     strings.TrimSpace(in.Title),
		Content:   in.Content,
		CreatedAt: time.Now().UTC(),
	}
	return s.repo.Create(n)
}

// List возвращает страницу заметок и общее число найденных (до пагинации).
func (s *NoteService) List(p core.ListParams) ([]core.Note, int, error) {
	all, err := s.repo.GetAll()
	if err != nil {
		return nil, 0, err
	}
	if q := strings.ToLower(strings.TrimSpace(p.Query)); q != "" {
		filtered := all[:0:0]
		for _, n := range all {
			if strings.Contains(strings.ToLower(n.Title), q) {
				filtered = append(filtered, n)
			}
		}
		all = filtered
	}
	total := len(all)
	if p.Limit <= 0 || p.Limit > 100 {
		p.Limit = 20
	}
	if p.Page <= 0 {
		p.Page = 1
	}
	start := (p.Page - 1) * p.Limit
	if start >= total {
		return []core.Note{}, total, nil
	}
	end := start + p.Limit
	if end > total {
		end = total
	}
	return all[start:end], total, nil
}

func (s *NoteService) Get(id int64) (core.Note, error) { return s.repo.GetByID(id) }

func (s *NoteService) Update(id int64, p core.NoteUpdate) (core.Note, error) {
	n, err := s.repo.GetByID(id)
	if err != nil {
		return core.Note{}, err
	}
	if p.Title == nil && p.Content == nil {
		return core.Note{}, fmt.Errorf("%w: nothing to update", core.ErrValidation)
	}
	if p.Title != nil {
		if err := validateTitle(*p.Title); err != nil {
			return core.Note{}, err
		}
		n.Title = strings.TrimSpace(*p.Title)
	}
	if p.Content != nil {
		n.Content = *p.Content
	}
	now := time.Now().UTC()
	n.UpdatedAt = &now
	return s.repo.Update(n)
}

func (s *NoteService) Delete(id int64) error { return s.repo.Delete(id) }
