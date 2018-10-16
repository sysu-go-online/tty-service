package controller

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

// InitDockerConnection inits the connection to the docker service with the first message received from client
func initDockerConnection(service string) (*websocket.Conn, error) {
	// Just handle command start with `go`
	conn, err := dialDockerService(service)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// DialDockerService create connection between web server and docker server
// Accept service type:
// tty debug
func dialDockerService(service string) (*websocket.Conn, error) {
	// Set up websocket connection
	dockerAddr := os.Getenv("DOCKER_ADDRESS")
	dockerPort := os.Getenv("DOCKER_PORT")
	if len(dockerAddr) == 0 {
		dockerAddr = "localhost"
	}
	if len(dockerPort) == 0 {
		dockerPort = "8888"
	}
	dockerPort = ":" + dockerPort
	dockerAddr = dockerAddr + dockerPort
	url := url.URL{Scheme: "ws", Host: dockerAddr, Path: "/" + service}
	conn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// ReadFromClient receive message from client connection
func readFromClient(clientChan chan<- RequestCommand, ws *websocket.Conn) {
	for {
		_, b, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				fmt.Fprintln(os.Stderr, "Remote user closed the connection")
				ws.Close()
				close(clientChan)
				break
			}
			close(clientChan)
			fmt.Fprintln(os.Stderr, "Can not read message.")
			return
		}
		// read json message from rws
		msg := RequestCommand{}
		if err := json.Unmarshal(b, &msg); err != nil {
			fmt.Fprintln(os.Stderr, "Can not parse data")
			ws.Close()
			close(clientChan)
			break
		}

		clientChan <- msg
	}
}

// getPwd return current path of given username
func getPwd(projectName string, username string, projectType int) string {
	// TODO: return according to the context
	return "/"
}

func getEnv(projectName string, username string, language int) []string {
	env := []string{}
	switch language {
	case 0:
		// golang
		env = append(env, "GOPATH=/root/go:/home/go")
	}
	return env
}
