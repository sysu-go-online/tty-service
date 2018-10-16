package router

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"github.com/sysu-go-online/tty-service/controller"
	"github.com/urfave/negroni"
)

var upgrader = websocket.Upgrader{}

// GetServer return web server
func GetServer() *negroni.Negroni {
	r := mux.NewRouter()

	r.HandleFunc("/tty", controller.WebSocketTermHandler)
	r.HandleFunc("/debug", controller.DebugHandler)

	// Use classic server and return it
	handler := cors.Default().Handler(r)
	s := negroni.Classic()
	s.UseHandler(handler)
	return s
}
