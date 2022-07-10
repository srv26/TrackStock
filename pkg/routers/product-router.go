package routers

import (
	"github.com/gorilla/mux"
	"github.com/sourav/TrackStock/pkg/controllers"
	"github.com/sourav/TrackStock/pkg/models"
)

// RegisterTrackStockRoutes registers Handle functions for specific routes
var RegisterTrackStockRoutes = func(router *mux.Router, dbStore models.Repository) {
	handler := new(controllers.Handler)
	handler.DbStorage = dbStore
	router.HandleFunc("/product/", handler.AddItem).Methods("POST")
	router.HandleFunc("/product/{id}", handler.UpdateItem).Methods("PUT")
	router.HandleFunc("/product/{id}", handler.GetItemById).Methods("GET")
	router.HandleFunc("/product/{id}", handler.DeleteItemById).Methods("DELETE")
}
