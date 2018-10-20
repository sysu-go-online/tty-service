package controller

// browser to service end

// RequestCommand stores command and jwt in every ws message
type RequestCommand struct {
	Message  string `json:"msg"`
	JWT      string `json:"jwt"`
	Project  string `json:"project"`
	username string
}

// TTYResponse stores data to be sent to the client
type TTYResponse struct {
	OK  bool   `json:"ok"`
	Msg string `json:"msg"`
}

// service end to docker end

// NewContainer is the JSON format between web server and docker server
type NewContainer struct {
	Image     string   `json:"image"`
	Command   string   `json:"command"`
	PWD       string   `json:"pwd"`
	ENV       []string `json:"env"`
	Mnt       []string `json:"mnt"`
	TargetDir []string `json:"target"`
	Network   []string `json:"network"`
}

// NewContainerRet is the create result returned by docker end
type NewContainerRet struct {
	ID  string `json:"id"`
	OK  bool   `json:"ok"`
	Msg string `json:"msg"`
}

// ByteStreamToDocker contains byte stream from user to container
type ByteStreamToDocker struct {
	ID  string `json:"id"`
	Msg string `json:"msg"`
}

// ByteStreamToUser stores byte stream from container to user
type ByteStreamToUser struct {
	OK  bool   `ok`
	Msg string `json:"msg"`
}
