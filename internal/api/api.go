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

type API struct {
	Address string
	Port    int
	Auth    bool
	APIKeys []APIKey
}

func New(c *config.Config) *API {
	api := &API{
		Address: c.Server.Address,
		Port:    c.Server.Port,
		Auth:    c.Auth.Enabled,
		APIKeys: []APIKey{},
	}

	for _, key := range c.Auth.APIKeys {
		api.APIKeys = append(api.APIKeys, APIKey(key.Key))
	}

	return api
}

func (a *API) Run() {
	router := mux.NewRouter().StrictSlash(true)

	router.Use(accessLog)
	router.Use(jsonContent)
	router.Use(handleError)

	if a.Auth {
		router.Use(a.isAuthorized)
	}

	router.HandleFunc("/", a.handleWOL).Methods("POST")

	http.ListenAndServe(fmt.Sprintf("%s:%d", a.Address, a.Port), router)
}

func (a *API) handleWOL(wr http.ResponseWriter, req *http.Request) {
	var data struct {
		Type      wol.WOLType
		MAC       string
		IP        string
		Port      int
		Interface string
	}

	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		panic(err)
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
