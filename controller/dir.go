package controller

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/websocket"
	"github.com/rjeczalik/notify"
	projectModel "github.com/sysu-go-online/project-service/model"
	"github.com/sysu-go-online/public-service/tools"
	userModel "github.com/sysu-go-online/user-service/model"
)

// DirRequest contains message send from user
type DirRequest struct {
	JWT     string `json:"jwt"`
	Project string `json:"project"`
}

// DirResponse stores response to the user
type DirResponse struct {
	OK   bool   `json:"ok"`
	Path string `json:"path"`
	Type int    `json:"type"`
}

// MonitorDirHandler monitor user dir status
func MonitorDirHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		fmt.Println(err)
		return
	}

	// wait for the first message
	msg := DirRequest{}
	err = ws.ReadJSON(&msg)
	if err != nil {
		fmt.Println(err)
		ws.Close()
		return
	}

	// parse jwt for auth and get relevant message
	// check token
	res := DirResponse{}
	ok, username := tools.GetUserNameFromToken(msg.JWT, AuthRedisClient)
	if !ok {
		fmt.Fprintln(os.Stderr, "Can not get user token information")
		res.OK = false
		ws.WriteJSON(r)
		ws.Close()
		return
	}

	// Get project information
	session := MysqlEngine.NewSession()
	u := userModel.User{Username: username}
	ok, err = u.GetWithUsername(session)
	if !ok {
		fmt.Fprintln(os.Stderr, "Can not get user information")
		res.OK = false
		ws.WriteJSON(r)
		ws.Close()
		return
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		res.OK = false
		ws.WriteJSON(r)
		ws.Close()
		return
	}
	p := projectModel.Project{Name: msg.Project, UserID: u.ID}
	has, err := p.GetWithUserIDAndName(session)
	if !has {
		fmt.Fprintln(os.Stderr, "Can not get project information")
		res.OK = false
		ws.WriteJSON(r)
		ws.Close()
		return
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		res.OK = false
		ws.WriteJSON(r)
		ws.Close()
		return
	}

	// listen to dir events
	path := filepath.Join("/home", username, "projects", p.Path, p.Name)
	c := make(chan notify.EventInfo, 1)
	if err := notify.Watch(path, c, notify.All); err != nil {
		log.Fatal(err)
	}
	defer notify.Stop(c)

	res.OK = true
	res.Type = -1
	err = ws.WriteJSON(res)
	if err != nil {
		log.Println(err)
		return
	}

	// Block until an event is received.
	for n := range c {
		res.Path = n.Path()
		res.OK = true
		switch n.Event() {
		case notify.Create:
			res.Type = 0
		case notify.Remove:
			res.Type = 1
		case notify.Write:
			res.Type = 2
		case notify.InMovedFrom:
			res.Type = 3
		case notify.InMovedTo:
			res.Type = 4
		default:
			res.Type = 5
		}
		err = ws.WriteJSON(res)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
