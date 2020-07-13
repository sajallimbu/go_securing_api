package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sajallimbu/go_securing_api/controllers"
	"github.com/sajallimbu/go_securing_api/utils/auth"
)

//Handlers ... our route handling function
func Handlers() *mux.Router {
	router := mux.NewRouter()
	router.Use(CommonMiddleware)

	// convenience function that gives us the address of our controllers
	uc := controllers.NewUserController()

	// Define routes
	router.HandleFunc("/", http.HandlerFunc(uc.TestAPI)).Methods("GET")
	router.HandleFunc("/api", http.HandlerFunc(uc.TestAPI)).Methods("GET")
	router.HandleFunc("/register", http.HandlerFunc(uc.CreateUser)).Methods("POST")
	router.HandleFunc("/login", http.HandlerFunc(uc.Login)).Methods("POST")

	// The routes above doesnt need authentication as they are available to the public
	// but we need certain routes to have locks so that only authorized users can access the apis
	// adding a pathprefix makes the underlying url be localhost:8080/auth/user....
	superRouter := router.PathPrefix("/auth").Subrouter()
	// here we use our JwtVerify middleware that we create to check if the token is valid or not
	superRouter.Use(auth.JwtVerify)
	superRouter.HandleFunc("/user", http.HandlerFunc(uc.FetchUsers)).Methods("GET")
	superRouter.HandleFunc("/user/{id}", http.HandlerFunc(uc.GetUser)).Methods("GET")
	superRouter.HandleFunc("/user/{id}", http.HandlerFunc(uc.UpdateUser)).Methods("PUT")
	superRouter.HandleFunc("/user/{id}", http.HandlerFunc(uc.DeleteUser)).Methods("DELETE")
	return router
}

// CommonMiddleware ... a convenience middleware that set the commonly used headers for us
// you can change the headers if you want
func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		next.ServeHTTP(w, r)
	})
}
