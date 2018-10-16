package controller

import (
	"net/http"
)

// Command is the JSON format between web server and docker server
type Command struct {
	Command     string   `json:"command"`
	PWD         string   `json:"pwd"`
	ENV         []string `json:"env"`
	UserName    string   `json:"user"`
	ProjectName string   `json:"project"`
	Entrypoint  []string `json:"entrypoint"`
	Type        string   `json:"type"`
	Ports       []int    `json:"port"`
}

// ErrorHandler is error handler for http
type ErrorHandler func(w http.ResponseWriter, r *http.Request) error

func (h ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// RequestCommand stores command and jwt in every ws message
type RequestCommand struct {
	Command     string
	JWT         string
	Project     string
	username    string
	projectType int
}

// TTYResponse stores data to be sent to the client
type TTYResponse struct {
	Type string
	Msg  string `json:"msg"`
}
