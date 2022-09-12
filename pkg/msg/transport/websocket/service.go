package websocket

import (
	"fmt"
	"time"
)

func Job(p *PoolStructure,  stopchan <-chan bool) {
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
