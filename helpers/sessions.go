package helpers

import (
	"net/http"
	"github.com/gorilla/securecookie"
)

var secret []byte = securecookie.GenerateRandomKey(32)
var userIdCookie *securecookie.SecureCookie = securecookie.New(secret, nil)
const cookieName string = "user_id"

func SignIn(w http.ResponseWriter, cookieValue string) {
	value := map[string]string{
		cookieName : cookieValue,
	}
	if encoded, err := userIdCookie.Encode(cookieName, value); err == nil {
		cookie := &http.Cookie{
			Name:  cookieName,
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}
}

func SignedIn(r *http.Request) bool {
	
	if cookie, err := r.Cookie(cookieName); err == nil {
		value := make(map[string]string)
		if(cookie != nil) {
			err = userIdCookie.Decode(cookieName, cookie.Value, &value)
			if (len(value["user_id"]) > 0 && err == nil) {
				return true
			}
		}
	}
	
	return false
}

func Logout(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:  cookieName,
		Value: "",
		Path:  "/",
	}
	http.SetCookie(w, cookie)
}