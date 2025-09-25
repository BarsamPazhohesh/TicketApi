// Package routes
package routes

import "strings"

type _Prefix struct {
	prefix string
}

type _APIPrefixes struct {
	Versions _Prefix
	Tickets  _Prefix
	Auth     _Prefix
	Captcha  _Prefix
}

var _APIRoutesPrefixes = _APIPrefixes{
	Tickets: _Prefix{prefix: "tickets/"},
	Auth:    _Prefix{prefix: "auth/"},
	Captcha: _Prefix{prefix: "captcha/"},
}

type HTTPMethod string

const (
	GetMethod    HTTPMethod = "GET"
	PostMethod   HTTPMethod = "POST"
	PutMethod    HTTPMethod = "PUT"
	DeleteMethod HTTPMethod = "DELETE"
)

type _APIRoute struct {
	Path        string
	method      string
	description string
}

type versions struct {
	GetCurrentVersion _APIRoute
}

type tickets struct {
	CreateTicket         _APIRoute
	GetTicketByID        _APIRoute
	GetTicketByTrackCode _APIRoute
	CreateChat           _APIRoute
	GetTicketsList       _APIRoute
}

type auth struct {
	LoginWithNoAuth         _APIRoute
	SignUp                  _APIRoute
	Login                   _APIRoute
	GetSingleUseToken       _APIRoute
	LoginWithSingleUseToken _APIRoute
}
type captcha struct {
	GetCaptcha    _APIRoute
	VerifyCaptcha _APIRoute
}

type _APIEndpoints struct {
	Versions versions
	Tickets  tickets
	Auth     auth
	Captcha  captcha
}

var APIRoutes = _APIEndpoints{
	Versions: versions{
		GetCurrentVersion: _APIRoute{Path: "", method: string(GetMethod)},
	},
	Tickets: tickets{
		CreateTicket:         _APIRoute{Path: mergeStrings(_APIRoutesPrefixes.Tickets.prefix, "CreateTicket/"), method: string(PostMethod)},
		GetTicketByID:        _APIRoute{Path: mergeStrings(_APIRoutesPrefixes.Tickets.prefix, ":id/"), method: string(GetMethod)},
		CreateChat:           _APIRoute{Path: mergeStrings(_APIRoutesPrefixes.Tickets.prefix, ":id/CreateChat/"), method: string(PostMethod)},
		GetTicketByTrackCode: _APIRoute{Path: mergeStrings(_APIRoutesPrefixes.Tickets.prefix, "GetTicketByTrackCode/"), method: string(PostMethod)},
		GetTicketsList:       _APIRoute{Path: mergeStrings(_APIRoutesPrefixes.Tickets.prefix, "GetTicketsList/"), method: string(PostMethod)},
	},
	Auth: auth{
		LoginWithNoAuth:         _APIRoute{Path: mergeStrings(_APIRoutesPrefixes.Auth.prefix, "LoginWithNoAuth/"), method: string(GetMethod)},
		SignUp:                  _APIRoute{Path: mergeStrings(_APIRoutesPrefixes.Auth.prefix, "SignUp/"), method: string(PostMethod)},
		Login:                   _APIRoute{Path: mergeStrings(_APIRoutesPrefixes.Auth.prefix, "Login/"), method: string(GetMethod)},
		GetSingleUseToken:       _APIRoute{Path: mergeStrings(_APIRoutesPrefixes.Auth.prefix, "GetSingleUseToken/"), method: string(PostMethod)},
		LoginWithSingleUseToken: _APIRoute{Path: mergeStrings(_APIRoutesPrefixes.Auth.prefix, "LoginWithSingleUseToken/"), method: string(GetMethod)},
	},
	Captcha: captcha{
		GetCaptcha:    _APIRoute{Path: mergeStrings(_APIRoutesPrefixes.Captcha.prefix, "GetCaptcha/"), method: string(GetMethod)},
		VerifyCaptcha: _APIRoute{Path: mergeStrings(_APIRoutesPrefixes.Captcha.prefix, "VerifyCaptcha/"), method: string(PostMethod)},
	},
}

func mergeStrings(string ...string) string {
	var builder strings.Builder
	for _, s := range string {
		builder.WriteString(s)
	}
	return builder.String()
}
