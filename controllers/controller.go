package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sajallimbu/go_securing_api/utils"

	"github.com/gorilla/mux"

	"github.com/dgrijalva/jwt-go"

	"golang.org/x/crypto/bcrypt"

	"github.com/sajallimbu/go_securing_api/models"
)

// UserController ... struct UserController
type UserController struct{}

var db = utils.ConnectDB()

// NewUserController ... returns the address of UserController Interface
func NewUserController() *UserController {
	return &UserController{}
}

// TestAPI ... should return confirmation that the api routes are working
func (uc UserController) TestAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "API running successfully")
}

// CreateUser ... create user function
func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request) {

	user := &models.User{}

	//Parse the request body for incoming user data
	json.NewDecoder(r.Body).Decode(user)

	//generate a hashkey for the user using his password
	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		json.NewEncoder(w).Encode(fmt.Sprintf("Password encryption failed: %s", err))
	}

	//set the generated hashkey as the password for the user
	user.Password = string(pass)

	//create a new user using the above credentials
	createdUser := db.Create(user)
	var errMessage = createdUser.Error

	//check error when creating the user
	if errMessage != nil {
		fmt.Fprintf(w, "Error creating user: %s", errMessage)
	}
	//send the newly created user credentials to the request on success
	json.NewEncoder(w).Encode(createdUser)
}

// Login ... login function
func (uc UserController) Login(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Invalid request"}
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp := FindUser(user.Email, user.Password)
	json.NewEncoder(w).Encode(resp)
}

// FindUser ... our authentication function that takes the user email and password then returns a response
func FindUser(email, password string) map[string]interface{} {
	user := &models.User{}

	// Query the database for matching email address and store the matched result in user
	if err := db.Where("Email = ?", email).First(user).Error; err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Email address not found"}
		return resp
	}

	// Add an expiry window to the JWT token
	expiresAt := time.Now().Add(time.Minute * 100000).Unix()

	// Compare the hash stored in the database with the password inputted by the user
	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match
		var resp = map[string]interface{}{"status": false, "message": "Incorrect password"}
		return resp
	}

	// Create a token model for signing and set an expiry claim
	tk := &models.Token{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	// Set the signing method
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	// Sign the token with your own unique secret
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		fmt.Println(err)
	}

	var resp = map[string]interface{}{"status": false, "message": "logged in"}
	// Add a couple of new key value pair in resp
	resp["token"] = tokenString
	resp["user"] = user

	return resp
}

// FetchUsers ... function that returns all users
func (uc UserController) FetchUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	db.Find(&users)
	json.NewEncoder(w).Encode(&users)
}

// GetUser ... function that returns the user data
func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	params := mux.Vars(r)
	var id = params["id"]
	db.First(&user, id)
	json.NewEncoder(w).Encode(&user)
}

// UpdateUser ... function that updates the user data
func (uc UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	params := mux.Vars(r)
	var id = params["id"]
	db.First(&user, id)
	json.NewDecoder(r.Body).Decode(&user)
	db.Save(&user)
	json.NewEncoder(w).Encode(&user)
}

// DeleteUser ... function that deletes a user
func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var id = params["id"]
	var user models.User
	db.First(&user, id)
	db.Delete(&user)
	json.NewEncoder(w).Encode("User deleted")
}