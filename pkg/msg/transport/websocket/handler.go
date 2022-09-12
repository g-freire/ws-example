package websocket

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func GetLastEvent(c *gin.Context) {
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if error != nil {
		http.NotFound(c.Writer, c.Request)
		fmt.Println(error)
		return
	}
	uuid := uuid.Must(uuid.NewV4(), error).String()
	client := &Client{Id: uuid, Ip: "", Socket: conn, ClientChan: make(chan interface{})}

	Pool.Register <- client
	fmt.Println("A new ws client has connected.", client.Id, client.Ip)

	go client.Read(c)
	go client.Write(client)
}
