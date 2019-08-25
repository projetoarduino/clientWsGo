package main

import (
	"log"
	"time"
	"os"
	"os/signal"
	"github.com/projetoarduino/clientws"
)

var ws = clientws.New("ws://192.168.1.16:8080")

func main(){
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)	

	go ws.Connect()
	
	ws.OnOpen = func(){
		log.Println("on Open")
		clientws.Send("Menssage inicial")
	}

	ws.OnMessage = func(message string){
		log.Println("on message")
		log.Println(message)	
	}

	ws.OnClose = func(){
		log.Println("connection is down.")				
		log.Println("Trying Reconnect.")
	}

	time.Sleep(5 * time.Second)
	clientws.Send("teste 1234")

	for {
		select {
			case <-interrupt:
				log.Println("interrupt")
				// ToDo graceful shutdown
			return
		}
	}	
}






