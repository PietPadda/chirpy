// main.go
package main

import (
	// std go libraries
	// for printing
	"database/sql"
	"log"      // for err logging
	"net/http" // http protocol
	"os"       // for io
	"strings"
	"time"

	// for conv itoa or atoi
	"sync/atomic" // allows safe incr + read of ints for goroutines

	// driver init
	"github.com/PietPadda/chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // postgresql driver
)

// STRUCTS
// stateful struct
type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	serverKey      string
}

// user database struct
type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

// MAIN
func main() {
	// load our .env file
	godotenv.Load() // loads it into the os environment for current program to use

	// get fields from .env file
	dbURL := os.Getenv("DB_URL")
	appPlatform := os.Getenv("PLATFORM")
	secretKey := strings.TrimSpace(os.Getenv("SECRET_KEY")) // remove whitespace from start and finish!
	// reaches into os env and gets the value at key

	// dbURL check
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}

	// Platform check
	if appPlatform == "" {
		log.Fatal("Platform is not set")
	}

	// secret key check
	if secretKey == "" {
		log.Fatal("SECRET_KEY is not set")
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
		db:             dbQueries,      // init the DBqueries for use in our handler
		platform:       appPlatform,    // init the platform for handler auth
		serverKey:      secretKey,      // init the server key for handler auth
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
	// register handlerReadiness, using /api/healthz system endpoint
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	// GET HTTP method routing only
	// healthz, because "system endpoint" convention!

	// register handlerMetrics, using /admin/metrics system endpoint
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics) // register func that receives apiCfg
	// GET HTTP method routing only
	// metrics, no z, as this is a conventional name!

	// register handlerAdminUsersReset, using /admin/reset system endpoint
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerAdminUsersReset) // register func that receives apiCfg
	// POST HTTP method routing only
	// reset, no z as this is a conventional name!

	// register handlerCreateChirp, using /api/chirps system endpoint
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp) // register func that receives apiCfg
	// POST HTTP method routing only

	// register handlerGetChirps, using /api/chirps system endpoint
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps) // register func that receives apiCfg
	// GET HTTP method routing only

	// register handlerGetChirp, using /api/chirps/{chirpID} system endpoint
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp) // register func that receives apiCfg
	// GET HTTP method routing only

	// register handlerCreateUser, using /api/users system endpoint
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser) // register func that receives apiCfg
	// POST HTTP method routing only

	// register handlerLoginUser, using /api/users system endpoint
	mux.HandleFunc("POST /api/login", apiCfg.handlerUserLogin) // register func that receives apiCfg
	// POST HTTP method routing only

	// register handlerRefresh, using /api/refresh system endpoint
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh) // register func that receives apiCfg
	// POST HTTP method routing only

	// register handlerRevoke, using /api/revoke system endpoint
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke) // register func that receives apiCfg
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
