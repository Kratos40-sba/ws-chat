package handlers

import (
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sort"
)

var wsChan = make(chan WsPayload)
var clients = make(map[WsConnection]string)
var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(),
)
var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Home(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "home.jet", nil)
	if err != nil {
		log.Println(err)

	}
}

type WsJsonResponse struct {
	Action         string   `json:"action"`
	Message        string   `json:"message"`
	MessageType    string   `json:"message_type"`
	ConnectedUsers []string `json:"connected_users"`
}
type WsPayload struct {
	Action   string       `json:"action"`
	UserName string       `json:"user_name"`
	Message  string       `json:"message"`
	Conn     WsConnection `json:"-"`
}
type WsConnection struct {
	*websocket.Conn
}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("Client connected to ws endpoint")
	var response WsJsonResponse
	response.Message = `<small>Connected to server</small>`
	conn := WsConnection{Conn: ws}
	clients[conn] = ""
	err = ws.WriteJSON(response)
	if err != nil {
		log.Println(err)
	}
	go ListenForWs(&conn)
}
func ListenForWs(conn *WsConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error ", fmt.Sprintf("%v", r))
		}
	}()
	var payload WsPayload
	for {
		err := conn.ReadJSON(&payload)
		if err != nil {

		} else {
			payload.Conn = *conn
			wsChan <- payload
		}
	}
}
func ListenToWsChannel() {
	var resp WsJsonResponse
	for {
		event := <-wsChan
		switch event.Action {
		case "username":
			clients[event.Conn] = event.UserName
			users := getUserList()
			resp.Action = "list_users_action"
			resp.ConnectedUsers = users
			broadCastToAll(resp)
		case "left":
			resp.Action = "list_users_action"

			delete(clients, event.Conn)

			users := getUserList()
			resp.ConnectedUsers = users
			broadCastToAll(resp)
		case "broadcast":
			resp.Action = "broadcast"
			resp.Message = fmt.Sprintf("<strong>%s</strong>: %s ", event.UserName, event.Message)
			broadCastToAll(resp)
		}
		//resp.Action = "Got here"
		//resp.Message = fmt.Sprintf("Some message , and action was %s",event.Action)
		//broadCastToAll(resp)
	}

}
func getUserList() []string {
	var userList []string
	for _, u := range clients {
		if u != "" {
			userList = append(userList, u)
		}

	}
	sort.Strings(userList)
	return userList
}
func broadCastToAll(response WsJsonResponse) {
	for c := range clients {
		err := c.WriteJSON(response)
		if err != nil {
			log.Println("websocket error")
			_ = c.Close()
			delete(clients, c)
		}
	}
}
func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	v, err := views.GetTemplate(tmpl)
	if err != nil {
		log.Println(err)
		return err
	}
	err = v.Execute(w, data, nil)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
