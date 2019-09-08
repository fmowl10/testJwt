package backend

import "log"

//https://dev.to/danielkun/go-asynchronous-and-safe-real-time-broadcasting-using-channels-and-websockets-4g5d

type Clients struct {
	clients   []chan<- string
	singChan  chan chan string
	unignChan chan chan string
}

func (c Clients) Broadcast(s string) {
	log.Println("client log: broadcasting: " + s)
	for _, v := range c.clients {
		v <- s
	}
}

func (c *Clients) Init() {
	log.Println("client log: Init")
	c.clients = make([]chan<- string, 0, 10)
	c.singChan = make(chan chan string)
	c.unignChan = make(chan chan string)
}

func (c Clients) Sign(clientChan chan string) {
	log.Println("client log: Signed")
	c.singChan <- clientChan
}

func (c Clients) Unsign(clientChan chan string) {
	log.Println("client log: Unsigned")
	c.unignChan <- clientChan
}

func (c *Clients) Hub() {
	for {
		log.Println("number of Chans: ", len(c.clients))
		select {
		case client := <-c.singChan:
			c.clients = append(c.clients, client)
		case clientChan := <-c.unignChan:
			for i, v := range c.clients {
				if v == clientChan {
					c.clients = append(c.clients[:i], c.clients[i+1:]...)
					close(clientChan)
				}
			}
		}
	}
}
