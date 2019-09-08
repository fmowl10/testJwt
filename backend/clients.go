package backend

import "log"

// Clients is collection of clients
// Clients do sign / unsign and boradcast message
type Clients struct {
	clients   []chan<- string
	singChan  chan chan string
	unignChan chan chan string
}

// Broadcast send "s" all clients
func (c Clients) Broadcast(s string) {
	log.Println(`client log: broadcasting: "` + s + `"`)
	for _, v := range c.clients {
		v <- s
	}
}

// Init intialize CLients
func (c *Clients) Init() {
	log.Println("client log: Init")
	c.clients = make([]chan<- string, 0, 10)
	c.singChan = make(chan chan string)
	c.unignChan = make(chan chan string)
}

// Sign add channel to c.clients
func (c Clients) Sign(clientChan chan string) {
	log.Println("client log: Signed")
	c.singChan <- clientChan
}

// Unsign remove channel from c.clients and close channel
func (c Clients) Unsign(clientChan chan string) {
	log.Println("client log: Unsigned")
	c.unignChan <- clientChan
}

// Hub do serve
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
