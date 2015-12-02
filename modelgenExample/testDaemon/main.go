package main

import (
	"net/http"

	contivModel "github.com/contiv/objmodel/tools/modelgen/modelgenExample"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

type ObjHandler struct {
	dummy string
}

func (self *ObjHandler) TenantCreate(tenant *contivModel.Tenant) error {
	log.Infof("Creating tenant: %+v", tenant)
	return nil
}

func (self *ObjHandler) TenantUpdate(tenant, params *contivModel.Tenant) error {
	log.Infof("Received TenantUpdate: %+v, params: %+v", tenant, params)
	return nil
}

func (self *ObjHandler) TenantDelete(tenant *contivModel.Tenant) error {
	log.Infof("Deleting tenant: %+v", tenant)
	return nil
}

func (self *ObjHandler) NetworkCreate(network *contivModel.Network) error {
	log.Infof("Creating network: %+v", network)
	return nil
}

func (self *ObjHandler) NetworkUpdate(network, params *contivModel.Network) error {
	log.Infof("Received NetworkUpdate: %+v, params: %+v", network, params)
	return nil
}

func (self *ObjHandler) NetworkDelete(network *contivModel.Network) error {
	log.Infof("Deleting network: %+v", network)
	return nil
}

var handler ObjHandler

func main() {
	// Initialize the models
	contivModel.Init()

	// Register Callbacks
	contivModel.RegisterNetworkCallbacks(&handler)
	contivModel.RegisterTenantCallbacks(&handler)

	// Create a router
	router := mux.NewRouter()

	// Register routes
	contivModel.AddRoutes(router)

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8000", router))
}
