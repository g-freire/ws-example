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

func getLastIBOPJob(host, db string, stopchan <-chan bool) {
	//defer func(){fmt.Print("FINAL WORK TO DO")}()
	for {
		select {
		default:
			Pool.Broadcast <- "result"
			time.Sleep(1000 * time.Millisecond)

		case <-stopchan:
			fmt.Print("\n CLOSING INFINITE QUERY LOOP \n")
			return
		}
	}
}

func (c *Client) Read(ctx *gin.Context) {
	defer func() {
		Pool.Unregister <- c
		c.Socket.Close()
	}()
	for {
		_, _, err := c.Socket.ReadMessage()
		if err != nil {
			Pool.Unregister <- c
			c.Socket.Close()
			break
		}
	}
}

func (c *Client) Write(client *Client) {
	defer func() {
		c.Socket.Close()
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
