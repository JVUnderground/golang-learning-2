package main

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"showAllDevelopers",
		"GET",
		"/developers",
		showAllDevelopers,
	},
	Route{
		"Developer",
		"GET",
		"/developers/{devId}",
		showDeveloper,
	},
	Route{
		"updateCart",
		"POST",
		"/update-cart",
		updateCart,
	},
	Route{
		"removeFromCart",
		"POST",
		"/remove-from-cart",
		removeFromCart,
	},
	Route{
		"showCart",
		"GET",
		"/check-out",
		showCart,
	},
	Route{
		"confirmPurchase",
		"POST",
		"/confirm-purchase",
		confirmPurchase,
	},
}
