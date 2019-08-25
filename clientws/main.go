package main

import (
	"log"
	"time"
	"os"
	"os/signal"
)

func main(){
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)	

	go Ws("ws://192.168.1.16:8080")	

	time.Sleep(5 * time.Second)
	Send("teste 1234")

	for {
		select {
			case <-interrupt:
				log.Println("interrupt")
				// ToDo graceful shutdown
			return
		}
	}	
}

func OnOpen(){
	log.Println("on Open")
	Send("Menssage inicial")
}

func OnMessage(message string){
	log.Println("on message")
	log.Println(message)	
}

func OnClose(){
	log.Println("connection is down.")				
	log.Println("Trying Reconnect.")
}

func Send(message string){
	channelWriteSocket <- message
}

