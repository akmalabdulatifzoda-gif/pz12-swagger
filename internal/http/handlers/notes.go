package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/akmalabdulatifzoda-gif/pz12-swagger/internal/core"
	"github.com/akmalabdulatifzoda-gif/pz12-swagger/internal/core/service"
)

type Handler struct {
	Svc *service.NoteService
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, core.ErrorResponse{Error: msg})
}

// handleErr переводит ошибки сервиса в HTTP-коды.
func handleErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, core.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, core.ErrValidation):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

func parseID(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id")
	}
	return id, nil
}

func queryInt(r *http.Request, key string) (int, error) {
	v := r.URL.Query().Get(key)
	if v == "" {
		return 0, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 0 {
		return 0, errors.New("invalid " + key)
	}
	return n, nil
}

// ListNotes godoc
//
//	@ID				listNotes
//
//	@Summary		Список заметок
//	@Description	Возвращает список заметок с пагинацией и фильтром по заголовку
//	@Tags			notes
//	@Produce		json
//	@Param			page	query		int		false	"Номер страницы (с 1)"		minimum(1)	default(1)
//	@Param			limit	query		int		false	"Размер страницы (1–100)"	minimum(1)	maximum(100)	default(20)
//	@Param			q		query		string	false	"Поиск по title (без учёта регистра)"
//	@Success		200		{array}		core.Note
//	@Header			200		{integer}	X-Total-Count	"Общее количество"
//	@Failure		400		{object}	core.ErrorResponse
//	@Failure		500		{object}	core.ErrorResponse
//	@Router			/notes [get]
func (h *Handler) ListNotes(w http.ResponseWriter, r *http.Request) {
	page, err := queryInt(r, "page")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	limit, err := queryInt(r, "limit")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	notes, total, err := h.Svc.List(core.ListParams{Page: page, Limit: limit, Query: r.URL.Query().Get("q")})
	if err != nil {
		handleErr(w, err)
		return
	}
	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	writeJSON(w, http.StatusOK, notes)
}

// CreateNote godoc
//
//	@ID				createNote
//
//	@Summary		Создать заметку
//	@Description	Создаёт новую заметку. Поле title обязательно.
//	@Tags			notes
//	@Accept			json
//	@Produce		json
//	@Param			input	body		core.NoteCreate	true	"Данные новой заметки"
//	@Success		201		{object}	core.Note
//	@Header			201		{string}	Location	"URL созданной заметки"
//	@Failure		400		{object}	core.ErrorResponse
//	@Failure		401		{object}	core.ErrorResponse
//	@Failure		500		{object}	core.ErrorResponse
//	@Security		BearerAuth
//	@Router			/notes [post]
func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var in core.NoteCreate
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	n, err := h.Svc.Create(in)
	if err != nil {
		handleErr(w, err)
		return
	}
	w.Header().Set("Location", "/api/v1/notes/"+strconv.FormatInt(n.ID, 10))
	writeJSON(w, http.StatusCreated, n)
}

// GetNote godoc
//
//	@ID				getNote
//
//	@Summary		Получить заметку
//	@Description	Возвращает данные заметки по идентификатору
//	@Tags			notes
//	@Produce		json
//	@Param			id	path		int	true	"ID заметки"	minimum(1)
//	@Success		200	{object}	core.Note
//	@Failure		400	{object}	core.ErrorResponse
//	@Failure		404	{object}	core.ErrorResponse
//	@Router			/notes/{id} [get]
func (h *Handler) GetNote(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	n, err := h.Svc.Get(id)
	if err != nil {
		handleErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, n)
}

// PatchNote godoc
//
//	@ID				patchNote
//
//	@Summary		Обновить заметку (частично)
//	@Description	Меняет только переданные поля; updatedAt выставляется автоматически
//	@Tags			notes
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int				true	"ID заметки"	minimum(1)
//	@Param			input	body		core.NoteUpdate	true	"Поля для обновления"
//	@Success		200		{object}	core.Note
//	@Failure		400		{object}	core.ErrorResponse
//	@Failure		401		{object}	core.ErrorResponse
//	@Failure		404		{object}	core.ErrorResponse
//	@Security		BearerAuth
//	@Router			/notes/{id} [patch]
func (h *Handler) PatchNote(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	var p core.NoteUpdate
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	n, err := h.Svc.Update(id, p)
	if err != nil {
		handleErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, n)
}

// DeleteNote godoc
//
//	@ID			deleteNote
//
//	@Summary	Удалить заметку
//	@Tags		notes
//	@Param		id	path	int	true	"ID заметки"	minimum(1)
//	@Success	204	"No Content"
//	@Failure	400	{object}	core.ErrorResponse
//	@Failure	401	{object}	core.ErrorResponse
//	@Failure	404	{object}	core.ErrorResponse
//	@Security	BearerAuth
//	@Router		/notes/{id} [delete]
func (h *Handler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.Svc.Delete(id); err != nil {
		handleErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
