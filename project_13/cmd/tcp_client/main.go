package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	addr := flag.String("addr", "localhost:9090", "адрес TCP-сервера")
	flag.Parse()

	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		log.Fatalf("подключиться к %s: %v", *addr, err)
	}
	defer conn.Close()

	fmt.Println("Подключено к серверу. Введите текст (quit — выход):")
	input := bufio.NewScanner(os.Stdin)
	responses := bufio.NewScanner(conn)
	for input.Scan() {
		message := input.Text()
		if _, err := fmt.Fprintf(conn, "%s\n", message); err != nil {
			log.Fatalf("отправить сообщение: %v", err)
		}
		if !responses.Scan() {
			if err := responses.Err(); err != nil {
				log.Fatalf("прочитать ответ: %v", err)
			}
			log.Fatal("сервер закрыл соединение без ответа")
		}
		fmt.Println("Сервер:", responses.Text())
		if strings.EqualFold(message, "quit") {
			return
		}
	}
	if err := input.Err(); err != nil {
		log.Fatalf("прочитать ввод: %v", err)
	}
}
