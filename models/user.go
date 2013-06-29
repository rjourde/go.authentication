package models

import (
	"net/http"
	"time"
	"appengine"
	"appengine/datastore"
	"errors"
)

type User struct {
	Id int64
	Email string
	Username string
	Created time.Time
}

type ExternalAuth struct {
	Id int64
	ExternalId string
	UserId int64
	Provider string
}

type GoogleUser struct {
	Id string
	Email string
	Name string
	GivenName string
	FamilyName string
}

func NewUser(r *http.Request, id string, email string, name string, provider string) (*User, error) {
	c := appengine.NewContext(r)
	// create new user
	userId, _, _ := datastore.AllocateIDs(c, "User", nil, 1)
	key := datastore.NewKey(c, "User", "", userId, nil)
	
	user := User{ userId, email, name, time.Now() }
	
	_, err := datastore.Put(c, key, &user)
	if err != nil {
		return nil, err;
	}
	// create external authentication
	externalAuthId, _, _ := datastore.AllocateIDs(c, "ExternalAuth", nil, 1)
	key = datastore.NewKey(c, "ExternalAuth", "", externalAuthId, nil)
	
	externalAuth := ExternalAuth{ externalAuthId, id, userId, provider }
	
	_, err = datastore.Put(c, key, &externalAuth)
	if err != nil {
		return nil, err;
	}
	
	return &user, err;
}

func GetUser(r *http.Request, id string, email string) (*User, error) {
	var user *User
	
	c := appengine.NewContext(r)
	// Fetch the page
	q:= datastore.NewQuery("User").Filter("Email =", email)
	
	var users []*User
	keys, err := q.GetAll(c, &users)
	if err != nil {
		return nil, err
	}
	if(keys == nil || users == nil) {
		return nil, errors.New(email + " doesn't exit")
	}
	
	user = users[0]
	
	return user, nil
}