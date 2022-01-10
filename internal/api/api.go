package api

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/deesel/wol"
	"github.com/deesel/wol/internal/config"
	"github.com/gorilla/mux"
)

// API holds configuration for API server
type API struct {
	Address string
	Port    int
	Auth    bool
	APIKeys []Key
}

// New creates new API server instance
func New(c *config.Config) *API {
	api := &API{
		Address: c.Server.Address,
		Port:    c.Server.Port,
		Auth:    c.Auth.Enabled,
		APIKeys: []Key{},
	}

	for _, key := range c.Auth.APIKeys {
		api.APIKeys = append(api.APIKeys, Key(key.Key))
	}

	return api
}

// Run runs API server instance
func (a *API) Run() error {
	router := mux.NewRouter().StrictSlash(true)

	router.Use(accessLog)
	router.Use(jsonContent)
	router.Use(handleError)

	if a.Auth {
		router.Use(a.isAuthorized)
	}

	router.HandleFunc("/", a.handleWOL).Methods("POST")

	return http.ListenAndServe(fmt.Sprintf("%s:%d", a.Address, a.Port), router)
}

func (a *API) handleWOL(wr http.ResponseWriter, req *http.Request) {
	var data struct {
		Type      wol.Type
		MAC       string
		IP        string
		Port      int
		Interface string
	}

	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	if data.Type == "" {
		data.Type = wol.UDP
	}

	if data.IP == "" {
		data.IP = "255.255.255.255"
	}

	if data.Port == 0 {
		data.Port = 9
	}

	mac, err := net.ParseMAC(data.MAC)
	if err != nil {
		panic(err)
	}

	var w *wol.WOL

	switch data.Type {
	case wol.Ethernet:
		w, err = wol.NewEther(mac, data.Interface)
	case wol.UDP:
		w, err = wol.NewUDP(mac, net.ParseIP(data.IP), data.Port)
	default:
		panic(ErrUnknownWOLType)
	}

	if err != nil {
		panic(err)
	}

	err = w.Send()
	if err != nil {
		panic(err)
	}
}
