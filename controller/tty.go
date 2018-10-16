package controller

import (
	"fmt"
	"net/http"
	"os"
	"strings"
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
	var mapping *types.PortMapping
	mapping = nil
	if *isFirst {
		// check token
		ok, username := tools.GetUserNameFromToken(connectContext.JWT, AuthRedisClient)
		connectContext.username = username
		if !ok {
			fmt.Fprintln(os.Stderr, "Can not get user token information")
			r.Type = "error"
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
			r.Type = "error"
			r.Msg = "Invalid user information"
			ws.WriteJSON(r)
			ws.Close()
			conn = nil
			return
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			r.Type = "error"
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
			r.Type = "error"
			r.Msg = "Can not get project information"
			ws.WriteJSON(r)
			ws.Close()
			conn = nil
			return
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			r.Type = "error"
			r.Msg = err.Error()
			ws.WriteJSON(r)
			ws.Close()
			conn = nil
			return
		}
		connectContext.projectType = p.Language

		// Check if the command is system command
		command := strings.Split(connectContext.Command, " ")
		if len(command) <= 0 {
			fmt.Fprintln(os.Stderr, "Can not parse command")
			r.Type = "error"
			r.Msg = "invalid command"
			ws.WriteJSON(r)
			ws.Close()
			conn = nil
			return
		}
		if command[0] == "go-online" {
			mapping = &types.PortMapping{}
			// parse command
			mapping, err = tools.ParseSystemCommand(command, DomainNameRedisClient)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				r.Type = "error"
				r.Msg = err.Error()
				ws.WriteJSON(r)
				ws.Close()
				conn = nil
				return
			}
			connectContext.Command = mapping.Command
		}

		tmp, err := initDockerConnection("tty")
		sConn = tmp
		if err != nil {
			fmt.Println("Can not connect to the docker service")
			r.Type = "error"
			r.Msg = "Server error"
			ws.WriteJSON(r)
			ws.Close()
		}
		// Listen message from docker service and send to client connection
		go sendTTYMsgToClient(ws, sConn, mapping)
	}

	if sConn == nil {
		fmt.Fprintf(os.Stderr, "Invalid command.")
		ws.WriteMessage(msgType, []byte("Invalid Command"))
		ws.Close()
		conn = nil
		return
	}

	// Send message to docker service
	handleTTYMessage(msgType, sConn, *isFirst, connectContext, mapping)
	*isFirst = false
	conn = sConn
	return
}

// SendMsgToClient send message to client
func sendTTYMsgToClient(cConn *websocket.Conn, sConn *websocket.Conn, mapping *types.PortMapping) {
	type dockerMsg struct {
		Msg  string `json:"msg"`
		ID   string `json:"id"`
		Type string `json:"type"`
	}
	for {
		msg := &dockerMsg{}
		err := sConn.ReadJSON(msg)
		r := &TTYResponse{}
		if err != nil {
			fmt.Println(err)
			// Server closed connection
			cConn.Close()
			return
		}
		// register for the first time
		if mapping != nil {
			err := RegisterPortAndDomainInfo(mapping, msg.ID)
			if err != nil {
				fmt.Println(err)
				r.Type = "error"
				r.Msg = err.Error()
				cConn.WriteJSON(r)
			}
			DOMAINNAME := os.Getenv("DOMAIN_NAME")
			if len(DOMAINNAME) != 0 {
				DOMAINNAME = "." + DOMAINNAME
			} else {
				DOMAINNAME = ".localhost"
			}
			r.Type = "dname"
			r.Msg = mapping.DomainName + DOMAINNAME
			cConn.WriteJSON(r)
			mapping = nil
		}
		r.Type = "tty"
		r.Msg = msg.Msg
		cConn.WriteJSON(r)
	}
}

// HandleMessage decide different operation according to the given json message
func handleTTYMessage(mType int, conn *websocket.Conn, isFirst bool, connectContext *RequestCommand, mapping *types.PortMapping) error {
	var workSpace *Command
	var err error
	if isFirst {
		projectName := connectContext.Project
		username := connectContext.username
		pwd := getPwd(projectName, username, connectContext.projectType)
		env := getEnv(projectName, username, connectContext.projectType)
		workSpace = &Command{
			Command:     connectContext.Command,
			PWD:         pwd,
			ENV:         env,
			UserName:    username,
			ProjectName: projectName,
			Type:        "tty",
		}
		if mapping != nil {
			workSpace.Ports = []int{mapping.Port}
		}
		fmt.Println(workSpace.ENV)
	}

	// Send message
	if isFirst {
		err = conn.WriteJSON(*workSpace)
	} else {
		err = conn.WriteJSON(connectContext)
	}
	if err != nil {
		return err
	}
	return nil
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
