package main

import (
	"log"
	"net/http"

	"github.com/fern4lvarez/piladb/config/vars"
	"github.com/fern4lvarez/piladb/pila"
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

// checkMaxSizeOfStack checks config value for MaxSizeOfStack and execute the
// wrapped handler if check is validated.
func (c *Conn) checkMaxSizeOfStack(handler stackHandlerFunc) stackHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, stack *pila.Stack) {
		if s := c.Config.Get(vars.MaxSizeOfStack); stack.Size() == s && s != -1 {
			log.Println(r.Method, r.URL, http.StatusNotAcceptable, vars.MaxSizeOfStack, "value reached")
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}

		handler(w, r, stack)
	}
}
