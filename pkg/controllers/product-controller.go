package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sourav/TrackStock/pkg/models"
	"github.com/sourav/TrackStock/pkg/utils"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type Handler struct {
	DbStorage models.Repository
}

// route POST /product/ product AddItem
// Add a products details
//
// responses:
//	200: ProductResponse
//  404: errorResponse

// AddItem handles POST requests and Add items in the database
func (h *Handler) AddItem(response http.ResponseWriter, r *http.Request) {

	addProduct := new(models.Product)
	if err := utils.ParseBody(r, addProduct); err != nil {
		res := Response{Success: false, Data: err.Error()}

		utils.ResponseWithError(response, http.StatusInternalServerError, res)
		return
	}
	err := h.DbStorage.AddItem(addProduct)

	if err != nil {
		res := Response{Success: false, Data: err.Error()}
		utils.ResponseWithError(response, http.StatusInternalServerError, res)
		return
	}
	response.WriteHeader(http.StatusNoContent)
}

// route GET /products/{id} products GetItemById
// Gets an input item
//
// responses:
//	200: ProductResponse
//  404 errorResponse

// GetItemById handles GET requests and Get item by id from the database
func (h *Handler) GetItemById(response http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	param := mux.Vars(r)
	id := param["id"]
	intId, _ := strconv.Atoi(id)
	pdt, err := h.DbStorage.GetItemById(intId)
	switch err {
	case nil:
	case models.ErrProductNotFound:
		res := Response{Success: false, Data: err.Error()}
		utils.ResponseWithError(response, http.StatusNotFound, res)
		return
	default:
		res := Response{Success: false, Data: err.Error()}
		utils.ResponseWithError(response, http.StatusInternalServerError, res)
		return
	}
	res := Response{Success: true, Data: pdt}
	utils.ResponseWithJson(response, http.StatusOK, res)
}

// route PUT /product/ product UpdateItem
// Gets an input item
//
// responses:
//	201: noContentResponse
//  404: errorResponse
// 422: error validation

// UpdateItem handles PUT requests to update an item
func (h *Handler) UpdateItem(response http.ResponseWriter, r *http.Request) {

	updateProduct := new(models.Product)
	if err := utils.ParseBody(r, updateProduct); err != nil {
		res := Response{Success: false, Data: err.Error()}

		utils.ResponseWithError(response, http.StatusInternalServerError, res)
		return
	}
	err := h.DbStorage.UpdateItem(updateProduct)
	switch err {
	case nil:
	case models.ErrProductNotFound:
		res := Response{Success: false, Data: err.Error()}
		utils.ResponseWithError(response, http.StatusNotFound, res)
		return
	default:
		res := Response{Success: false, Data: err.Error()}
		utils.ResponseWithError(response, http.StatusInternalServerError, res)
		return
	}
	response.WriteHeader(http.StatusNoContent)
}

// route DELETE /product/{id} product DeleteItem
// Update a item list
//
// responses:
//	201: noContentResponse
//  404: errorResponse
// 422: error validation

// DeleteItemById handles DELETE requests to remove an item from database
func (h *Handler) DeleteItemById(response http.ResponseWriter, r *http.Request) {

	param := mux.Vars(r)
	id := param["id"]
	intId, _ := strconv.Atoi(id)
	err := h.DbStorage.DeleteItemById(intId)
	switch err {
	case nil:
	case models.ErrProductNotFound:
		res := Response{Success: false, Data: err.Error()}
		utils.ResponseWithError(response, http.StatusNotFound, res)
		return
	default:
		res := Response{Success: false, Data: err.Error()}
		utils.ResponseWithError(response, http.StatusInternalServerError, res)
		return
	}
	response.WriteHeader(http.StatusNoContent)
}
