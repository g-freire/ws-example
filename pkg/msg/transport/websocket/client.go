package websocket

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"time"
)

type Client struct {
	Id         string
	Ip         string
	Socket     *websocket.Conn
	ClientChan chan interface{}
}

func getLastJob(p *PoolStructure,  stopchan <-chan bool) {
	//defer func(){fmt.Print("FINAL WORK TO DO")}()
	for {
		select {
		default:
			p.Broadcast <- "TEST"
			time.Sleep(time.Second)

		case <-stopchan:
			fmt.Print("\n CLOSING INFINITE QUERY LOOP \n")
			return
		}
	}
}

func (c *Client) Read(h *Handler, ctx *gin.Context) {
	h.pool.WG.Add(1)
	defer func() {
		h.pool.Unregister <- c
		c.Socket.Close()
		h.pool.WG.Done()
	}()
	for {
		_, _, err := c.Socket.ReadMessage()
		if err != nil {
			h.pool.Unregister <- c
			c.Socket.Close()
			break
		}
	}
}

func (c *Client) Write(h *Handler, client *Client) {
	h.pool.WG.Add(1)
	defer func() {
		c.Socket.Close()
		h.pool.WG.Done()
	}()

	for {
		select {
		case message, ok := <-c.ClientChan:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Socket.WriteJSON(message)
			//fmt.Println("SENDING ", msg.Type," IP: ", client.Ip," ID ", client.Id)
		}
	}
}
