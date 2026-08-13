package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

func main() {
	addrFlag := flag.String("addr", ":9091", "UDP-адрес для прослушивания")
	flag.Parse()

	addr, err := net.ResolveUDPAddr("udp", *addrFlag)
	if err != nil {
		log.Fatalf("разобрать UDP-адрес %s: %v", *addrFlag, err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("запустить UDP-сервер на %s: %v", *addrFlag, err)
	}
	defer conn.Close()

	log.Printf("UDP-сервер запущен на %s", *addrFlag)
	buffer := make([]byte, 64*1024)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("прочитать датаграмму: %v", err)
			continue
		}

		message := string(buffer[:n])
		log.Printf("[%s] получено: %s", remoteAddr, message)
		response := fmt.Sprintf("Получено %d байт: %s", n, message)
		if _, err := conn.WriteToUDP([]byte(response), remoteAddr); err != nil {
			log.Printf("[%s] отправить ответ: %v", remoteAddr, err)
		}
	}
}
