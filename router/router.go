package router

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"morph-tracking-api/database"
	"net/http"
)

type Deps struct {
	db *database.AxieDB
}

func Listen(db *database.AxieDB) {
	r := setupRouter(db)
	cors := setupCORS()
	// Start server
	fmt.Println("Starting server at port 3000")
	log.Fatal(http.ListenAndServe(":3000", cors(r)))
}

// setupRouter Sets up the router along with the http handlers of the API
func setupRouter(db *database.AxieDB) *mux.Router {
	deps := Deps{db}
	router := mux.NewRouter()
	router.Handle("/favicon.ico", http.NotFoundHandler())
	router.HandleFunc("/", deps.MorphHandler)
	return router
}

// SetupCORS Sets up the CORS policy for the API
func setupCORS() func(http.Handler) http.Handler {
	return handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET"}),
		handlers.AllowedOrigins([]string{"*"}),
	)
}
