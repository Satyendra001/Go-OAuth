package user

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/markbates/goth/gothic"
	"github.com/satyendra001/mdm-oauth/auth"
)

// Create a function to take handler as a receiver and register some routes
func (h *UserStore) RegisterOauthRoutes(router *mux.Router) {
	router.HandleFunc("/test", h.testRoute).Methods("GET")                              // Test Route
	router.HandleFunc("/auth/{provider}", h.HandleAuthProvider).Methods("GET")          // Handle the auth provider call
	router.HandleFunc("/auth/{provider}/callback", h.handleAuthCallback).Methods("GET") // Handle the auth callback from the Oauth provider

}

func (h *UserStore) testRoute(w http.ResponseWriter, r *http.Request) {
	log.Println("Test route called")
	h.GetAllUserInfo()

	fmt.Fprint(w, "Test Successfull!!!")
}

func (h *UserStore) HandleAuthProvider(w http.ResponseWriter, r *http.Request) {

	// Start a new goth service which will handle the specified auth providers
	auth.NewAuth()

	// try to get the user without re-authenticating
	if gothUser, err := gothic.CompleteUserAuth(w, r); err == nil {
		log.Println("User Auth completed...", gothUser)
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}

func (h *UserStore) handleAuthCallback(w http.ResponseWriter, r *http.Request) {
	log.Println("callback Hit reveived from the Oauth provider...")
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	user.Email = "satyendra@infilect.com"

	fmt.Println("Extracted User details ==>", user.Description, user.UserID, user.Email)
	// 1. Query the DB in public schema to check if the user exists?
	log.Println("User email ==> ", user.Email)

	var userId int

	dbUser, err := h.GetUserByEmail(user.Email)
	if err != nil {

		//2. Create a new User
		log.Println("User not found with email ", user.Email, ". Creating a new User...")

		userId, err = h.CreateUser(&user)

		if err != nil {
			log.Fatal("unable to create user. Error ==>", err.Error())
		}

	} else {
		userId = dbUser.Id
	}

	// 3. Get the token for the user
	token := h.GetToken(userId)

	// 4. Create token for the user if it is empty
	if token == "" {
		token, err = h.CreateToken(userId)

		if err != nil {
			log.Fatal("Error in new token creation. Err -->", err.Error())
		}
	}

	log.Println("Setting the token in Cookie", token)
	cookie := &http.Cookie{
		Name:     "db_token",
		Value:    token, // Plaintext value
		Path:     "/",   // Cookie path
		Domain:   ".dmt.localhost",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}

	http.SetCookie(w, cookie)

	// utils.WriteJSON(w, 200, map[string]string{"token": token})
	http.Redirect(w, r, "http://dmt.localhost:5173/dashboard", http.StatusSeeOther)
	// http.Redirect(w, r, "http://localhost:5173/home", http.StatusSeeOther)
}
