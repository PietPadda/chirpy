// main.go
package main

import (
	// std go libraries
	// for printing
	"database/sql"
	"log"      // for err logging
	"net/http" // http protocol
	"os"       // for io

	// for conv itoa or atoi
	"sync/atomic" // allows safe incr + read of ints for goroutines

	// driver init
	"github.com/PietPadda/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // postgresql driver
)

// STRUCTS
// stateful struct
type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

// MAIN
func main() {
	// load our .env file
	godotenv.Load() // loads it into the os environment for current program to use

	// get db_url for .env file
	dbURL := os.Getenv("DB_URL") // reaches into os env and gets the value at key

	// dbURL check
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}

	// open connection to your database using the DBUrl and driver
	db, err := sql.Open("postgres", dbURL)

	// db connection check
	if err != nil {
		log.Fatal("error connecting to db:", err)
	}

	// use SQLC database package
	dbQueries := database.New(db)

	// set constants
	const filepathRoot = "." // used constant
	const port = "8080"

	// create server mux for routing http requests
	mux := http.NewServeMux()

	// create apiConfig instance
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{}, // explicitly set to 0
		db:             dbQueries,      // init the dbqueries for use in our handler
	}

	// create the file server handle
	fsHandler := apiCfg.middlewareMetricsInc(
		http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot))),
	)
	// http.Dir(".") -  tells server to run files "here"
	// http.FileServer(...) - looks for the index.html
	// stripping the /app/ to just . -- /app/ is just there for cleaner url

	// apply a fileserver handler to mux
	mux.Handle("/app/", fsHandler)
	// mux.Handle("/app/", ...) -- server handle all requests

	// REGISTER HANDLERS
	// register handlerReadiness, using api/healthz system endpoint
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	// GET HTTP method routing only
	// healthz, because "system endpoint" convention!

	// register handlerMetrics, using admin/metrics system endpoint
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics) // register func that receives apiCfg
	// GET HTTP method routing only
	// metrics, no z, as this is a conventional name!

	// register handlerMetricsReset, using admin/reset system endpoint
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerMetricsReset) // register func that receives apiCfg
	// POST HTTP method routing only
	// reset, no z as this is a conventional name!

	// register handlerValidateChirp, using api/validate_chirps system endpoint
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	// POST HTTP method routing only

	// create Server struct for config
	server := &http.Server{ //ptr is more efficient than new copy
		Addr:    ":" + port, //server listens to port 8080 for all incoming requests
		Handler: mux,        // mux will "handle" our http request
	}

	// print server running msg
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	// print msg before blocking with "Listen"

	// start the server
	err = server.ListenAndServe() // "listens" to http requests on addr and let's mux handle them
	// listen blocks the server to prevent ending main func

	// server start check
	if err != nil {
		// log the err and terminate server
		log.Fatal("error starting server:", err)
	}
}
