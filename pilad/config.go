package main

import (
	"log"
	"net/http"

	"github.com/fern4lvarez/piladb/config/vars"
	"github.com/fern4lvarez/piladb/pila"
)

// stackHandlerFunc represents a Handler of a Stack.
type stackHandlerFunc func(w http.ResponseWriter, r *http.Request, stack *pila.Stack)

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
