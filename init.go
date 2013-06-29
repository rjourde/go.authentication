package main 

import (
	"net/http"
	"controllers"
)

func init() {
	http.HandleFunc("/", controllers.Index)
	http.HandleFunc("/signin_google", controllers.SigninWithGoogle)
	http.HandleFunc("/logout", controllers.Logout)
}