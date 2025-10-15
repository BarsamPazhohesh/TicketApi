// Package routes
package routes

import (
	"strings"
)

type _Prefix struct {
	prefix string
}

type _APIPrefixes struct {
	Versions   _Prefix
	Tickets    _Prefix
	Auth       _Prefix
	Captcha    _Prefix
	User       _Prefix
	Department _Prefix
	Files      _Prefix
}

var _APIRoutesPrefixes = _APIPrefixes{
	Tickets:    _Prefix{prefix: "tickets/"},
	Auth:       _Prefix{prefix: "auth/"},
	Captcha:    _Prefix{prefix: "captcha/"},
	User:       _Prefix{prefix: "users/"},
	Department: _Prefix{prefix: "departments/"},
	Files:      _Prefix{prefix: "files/"},
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
	Status      bool
}

type versions struct {
	GetCurrentVersion _APIRoute
}

type tickets struct {
	CreateTicket               _APIRoute
	GetTicketByID              _APIRoute
	GetTicketByTrackCode       _APIRoute
	CreateChat                 _APIRoute
	GetTicketsList             _APIRoute
	GetAllActiveTicketTypes    _APIRoute
	GetAllActiveTicketStatuses _APIRoute
}

type departments struct {
	GetAllActiveDepartments _APIRoute
}

type auth struct {
	LoginWithNoAuth         _APIRoute
	SignUp                  _APIRoute
	Login                   _APIRoute
	GetSingleUseToken       _APIRoute
	LoginWithSingleUseToken _APIRoute
}

type users struct {
	GetUserByUsername _APIRoute
	GetUserByID       _APIRoute
	GetUsersByIDs     _APIRoute
}

type captcha struct {
	GetCaptcha    _APIRoute
	VerifyCaptcha _APIRoute
}

type files struct {
	UploadTicketFile   _APIRoute
	DownloadTicketFile _APIRoute
}
type _APIEndpoints struct {
	Versions    versions
	Tickets     tickets
	Files       files
	Auth        auth
	Captcha     captcha
	Users       users
	Departments departments
}

