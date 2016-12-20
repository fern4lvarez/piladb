package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/fern4lvarez/piladb/config/vars"
	"github.com/fern4lvarez/piladb/pila"
	"github.com/gorilla/mux"
)

// These vars represent the command line flags.
// They are only used to initialize the Connection
// Config at pilad start-up.
var (
	maxStackSizeFlag                  int
	readTimeoutFlag, writeTimeoutFlag int
	portFlag                          int
	versionFlag                       bool
)

func init() {
	flag.IntVar(&maxStackSizeFlag, "max-stack-size", vars.MaxStackSizeDefault, "Max size of Stacks")
	flag.IntVar(&readTimeoutFlag, "read-timeout", vars.ReadTimeoutDefault, "Read request timeout")
	flag.IntVar(&writeTimeoutFlag, "write-timeout", vars.WriteTimeoutDefault, "Write response timeout")
	flag.IntVar(&portFlag, "port", vars.PortDefault, "Port number")
	flag.BoolVar(&versionFlag, "v", false, "Version")
}

type flagKey struct {
	flag interface{}
	key  string
}

// buildConfig sets non-default config values to the Connection
// reading from environment variables and cli flags.
func (c *Conn) buildConfig() {
	flagKeys := []flagKey{
		{maxStackSizeFlag, vars.MaxStackSize},
		{readTimeoutFlag, vars.ReadTimeout},
		{writeTimeoutFlag, vars.WriteTimeout},
		{portFlag, vars.Port},
	}

	for _, fk := range flagKeys {
		if e := os.Getenv(vars.Env(fk.key)); e != "" {
			if i, err := strconv.Atoi(e); err != nil {
				c.Config.Set(fk.key, vars.DefaultInt(fk.key))
			} else {
				c.Config.Set(fk.key, i)

			}
			continue
		}
		c.Config.Set(fk.key, fk.flag)
	}
}

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
