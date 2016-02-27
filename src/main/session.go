package main

import (
	"encoding/gob"

	"github.com/gorilla/sessions"
)

type Cart struct {
	Selected_Developers Developers
	Total_Price         float64
}

type M map[string]interface{}

var store = sessions.NewCookieStore([]byte("biscoitos-muito-secretos"))

func init() {
	gob.Register(&Cart{})
	gob.Register(&M{})
}
