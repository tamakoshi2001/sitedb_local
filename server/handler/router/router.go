package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tamakoshi2001/gextension/handler"
	"github.com/tamakoshi2001/gextension/handler/middleware"
	"github.com/tamakoshi2001/gextension/model"
	"github.com/tamakoshi2001/gextension/service"
)

func NewRouter(sites *[]model.Site, vectors *[]model.Vector) *mux.Router {
	// register routes
	r := mux.NewRouter()

	siteService := service.NewSiteService(sites, vectors)
	siteHandler := handler.NewSiteHandler(siteService)
	r.Handle("/site", middleware.Recovery(http.HandlerFunc(siteHandler.HandlePost))).Methods("POST")
	r.Handle("/site", middleware.Recovery(http.HandlerFunc(siteHandler.HandleGet))).Methods("GET")
	r.Handle("/site", middleware.Recovery(http.HandlerFunc(siteHandler.HandleDelete))).Methods("DELETE")

	return r
}
