package main

import (
	"log"
	"time"
	"github.com/gorilla/websocket"
)

var channelWriteSocket chan string

func Ws(url string){
	channelReadSocket := make(chan string)
	channelWriteSocket = make(chan string)
	channelHealthCheck := make(chan bool)	

	c, err := RegisterWebsocketServer(url)	

	for err {
		log.Println("error on connect trying Connect...")
		c, err = RegisterWebsocketServer(url)
		time.Sleep(5 * time.Second)
	}

	ReadSocketMessage(c, channelReadSocket)
	HealthCheck(c, channelHealthCheck)
	SendSocketMessage(c, channelWriteSocket)
	OnOpen()
	
	for {
		select {
		case message := <-channelReadSocket:
			log.Println("----------------------- Processando Mensagens Websocket -----------------------\n")
			OnMessage(message)

		case t2 := <-channelHealthCheck:
			if t2 == true {				
				OnClose()

				c, err = RegisterWebsocketServer(url)
				for err {
					c, err = RegisterWebsocketServer(url)
					
					time.Sleep(5 * time.Second)
					if err == false {
						log.Println("t2 Connected...") //esse cara aqui parece que não reconecta não						
					}
				}
				// Must register all methods after reconnect.
				channelWriteSocket = make(chan string)

				ReadSocketMessage(c, channelReadSocket)
				HealthCheck(c, channelHealthCheck)
				SendSocketMessage(c, channelWriteSocket)
				OnOpen()	
			} else {
				log.Println("connection is up.")
			}		
	
		}
	}
}

func RegisterWebsocketServer(url string) (*websocket.Conn, bool){
	c, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		log.Println("Error to connect", err)
		time.Sleep(1 * time.Second)
		return nil, true
	}

	log.Println("Connected")
	return c, false
}

func ReadSocketMessage(c *websocket.Conn, ch chan string) {
	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("error conection --->:", err)
				return // kill go routine
			}
			ch <- string(message)
		}
	}()
}

func SendSocketMessage(c *websocket.Conn, sendSocketMsg chan string) {
	go func() {
		log.Println("Inicializando Fila de mensagens")
		for {
			c.WriteMessage(websocket.TextMessage, []byte(<-sendSocketMsg))
		}
	}()
}

func HealthCheck(c *websocket.Conn, ch chan bool) {
	go func() {
		for {
			time.Sleep(15 * time.Second)

			err := c.WriteMessage(websocket.TextMessage, []byte("ping"))
			if err != nil {
				log.Println("error trying reconnection...", err)
				ch <- true
				return // kill go routine
			} else {
				log.Println("connection is up.")
				ch <- false
			}
		}
	}()
}