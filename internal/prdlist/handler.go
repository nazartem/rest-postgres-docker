package prdlist

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
	prdListsURL = "/prdlists"
	prdListURL  = "/prdlists/:uuid"
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
	router.HandlerFunc(http.MethodGet, prdListURL, apperror.Middleware(h.GetProductList))
	router.HandlerFunc(http.MethodGet, prdListsURL, apperror.Middleware(h.GetAllProductLists))
	router.HandlerFunc(http.MethodPost, prdListsURL, apperror.Middleware(h.CreateProductList))
	router.HandlerFunc(http.MethodPatch, prdListURL, apperror.Middleware(h.UpdateProductList))
	router.HandlerFunc(http.MethodDelete, prdListURL, apperror.Middleware(h.DeleteProductList))
}

func (h *handler) GetProductList(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info.Println("GET PRODUCT LIST")
	w.Header().Set("Content-Type", "application/json")

	h.logger.Info.Println("get uuid from context")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	productListUUID := params.ByName("uuid")
	if productListUUID == "" {
		return apperror.BadRequestError("uuid query parameter is required and must be a comma separated integers")
	}

	productList, err := h.repository.FindOne(r.Context(), productListUUID)
	if err != nil {
		return err
	}
	productListBytes, err := json.Marshal(productList)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(productListBytes)

	return nil
}

func (h *handler) GetAllProductLists(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info.Println("GET ALL PRODUCT LISTS")
	w.Header().Set("Content-Type", "application/json")

	h.logger.Info.Println("get category_uuid from URL")

	productLists, err := h.repository.FindAll(r.Context())
	if err != nil {
		return err
	}

	productListsBytes, err := json.Marshal(productLists)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(productListsBytes)

	return nil
}

func (h *handler) CreateProductList(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info.Println("CREATE PRODUCT LIST")
	w.Header().Set("Content-Type", "application/json")

	var pl ProductList

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&pl); err != nil {
		return apperror.BadRequestError("invalid data")
	}

	err := h.repository.Create(r.Context(), &pl)
	if err != nil {
		return err
	}

	productListUUID := pl.ID
	w.Header().Set("Location", fmt.Sprintf("%s/%v", prdListsURL, productListUUID))
	w.WriteHeader(http.StatusCreated)

	return nil
}

func (h *handler) UpdateProductList(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info.Println("UPDATE PRODUCT LIST")
	w.Header().Set("Content-Type", "application/json")

	h.logger.Info.Println("get uuid from context")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	productListUUID := params.ByName("uuid")
	if productListUUID == "" {
		return apperror.BadRequestError("uuid query parameter is required and must be a comma separated integers")
	}
	id, err := strconv.Atoi(productListUUID)
	if err != nil {
		return err
	}

	var pl ProductList
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&pl); err != nil {
		return apperror.BadRequestError("invalid data")
	}

	pl.ID = id

	err = h.repository.Update(r.Context(), pl)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (h *handler) DeleteProductList(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info.Println("DELETE PRODUCT LIST")
	w.Header().Set("Content-Type", "application/json")

	h.logger.Info.Println("get uuid from context")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	productListUUID := params.ByName("uuid")
	if productListUUID == "" {
		return apperror.BadRequestError("uuid query parameter is required and must be a comma separated integers")
	}

	err := h.repository.Delete(r.Context(), productListUUID)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}
