package main

import "net/http"

// Route is a struct that contains the main parameters of an http route and the handler that handles it
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes is an array of Route
type Routes []Route

var routes = Routes{
	Route{
		"index",
		"GET",
		"/",
		index,
	},
	Route{
		"profiles",
		"POST",
		"/profiles",
		createProfileHandler,
	},
	Route{
		"recognize",
		"POST",
		"/recognize",
		recognizeHandler,
	},
}
