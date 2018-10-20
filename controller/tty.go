package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/websocket"
	projectModel "github.com/sysu-go-online/project-service/model"
	"github.com/sysu-go-online/public-service/tools"
	"github.com/sysu-go-online/public-service/types"
	"github.com/sysu-go-online/tty-service/model"
	userModel "github.com/sysu-go-online/user-service/model"
)

// WebSocketTermHandler is a middle way handler to connect web app with docker service
func WebSocketTermHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	// Set TextMessage as default
	msgType := websocket.TextMessage
	clientMsg := make(chan RequestCommand)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()

	// Open a goroutine to receive message from client connection
	go readFromClient(clientMsg, ws)

	// keep connection
	go func() {
		for {
			timer := time.NewTimer(time.Second * 2)
			<-timer.C
			err := ws.WriteControl(websocket.PingMessage, []byte("ping"), time.Time{})
			if err != nil {
				timer.Stop()
				return
			}
		}
	}()

	// Handle messages from the channel
	isFirst := true
	var sConn *websocket.Conn
	for msg := range clientMsg {
		conn := handlerClientTTYMsg(&isFirst, ws, sConn, msgType, &msg)
		sConn = conn
	}
	if sConn != nil {
		sConn.Close()
	}
}

// HandlerClientMsg handle message from client and send it to docker service
func handlerClientTTYMsg(isFirst *bool, ws *websocket.Conn, sConn *websocket.Conn, msgType int, connectContext *RequestCommand) (conn *websocket.Conn) {
	r := &TTYResponse{}
	var id string
	if *isFirst {
		// check token
		ok, username := tools.GetUserNameFromToken(connectContext.JWT, AuthRedisClient)
		connectContext.username = username
		if !ok {
			fmt.Fprintln(os.Stderr, "Can not get user token information")
			r.OK = false
			r.Msg = "Invalid token"
			ws.WriteJSON(r)
			ws.Close()
			conn = nil
			return
		}

		// Get project information
		session := MysqlEngine.NewSession()
		u := userModel.User{Username: username}
		ok, err := u.GetWithUsername(session)
		if !ok {
			fmt.Fprintln(os.Stderr, "Can not get user information")
			r.OK = false
			r.Msg = "Invalid user information"
			ws.WriteJSON(r)
			ws.Close()
			conn = nil
			return
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			r.OK = false
			r.Msg = err.Error()
			ws.WriteJSON(r)
			ws.Close()
			conn = nil
			return
		}
		p := projectModel.Project{Name: connectContext.Project, UserID: u.ID}
		has, err := p.GetWithUserIDAndName(session)
		if !has {
			fmt.Fprintln(os.Stderr, "Can not get project information")
			r.OK = false
			r.Msg = "Can not get project information"
			ws.WriteJSON(r)
			ws.Close()
			conn = nil
			return
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			r.OK = false
			r.Msg = err.Error()
			ws.WriteJSON(r)
			ws.Close()
			conn = nil
			return
		}

		// send request for start a container
		userHome := filepath.Join("/home", username)
		body := NewContainer{
			Image:     "go-online-golang_image",
			Command:   "bash",
			PWD:       "/root/",
			ENV:       []string{},
			Mnt:       []string{userHome},
			TargetDir: []string{"/root"},
			Network:   []string{},
		}
		b, err := json.Marshal(body)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			r.OK = false
			r.Msg = err.Error()
			ws.WriteJSON(r)
			ws.Close()
			conn = nil
			return
		}
		id, err = startContainer(b)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			r.OK = false
			r.Msg = err.Error()
			ws.WriteJSON(r)
			ws.Close()
			conn = nil
			return
		}

		// TODO: write container information to the redis
		// connect to the container
		tmp, err := initDockerConnection("tty", id)
		sConn = tmp
		if err != nil {
			fmt.Println("Can not connect to the docker service")
			r.OK = false
			r.Msg = "Server error"
			ws.WriteJSON(r)
			ws.Close()
			return
		}
		// Listen message from docker service and send to client connection
		go sendTTYMsgToClient(ws, sConn)
	}

	if sConn == nil {
		fmt.Fprintf(os.Stderr, "Invalid command.")
		ws.WriteMessage(msgType, []byte("Invalid Command"))
		ws.Close()
		conn = nil
		return
	}

	// Send message to docker service
	handleTTYMessage(msgType, sConn, id, connectContext.Message)
	*isFirst = false
	conn = sConn
	return
}

// SendMsgToClient send message to client
func sendTTYMsgToClient(cConn *websocket.Conn, sConn *websocket.Conn) {
	for {
		r := &TTYResponse{}
		err := sConn.ReadJSON(r)
		if err != nil {
			fmt.Println(err)
			// Server closed connection
			cConn.Close()
			return
		}
		cConn.WriteJSON(r)
	}
}

// HandleMessage decide different operation according to the given json message
func handleTTYMessage(mType int, conn *websocket.Conn, id string, msg string) error {
	Msg := ByteStreamToDocker{
		ID:  id,
		Msg: msg,
	}

	return conn.WriteJSON(Msg)
}

// RegisterPortAndDomainInfo register port
func RegisterPortAndDomainInfo(mapping *types.PortMapping, containerName string) error {
	CONSULADDRESS := os.Getenv("CONSUL_ADDRESS")
	if len(CONSULADDRESS) == 0 {
		CONSULADDRESS = "localhost"
	}
	CONSULPORT := os.Getenv("CONSUL_PORT")
	if len(CONSULPORT) == 0 {
		CONSULPORT = "8500"
	}
	if CONSULPORT[0] != ':' {
		CONSULPORT = ":" + CONSULPORT
	}

	url := "http://" + CONSULADDRESS + CONSULPORT + "/v1/kv/upstreams/"
	err := model.AddDomainName(mapping.DomainName, DomainNameRedisClient)
	if err != nil {
		return err
	}
	DOMAINNAME := os.Getenv("DOMAIN_NAME")
	if len(DOMAINNAME) != 0 {
		DOMAINNAME = "." + DOMAINNAME
	}
	req := model.RegisterConsulParam{
		Key:   mapping.DomainName,
		Value: fmt.Sprintf("%s:%d", containerName[0:12], mapping.Port),
	}
	return req.RegisterToConsul(url)
}
