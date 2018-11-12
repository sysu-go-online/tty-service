package controller

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	State int    `json:"state"`
	Type string `json:"type"`
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

	// listening for close event
	stop := make(chan bool, 0)
	// keep connection
	go func() {
		for {
			timer := time.NewTimer(time.Second * 2)
			<-timer.C
			err := ws.WriteControl(websocket.PingMessage, []byte("ping"), time.Time{})
			if err != nil {
				stop <- true
				timer.Stop()
				return
			}
		}
	}()

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
	if err := notify.Watch(path+"/...", c, notify.All); err != nil {
		log.Fatal(err)
		return
	}
	defer notify.Stop(c)

	res.OK = true
	res.State = -1
	err = ws.WriteJSON(res)
	if err != nil {
		log.Println(err)
		return
	}

	// Block until an event is received.
	go func() {
		for n := range c {
			res.Path = strings.Split(n.Path(), filepath.Join("/home", username, "projects/", p.Path, p.Name))[1][1:]
			
			fileInfo,err := os.Stat(res.Path);
			if err != nil {
				log.Println(err)
				return
			}
			if fileInfo.IsDir() {
				res.Type = "dir"
			}else {
				res.Type = "file"
			}
			res.OK = true
			switch n.Event() {
			case notify.Create:
				res.State = 0
			case notify.Remove:
				res.State = 1
			case notify.Write:
				res.State = 2
			case notify.InMovedFrom:
				res.State = 3
			case notify.InMovedTo:
				res.State = 4
			default:
				res.State = 5
			}
			err = ws.WriteJSON(res)
			if err != nil {
				stop <- true
				log.Println(err)
				return
			}
		}
	}()
	<-stop
}
