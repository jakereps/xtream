package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

type wrIntercept struct {
	http.ResponseWriter
	status int
}

func (wr *wrIntercept) WriteHeader(statusCode int) {
	wr.status = statusCode
	wr.ResponseWriter.WriteHeader(statusCode)
}

func (wr *wrIntercept) Status() int {
	return wr.status
}

func main() {
	conf := flag.String("config", "/etc/xtream/xtream.conf", "stream data config path")
	flag.Parse()

	var err error
	defer func() {
		if err != nil {
			os.Exit(1)
		}
	}()

	sd, err = loadStreamData(*conf)
	if err != nil {
		log.Printf("failed reading config from: %s", *conf)
		return
	}

	// kiss routing
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wi := &wrIntercept{w, http.StatusOK}
		switch r.URL.Path {
		case "/":
			root(wi, r)
		case "/authz":
			handleAuthz(wi, r)
		default:
			http.Error(wi, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
		log.Printf("handled request %s - %s (%d - %s)", r.Method, r.URL.Path, wi.Status(), http.StatusText(wi.Status()))
	})

	server := http.Server{
		Addr:    "localhost:8000",
		Handler: handler,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Printf("server error: %s", err)
		return
	}
}

func root(w http.ResponseWriter, r *http.Request) {

}

func handleAuthz(w http.ResponseWriter, r *http.Request) {
	app := r.URL.Query().Get("app")
	if app == "" {
		log.Println("missing app name")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	sk := r.URL.Query().Get("name")
	if sk == "" {
		log.Println("missing stream key")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	appData, ok := sd.App(sk)
	if !ok {
		log.Println("invalid sk provided")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	if app != appData.Name {
		http.Redirect(w, r, fmt.Sprintf("rtmp://xtream/%s", appData.Name), http.StatusFound)
		return
	}
}
