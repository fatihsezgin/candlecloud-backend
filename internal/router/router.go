package router

import (
	"net/http"

	"github.com/fatihsezgin/candlecloud-backend/internal/api"
	"github.com/fatihsezgin/candlecloud-backend/internal/storage"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type Router struct {
	router *mux.Router
	store  storage.Store
}

func New(s storage.Store) *Router {
	r := &Router{
		router: mux.NewRouter(),
		store:  s,
	}
	r.initRoutes()
	return r
}

// ServeHTTP ...
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

func (r *Router) initRoutes() {

	// API Router Group
	apiRouter := mux.NewRouter().PathPrefix("/api").Subrouter()

	// User endpoints
	apiRouter.HandleFunc("/users", api.FindAllUsers(r.store)).Methods(http.MethodGet)
	apiRouter.HandleFunc("/users", api.CreateUser(r.store)).Methods(http.MethodPost)

	// Product endpoints
	apiRouter.HandleFunc("/products", api.CreateProduct(r.store)).Methods(http.MethodPost)
	apiRouter.HandleFunc("/products", api.GetAllProducts(r.store)).Methods(http.MethodGet)

	// Auth endpoints
	authRouter := mux.NewRouter().PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/signup", api.Signup(r.store)).Methods(http.MethodPost)

	n := negroni.Classic()
	n.Use(negroni.HandlerFunc(CORS))

	// TODO refactor after the token issues has done.
	r.router.PathPrefix("/api").Handler(n.With(
		LimitHandler(),
		negroni.Wrap(apiRouter),
	))

	r.router.PathPrefix("/auth").Handler(n.With(
		LimitHandler(),
		negroni.Wrap(authRouter),
	))
}
