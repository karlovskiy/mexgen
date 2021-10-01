package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	_ "github.com/heroku/x/hmetrics/onload"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"time"
)

func main() {
	portStr := os.Getenv("PORT")
	if portStr == "" {
		log.Fatal("PORT must be set")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("bad PORT: %s: %v", portStr, err)
	}
	log.Printf("PORT=%d", port)
	log.Printf("WKHTMLTOPDF_PATH=%s", os.Getenv("WKHTMLTOPDF_PATH"))
	rand.Seed(time.Now().UnixNano())
	mux := http.NewServeMux()

	mux.Handle("/multdiv",
		handlers.LoggingHandler(os.Stdout, Handler("GET", http.HandlerFunc(multDivHandler))))

	mux.Handle("/addsub",
		handlers.LoggingHandler(os.Stdout, Handler("GET", http.HandlerFunc(addSubHandler))))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}

func multDivHandler(w http.ResponseWriter, r *http.Request) {
	generateHandler(w, "multdiv", &MultDiv{
		Rows: 30,
		Cols: 8,
	})
}

func addSubHandler(w http.ResponseWriter, r *http.Request) {
	generateHandler(w, "addsub", &AddSub{
		Rows: 15,
		Cols: 16,
	})
}

func generateHandler(w http.ResponseWriter, name string, generator ExampleGenerator) {
	w.Header().Set("Content-Type", "application/pdf")
	filename := fmt.Sprintf("attachment; filename=%s-%s.pdf", name, time.Now().UTC().Format("20060102150405"))
	w.Header().Set("Content-Disposition", filename)
	if err := GeneratePdf(w, generator); err != nil {
		log.Printf("error generating %s: %v", filename, err)
		httpError(w, http.StatusInternalServerError)
	}
}

func Handler(method string, h http.Handler) http.Handler {
	return handler{method, h}
}

type handler struct {
	method  string
	handler http.Handler
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqDump, err := httputil.DumpRequest(r, false)
	if err != nil {
		log.Printf("error dumping req: %v\n", err)
	} else {
		log.Printf("\n%s\n", reqDump)
	}
	if r.Method != h.method {
		httpError(w, http.StatusMethodNotAllowed) // 405
		return
	}
	h.handler.ServeHTTP(w, r)
}

func httpError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
