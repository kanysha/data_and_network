package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
)

func handleConn(conn net.Conn) {
	defer conn.Close()

	remoteAddr := conn.RemoteAddr().String()
	log.Printf("подключён клиент %s", remoteAddr)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		log.Printf("[%s] получено: %s", remoteAddr, message)

		if _, err := fmt.Fprintf(conn, "ECHO: %s\n", strings.ToUpper(message)); err != nil {
			log.Printf("[%s] ошибка отправки ответа: %v", remoteAddr, err)
			return
		}
		if strings.EqualFold(message, "quit") {
			log.Printf("[%s] клиент завершил сеанс", remoteAddr)
			return
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("[%s] ошибка чтения: %v", remoteAddr, err)
	}
}

func main() {
	addr := flag.String("addr", ":9090", "TCP-адрес для прослушивания")
	flag.Parse()

	listener, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("запустить TCP-сервер на %s: %v", *addr, err)
	}
	defer listener.Close()

	log.Printf("TCP-сервер запущен на %s", *addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("принять соединение: %v", err)
			continue
		}
		go handleConn(conn)
	}
}
