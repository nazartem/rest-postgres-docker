package buyer

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
	buyersURL = "/buyers"
	buyerURL  = "/buyers/:uuid"
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
	router.HandlerFunc(http.MethodGet, buyerURL, apperror.Middleware(h.GetBuyer))
	router.HandlerFunc(http.MethodGet, buyersURL, apperror.Middleware(h.GetAllBuyers))
	router.HandlerFunc(http.MethodPost, buyersURL, apperror.Middleware(h.CreateBuyer))
	router.HandlerFunc(http.MethodPatch, buyerURL, apperror.Middleware(h.UpdateBuyer))
	router.HandlerFunc(http.MethodDelete, buyerURL, apperror.Middleware(h.DeleteBuyer))
}

func (h *handler) GetBuyer(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info.Println("GET BUYER")
	w.Header().Set("Content-Type", "application/json")

	h.logger.Info.Println("get uuid from context")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	buyerUUID := params.ByName("uuid")
	if buyerUUID == "" {
		return apperror.BadRequestError("uuid query parameter is required and must be a comma separated integers")
	}

	buyer, err := h.repository.FindOne(r.Context(), buyerUUID)
	if err != nil {
		return err
	}
	buyerBytes, err := json.Marshal(buyer)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(buyerBytes)

	return nil
}

func (h *handler) GetAllBuyers(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info.Println("GET ALL BUYERS")
	w.Header().Set("Content-Type", "application/json")

	h.logger.Info.Println("get category_uuid from URL")

	buyers, err := h.repository.FindAll(r.Context())
	if err != nil {
		return err
	}

	buyersBytes, err := json.Marshal(buyers)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(buyersBytes)

	return nil
}

func (h *handler) CreateBuyer(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info.Println("CREATE BUYER")
	w.Header().Set("Content-Type", "application/json")

	var br Buyer

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&br); err != nil {
		return apperror.BadRequestError("invalid data")
	}

	err := h.repository.Create(r.Context(), &br)
	if err != nil {
		return err
	}

	buyerUUID := br.ID
	w.Header().Set("Location", fmt.Sprintf("%s/%v", buyersURL, buyerUUID))
	w.WriteHeader(http.StatusCreated)

	return nil
}

func (h *handler) UpdateBuyer(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info.Println("UPDATE BUYER")
	w.Header().Set("Content-Type", "application/json")

	h.logger.Info.Println("get uuid from context")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	buyerUUID := params.ByName("uuid")
	if buyerUUID == "" {
		return apperror.BadRequestError("uuid query parameter is required and must be a comma separated integers")
	}
	id, err := strconv.Atoi(buyerUUID)
	if err != nil {
		return err
	}

	var br Buyer
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&br)
	if err != nil {
		return apperror.BadRequestError("invalid data")
	}

	br.ID = id

	err = h.repository.Update(r.Context(), br)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (h *handler) DeleteBuyer(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info.Println("DELETE BUYER")
	w.Header().Set("Content-Type", "application/json")

	h.logger.Info.Println("get uuid from context")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	buyerUUID := params.ByName("uuid")
	if buyerUUID == "" {
		return apperror.BadRequestError("uuid query parameter is required and must be a comma separated integers")
	}

	err := h.repository.Delete(r.Context(), buyerUUID)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}
