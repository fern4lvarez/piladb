package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fern4lvarez/piladb/config/vars"
	"github.com/fern4lvarez/piladb/pila"
	"github.com/gorilla/mux"
)

// stackHandlerFunc represents a Handler of a Stack.
type stackHandlerFunc func(w http.ResponseWriter, r *http.Request, stack *pila.Stack)

// configHandler handles a request to the Conn configuration.
func (c *Conn) configHandler(w http.ResponseWriter, r *http.Request) {
	res, err := c.Config.Values.StacksKV().ToJSON()
	if err != nil {
		log.Println(r.Method, r.URL, http.StatusBadRequest,
			"error on response serialization:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
	log.Println(r.Method, r.URL, http.StatusOK)
}

// configKeyHandler handles a config value.
func (c *Conn) configKeyHandler(configKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		// we override the mux vars to be able to test
		// an arbitrary configKey
		if configKey != "" {
			vars = map[string]string{
				"key": configKey,
			}
		}
		value := c.Config.Get(vars["key"])
		if value == nil {
			c.goneHandler(w, r, fmt.Sprintf("%s is not set", vars["key"]))
			return
		}

		var element pila.Element
		if r.Method == "GET" {
			value := c.Config.Get(vars["key"])
			element.Value = value
		}
		if r.Method == "POST" {
			if r.Body == nil {
				log.Println(r.Method, r.URL, http.StatusBadRequest,
					"no element provided")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			err := element.Decode(r.Body)
			if err != nil {
				log.Println(r.Method, r.URL, http.StatusBadRequest,
					"error on decoding element:", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			c.Config.Set(vars["key"], element.Value)
		}

		log.Println(r.Method, r.URL, http.StatusOK, element.Value)
		w.Header().Set("Content-Type", "application/json")

		b, err := element.ToJSON()
		if err != nil {
			log.Println(r.Method, r.URL, http.StatusBadRequest,
				"error on decoding element:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Write(b)
	})
}

// checkMaxStackSize checks config value for MaxStackSize and execute the
// wrapped handler if check is validated.
func (c *Conn) checkMaxStackSize(handler stackHandlerFunc) stackHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, stack *pila.Stack) {
		if s := c.Config.MaxStackSize(); stack.Size() >= s && s != -1 {
			log.Println(r.Method, r.URL, http.StatusNotAcceptable, vars.MaxStackSize, "value reached")
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}

		handler(w, r, stack)
	}
}
