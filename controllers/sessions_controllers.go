package controllers

import (
	"net/http"
	"html/template"
	"appengine"
	"code.google.com/p/goauth2/oauth"
	"appengine/urlfetch"
	"encoding/json"
	"models"
	"io/ioutil"
	"helpers"
	"fmt"
)

// Set up a configuration.
var config = &oauth.Config{
		ClientId:     "58551267124.apps.googleusercontent.com",
		ClientSecret: "WsXC7Ea1LKWdz0KoQhE-yeQG",
		RedirectURL:  "http://localhost:8080/signin_google",
		Scope:        "https://www.googleapis.com/auth/userinfo.profile https://www.googleapis.com/auth/userinfo.email",
		AuthURL:      "https://accounts.google.com/o/oauth2/auth",
		TokenURL:     "https://accounts.google.com/o/oauth2/token",
		TokenCache:   nil,
}

var currentUser *models.User

func Signin(w http.ResponseWriter, r *http.Request) {
	if !helpers.SignedIn(r) {
		renderSigninPage(w)
	} else {
		// redirect to home
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
}

func renderSigninPage(w http.ResponseWriter) {
	t, _ := template.ParseFiles("app/templates/header.html", 
								"app/templates/signin.html",
								"app/templates/footer.html")

	if err := t.ExecuteTemplate(w, "tmpl_signin", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func authorize(w http.ResponseWriter, r *http.Request, t *oauth.Transport) (*oauth.Token, error) {
	//Get the code from the response
	code := r.FormValue("code")
	
	if code == "" {
		// Get an authorization code from the data provider.
		// ("Please ask the user if I can access this resource.")
		url := config.AuthCodeURL("")
		http.Redirect(w, r, url, http.StatusFound)
		return nil, nil
	}
	// Exchange the authorization code for an access token.
	// ("Here's the code you gave the user, now give me a token!")
	return t.Exchange(code)
}

func getCurrentUser(r *http.Request, t *oauth.Transport) (*models.User, error) {
	// Make the request.
	request, err := t.Client().Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
	
	if err != nil {
		return nil, err
	}
	
	if userInfo, err := ioutil.ReadAll(request.Body); err == nil {
		var gu models.GoogleUser
		
		if err := json.Unmarshal(userInfo, &gu); err == nil {
			// create new user if he doesn't already exist
			if currentUser, err = models.GetUser(r, gu.Id, gu.Email); err != nil {
				currentUser, err = models.NewUser(r, gu.Id, gu.Email, gu.Name, "google")
			}
			
			return currentUser, err
		}	
	}
	
	return nil, err
}

func SigninWithGoogle(w http.ResponseWriter, r *http.Request) {
	if !helpers.SignedIn(r) {
		c := appengine.NewContext(r)
		
		// Set up a Transport using the config.
		transport := &oauth.Transport{Config: config, Transport: &urlfetch.Transport{Context: c}}
		
		// get the access token
		token, err := authorize(w, r, transport)
		if err != nil {
			c.Debugf("%q", err)
		}
		
		// store token to authenticate
		transport.Token = token
		
		if currentUser, _ := getCurrentUser(r, transport); currentUser != nil {
			//save the session
			helpers.SignIn(w, fmt.Sprintf("%d", currentUser.Id))	
		}
	}
	
	// redirect to home
	http.Redirect(w, r, "/", http.StatusFound)
	return
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// clear the session
	helpers.Logout(w)
	// clear the current user
	currentUser = nil
	// redirect to home
	http.Redirect(w, r, "/", http.StatusFound)
	return
}