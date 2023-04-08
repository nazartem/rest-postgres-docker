package product

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
	productsURL = "/products"
	productURL  = "/products/:uuid"
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
	router.HandlerFunc(http.MethodGet, productURL, apperror.Middleware(h.GetProduct))
	router.HandlerFunc(http.MethodGet, productsURL, apperror.Middleware(h.GetAllProducts))
	router.HandlerFunc(http.MethodPost, productsURL, apperror.Middleware(h.CreateProduct))
	router.HandlerFunc(http.MethodPatch, productURL, apperror.Middleware(h.UpdateProduct))
	router.HandlerFunc(http.MethodDelete, productURL, apperror.Middleware(h.DeleteProduct))
}

func (h *handler) GetProduct(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info.Println("GET PRODUCT")
	w.Header().Set("Content-Type", "application/json")

	h.logger.Info.Println("get uuid from context")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	productUUID := params.ByName("uuid")
	if productUUID == "" {
		return apperror.BadRequestError("uuid query parameter is required and must be a comma separated integers")
	}
	h.logger.Info.Printf("get param: %v", productUUID)

	product, err := h.repository.FindOne(r.Context(), productUUID)
	if err != nil {
		return err
	}
	productBytes, err := json.Marshal(product)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(productBytes)

	return nil
}

func (h *handler) GetAllProducts(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info.Println("GET ALL PRODUCTS")
	w.Header().Set("Content-Type", "application/json")

	h.logger.Info.Println("get category_uuid from URL")

	products, err := h.repository.FindAll(r.Context())
	if err != nil {
		return err
	}

	productsBytes, err := json.Marshal(products)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(productsBytes)

	return nil
}

func (h *handler) CreateProduct(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info.Println("CREATE PRODUCT")
	w.Header().Set("Content-Type", "application/json")

	var prd Product

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&prd); err != nil {
		return apperror.BadRequestError("invalid data")
	}

	err := h.repository.Create(r.Context(), &prd)
	if err != nil {
		return err
	}

	productUUID := prd.ID
	w.Header().Set("Location", fmt.Sprintf("%s/%v", productsURL, productUUID))
	w.WriteHeader(http.StatusCreated)

	return nil
}

func (h *handler) UpdateProduct(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info.Println("UPDATE PRODUCT")
	w.Header().Set("Content-Type", "application/json")

	h.logger.Info.Println("get uuid from context")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	productUUID := params.ByName("uuid")
	if productUUID == "" {
		return apperror.BadRequestError("uuid query parameter is required and must be a comma separated integers")
	}

	id, err := strconv.Atoi(productUUID)
	if err != nil {
		return err
	}

	var prd Product
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&prd); err != nil {
		return apperror.BadRequestError("invalid data")
	}
	prd.ID = id

	err = h.repository.Update(r.Context(), prd)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (h *handler) DeleteProduct(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info.Println("DELETE PRODUCT")
	w.Header().Set("Content-Type", "application/json")

	h.logger.Info.Println("get uuid from context")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	productUUID := params.ByName("uuid")
	if productUUID == "" {
		return apperror.BadRequestError("uuid query parameter is required and must be a comma separated integers")
	}

	err := h.repository.Delete(r.Context(), productUUID)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}
