package websocket

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"ws-example/internal/config"
)

type Handler struct {
	pool PoolStructure
	conf *config.Config
}

func NewHandler(pool PoolStructure, conf *config.Config) *Handler {
	return &Handler{pool: pool, conf: conf}
}


func (h *Handler)GetLastEvent(c *gin.Context) {
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if error != nil {
		http.NotFound(c.Writer, c.Request)
		fmt.Println(error)
		return
	}
	uuid := uuid.Must(uuid.NewV4(), error).String()
	client := &Client{Id: uuid, Ip: "", Socket: conn, ClientChan: make(chan interface{})}

	h.pool.Register <- client
	fmt.Println("A new ws client has connected.", client.Id, client.Ip)

	go client.Read(h, c)
	go client.Write(h,client)
}
