package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

// InitDockerConnection inits the connection to the docker service with the first message received from client
func initDockerConnection(service string, id string) (*websocket.Conn, error) {
	conn, err := dialDockerService(service, id)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// DialDockerService create connection between web server and docker server
// Accept service type:
// tty debug
func dialDockerService(service string, id string) (*websocket.Conn, error) {
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
	url := url.URL{Scheme: "ws", Host: dockerAddr, Path: "/tty"}
	conn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func startContainer(b []byte) (string, error) {
	// get addr
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
	url := url.URL{Scheme: "http", Host: dockerAddr, Path: "/create"}
	resp, err := http.Post(url.String(), "application/x-www-form-urlencoded", strings.NewReader(string(b)))
	if err != nil {
		log.Println(err)
		return "", err
	}
	res := NewContainerRet{}
	retBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(retBody, &res)
	if err != nil {
		return "", err
	}
	if !res.OK {
		return "", errors.New(res.Msg)
	}
	return res.ID, nil
}

func resizeContainer(r *ResizeContainer) error {
	// get addr
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
	url := url.URL{Scheme: "http", Host: dockerAddr, Path: "/resize"}
	b, err := json.Marshal(r)
	if err != nil {
		log.Println(err)
		return err
	}
	resp, err := http.Post(url.String(), "application/x-www-form-urlencoded", strings.NewReader(string(b)))
	if err != nil {
		log.Println(err)
		return err
	}
	ret := ResizeContainerRet{}
	retBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	err = json.Unmarshal(retBody, &ret)
	if err != nil {
		log.Println(err)
		return err
	}
	if ret.OK {
		return nil
	}
	log.Println(ret.Msg)
	return errors.New(ret.Msg)
}

// ReadFromClient receive message from client connection
func readFromClient(clientChan chan<- RequestCommand, ws *websocket.Conn) {
	for {
		_, b, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				log.Println("Remote user closed the connection")
				ws.Close()
				close(clientChan)
				break
			}
			close(clientChan)
			log.Println(err)
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