var APIRoutes = _APIEndpoints{
	Versions: versions{
		GetCurrentVersion: _APIRoute{Path: "", method: string(GetMethod)},
	},
	Tickets: tickets{
		CreateTicket:               _APIRoute{Path: mergeStrings(_APIRoutesPrefixes.Tickets.prefix, "CreateTicket/"), method: string(PostMethod), Status: true},
		GetTicketByID:              _APIRoute{Path: mergeStrings(_APIRoutesPrefixes.Tickets.prefix, "GetTicketByID/"), method: string(GetMethod), Status: true},
		CreateChat:                 _APIRoute{Path: mergeStrings(_APIRoutesPrefixes.Tickets.prefix, ":id/CreateChat/"), method: string(PostMethod), Status: true},
		GetTicketByTrackCode:       _APIRoute{Path: mergeStrings(_APIRoutesPrefixes.Tickets.prefix, "GetTicketByTrackCode/"), method: string(PostMethod), Status: true},
		GetTicketsList:             _APIRoute{Path: mergeStrings(_APIRoutesPrefixes.Tickets.prefix, "GetTicketsList/"), method: string(PostMethod), Status: true},
		GetAllActiveTicketTypes:    _APIRoute{Path: mergeStrings(_APIRoutesPrefixes.Tickets.prefix, "GetAllActiveTicketTypes/"), method: string(GetMethod), Status: true},
		GetAllActiveTicketStatuses: _APIRoute{Path: mergeStrings(_APIRoutesPrefixes.Tickets.prefix, "GetAllActiveTicketStatuses/"), method: string(GetMethod), Status: true},
	},
	Auth: auth{
		LoginWithNoAuth:         _APIRoute{Path: mergeStrings(_APIRoutesPrefixes.Auth.prefix, "LoginWithNoAuth/"), method: string(GetMethod), Status: true},
		SignUp:                  _APIRoute{Path: mergeStrings(_APIRoutesPrefixes.Auth.prefix, "SignUp/"), method: string(PostMethod), Status: false},
		Login:                   _APIRoute{Path: mergeStrings(_APIRoutesPrefixes.Auth.prefix, "Login/"), method: string(PostMethod), Status: false},
		GetSingleUseToken:       _APIRoute{Path: mergeStrings(_APIRoutesPrefixes.Auth.prefix, "GetSingleUseToken/"), method: string(PostMethod), Status: true},
		LoginWithSingleUseToken: _APIRoute{Path: mergeStrings(_APIRoutesPrefixes.Auth.prefix, "LoginWithSingleUseToken/"), method: string(GetMethod), Status: true},
	},
	Captcha: captcha{
		GetCaptcha:    _APIRoute{Path: mergeStrings(_APIRoutesPrefixes.Captcha.prefix, "GetCaptcha/"), method: string(GetMethod), Status: true},
		VerifyCaptcha: _APIRoute{Path: mergeStrings(_APIRoutesPrefixes.Captcha.prefix, "VerifyCaptcha/"), method: string(PostMethod), Status: true},
	},
	Departments: departments{
		GetAllActiveDepartments: _APIRoute{Path: mergeStrings(_APIRoutesPrefixes.Department.prefix, "GetAllActiveDepartments/"), method: string(GetMethod), Status: true},
	},
	Users: users{
		GetUserByUsername: _APIRoute{Path: mergeStrings(_APIRoutesPrefixes.User.prefix, "GetUserByUsername/"), method: string(PostMethod), Status: true},
		GetUserByID:       _APIRoute{Path: mergeStrings(_APIRoutesPrefixes.User.prefix, "GetUserByID/"), method: string(PostMethod), Status: true},
		GetUsersByIDs:     _APIRoute{Path: mergeStrings(_APIRoutesPrefixes.User.prefix, "GetUsersByIDs/"), method: string(PostMethod), Status: true},
	},
	Files: files{
		UploadTicketFile:   _APIRoute{Path: mergeStrings(_APIRoutesPrefixes.Files.prefix, "UploadTicketFile/"), method: string(PostMethod), Status: true},
		DownloadTicketFile: _APIRoute{Path: mergeStrings(_APIRoutesPrefixes.Files.prefix, "DownloadTicketFile/:objectName"), method: string(GetMethod), Status: true},
	},
}

func mergeStrings(string ...string) string {
	var builder strings.Builder
	for _, s := range string {
		builder.WriteString(s)
	}
	return builder.String()
}

func IsRouteEnabled(path, method string) bool {
	allRoutes := []_APIRoute{
		APIRoutes.Tickets.CreateTicket,
		APIRoutes.Tickets.GetTicketByID,
		APIRoutes.Tickets.GetTicketByTrackCode,
		APIRoutes.Tickets.CreateChat,
		APIRoutes.Tickets.GetTicketsList,
		APIRoutes.Tickets.GetAllActiveTicketTypes,
		APIRoutes.Tickets.GetAllActiveTicketStatuses,
		APIRoutes.Auth.LoginWithNoAuth,
		APIRoutes.Auth.SignUp,
		APIRoutes.Auth.Login,
		APIRoutes.Auth.GetSingleUseToken,
		APIRoutes.Auth.LoginWithSingleUseToken,
		APIRoutes.Captcha.GetCaptcha,
		APIRoutes.Captcha.VerifyCaptcha,
		APIRoutes.Departments.GetAllActiveDepartments,
		APIRoutes.Users.GetUserByUsername,
		APIRoutes.Users.GetUserByID,
		APIRoutes.Users.GetUsersByIDs,
		APIRoutes.Files.DownloadTicketFile,
		APIRoutes.Files.UploadTicketFile,
	}
	for _, r := range allRoutes {
		if r.Path == path && r.method == method {
			return r.Status
		}
	}
	return true // default allow if not listed
}
