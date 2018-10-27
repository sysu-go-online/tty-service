package controller

import (
	"reflect"
	"testing"

	"github.com/gorilla/websocket"
)

func Test_initDockerConnection(t *testing.T) {
	type args struct {
		service string
		id      string
	}
	tests := []struct {
		name    string
		args    args
		want    *websocket.Conn
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := initDockerConnection(tt.args.service, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("initDockerConnection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("initDockerConnection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dialDockerService(t *testing.T) {
	type args struct {
		service string
		id      string
	}
	tests := []struct {
		name    string
		args    args
		want    *websocket.Conn
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dialDockerService(tt.args.service, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("dialDockerService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dialDockerService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_startContainer(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := startContainer(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("startContainer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("startContainer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readFromClient(t *testing.T) {
	type args struct {
		clientChan chan<- RequestCommand
		ws         *websocket.Conn
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			readFromClient(tt.args.clientChan, tt.args.ws)
		})
	}
}

func Test_getPwd(t *testing.T) {
	type args struct {
		projectName string
		username    string
		projectType int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPwd(tt.args.projectName, tt.args.username, tt.args.projectType); got != tt.want {
				t.Errorf("getPwd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getEnv(t *testing.T) {
	type args struct {
		projectName string
		username    string
		language    int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getEnv(tt.args.projectName, tt.args.username, tt.args.language); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
