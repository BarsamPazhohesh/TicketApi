// Package routes
package routes

type _APIRoute struct {
	Path        string
	method      string
	description string
}

type versions struct {
	GetCurrentVersion _APIRoute
}

type tickets struct {
	CreateTicket  _APIRoute
	GetTicketByID _APIRoute
	CreateChat    _APIRoute
}

type auth struct {
	LoginWithNoAuth _APIRoute
}

type _APIEndpoints struct {
	Versions versions
	Tickets  tickets
	Auth     auth
}

var APIRoutes = _APIEndpoints{
	Versions: versions{
		GetCurrentVersion: _APIRoute{Path: "/", method: "Get"},
	},
	Tickets: tickets{
		CreateTicket:  _APIRoute{Path: "/tickets", method: "POST"},
		GetTicketByID: _APIRoute{Path: "/tickets/:id", method: "GET"},
		CreateChat:    _APIRoute{Path: "/tickets/:id/chat", method: "POST"},
	},
	Auth: auth{
		LoginWithNoAuth: _APIRoute{Path: "/auth/LoginWithNoAuth", method: "POST"},
	},
}
