package websocket

import (
	"fmt"
	"sync"
)

type PoolStructure struct {
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan interface{}
	Clients    map[*Client]bool
	WG sync.WaitGroup
}

func NewPool(wg sync.WaitGroup) *PoolStructure {
	return &PoolStructure{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan interface{}),
		Clients:    make(map[*Client]bool),
		WG: wg,
	}
}


func (pool *PoolStructure) Start(host, db string) {
	var jobChan = make(chan bool)
	var pollSize int
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			pollSize = len(pool.Clients)
			fmt.Print("## POOL SIZE NOW: ", pollSize, "\n")
			if pollSize == 1 {
				fmt.Print("** CREATING THE SINGLETON JOB \n")
				jobChan = make(chan bool)
				pool.WG.Add(1)
				defer pool.WG.Done()
				go Job(pool, jobChan)
			}
		case client := <-pool.Unregister:
			if _, ok := pool.Clients[client]; ok {
				close(client.ClientChan)
				delete(pool.Clients, client)
				fmt.Println("A socket has disconnected.")
				pollSize = len(pool.Clients)
				fmt.Print("## POOL SIZE NOW: ", pollSize, "\n")
				if pollSize == 0 {
					fmt.Print("** CLOSING THE SINGLETON JOB \n")
					close(jobChan)
					//jobChan <- true
				}
			}
		case message := <-pool.Broadcast:
			for client := range pool.Clients {
				select {
				case client.ClientChan <- message:
				default:
					close(client.ClientChan)
					delete(pool.Clients, client)
				}
			}
		}
	}
}
