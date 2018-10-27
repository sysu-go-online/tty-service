package controller

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/sysu-go-online/public-service/types"
)

func TestWebSocketTermHandler(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WebSocketTermHandler(tt.args.w, tt.args.r)
		})
	}
}

func Test_handlerClientTTYMsg(t *testing.T) {
	type args struct {
		isFirst        *bool
		ws             *websocket.Conn
		sConn          *websocket.Conn
		msgType        int
		connectContext *RequestCommand
	}
	tests := []struct {
		name     string
		args     args
		wantConn *websocket.Conn
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotConn := handlerClientTTYMsg(tt.args.isFirst, tt.args.ws, tt.args.sConn, tt.args.msgType, tt.args.connectContext); !reflect.DeepEqual(gotConn, tt.wantConn) {
				t.Errorf("handlerClientTTYMsg() = %v, want %v", gotConn, tt.wantConn)
			}
		})
	}
}

func Test_sendTTYMsgToClient(t *testing.T) {
	type args struct {
		cConn *websocket.Conn
		sConn *websocket.Conn
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sendTTYMsgToClient(tt.args.cConn, tt.args.sConn)
		})
	}
}

func Test_handleTTYMessage(t *testing.T) {
	type args struct {
		mType int
		conn  *websocket.Conn
		id    string
		msg   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := handleTTYMessage(tt.args.mType, tt.args.conn, tt.args.id, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("handleTTYMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRegisterPortAndDomainInfo(t *testing.T) {
	type args struct {
		mapping       *types.PortMapping
		containerName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RegisterPortAndDomainInfo(tt.args.mapping, tt.args.containerName); (err != nil) != tt.wantErr {
				t.Errorf("RegisterPortAndDomainInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
