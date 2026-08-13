package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	addrFlag := flag.String("addr", "localhost:9091", "адрес UDP-сервера")
	flag.Parse()

	addr, err := net.ResolveUDPAddr("udp", *addrFlag)
	if err != nil {
		log.Fatalf("разобрать UDP-адрес %s: %v", *addrFlag, err)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("подключиться к %s: %v", *addrFlag, err)
	}
	defer conn.Close()

	messages := []string{
		"Привет, UDP!",
		"Это датаграмма",
		"Никаких гарантий доставки",
	}
	buffer := make([]byte, 64*1024)
	for _, message := range messages {
		if err := conn.SetDeadline(time.Now().Add(3 * time.Second)); err != nil {
			log.Fatalf("установить таймаут: %v", err)
		}
		fmt.Println("Отправляю:", message)
		if _, err := conn.Write([]byte(message)); err != nil {
			log.Fatalf("отправить датаграмму: %v", err)
		}

		n, err := conn.Read(buffer)
		if err != nil {
			log.Fatalf("прочитать ответ: %v", err)
		}
		fmt.Printf("Ответ: %s\n\n", buffer[:n])
	}
}
