package note

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"restapi-lesson/internal/apperror"
	"restapi-lesson/internal/handlers"
	"restapi-lesson/internal/logging"
	"strconv"
)

const (
	notesURL = "/notes"
	noteURL  = "/notes/:uuid"
)

type handler struct {
	logger     *logging.Logger
	repository Repository
}

func NewHandler(repository Repository, logger *logging.Logger) handlers.Handler {
	return &handler{
		repository: repository,
		logger:     logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, noteURL, apperror.Middleware(h.GetNote))
	router.HandlerFunc(http.MethodGet, notesURL, apperror.Middleware(h.GetAllNotes))
	router.HandlerFunc(http.MethodPost, notesURL, apperror.Middleware(h.CreateNote))
	router.HandlerFunc(http.MethodPatch, noteURL, apperror.Middleware(h.UpdateNote))
	router.HandlerFunc(http.MethodDelete, noteURL, apperror.Middleware(h.DeleteNote))
}

func (h *handler) GetNote(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info.Println("GET NOTE")
	w.Header().Set("Content-Type", "application/json")

	h.logger.Info.Println("get uuid from context")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	noteUUID := params.ByName("uuid")
	if noteUUID == "" {
		return apperror.BadRequestError("uuid query parameter is required and must be a comma separated integers")
	}

	note, err := h.repository.FindOne(r.Context(), noteUUID)
	if err != nil {
		return err
	}

	noteBytes, err := json.Marshal(note)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(noteBytes)

	return nil
}

func (h *handler) GetAllNotes(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info.Println("GET ALL NOTES")
	w.Header().Set("Content-Type", "application/json")

	h.logger.Info.Println("get category_uuid from URL")

	notes, err := h.repository.FindAll(r.Context())
	if err != nil {
		return err
	}

	notesBytes, err := json.Marshal(notes)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(notesBytes)

	return nil
}

func (h *handler) CreateNote(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info.Println("CREATE NOTE")
	w.Header().Set("Content-Type", "application/json")

	var nt Note

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&nt); err != nil {
		return apperror.BadRequestError("invalid data")
	}

	err := h.repository.Create(r.Context(), &nt)
	if err != nil {
		return err
	}

	noteNumber := nt.Number
	w.Header().Set("Location", fmt.Sprintf("%s/%v", notesURL, noteNumber))
	w.WriteHeader(http.StatusCreated)

	return nil
}

func (h *handler) UpdateNote(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info.Println("UPDATE NOTE")
	w.Header().Set("Content-Type", "application/json")

	h.logger.Info.Println("get uuid from context")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	noteNumber := params.ByName("uuid")
	if noteNumber == "" {
		return apperror.BadRequestError("uuid query parameter is required and must be a comma separated integers")
	}
	number, err := strconv.Atoi(noteNumber)
	if err != nil {
		return err
	}

	var nt Note
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&nt); err != nil {
		return apperror.BadRequestError("invalid data")
	}

	nt.Number = number

	err = h.repository.Update(r.Context(), nt)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (h *handler) DeleteNote(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info.Println("DELETE NOTE")
	w.Header().Set("Content-Type", "application/json")

	h.logger.Info.Println("get uuid from context")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	noteNumber := params.ByName("uuid")
	if noteNumber == "" {
		return apperror.BadRequestError("uuid query parameter is required and must be a comma separated integers")
	}

	err := h.repository.Delete(r.Context(), noteNumber)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}
