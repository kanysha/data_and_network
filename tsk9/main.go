package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func runDNSExercise(domain string) {
	fmt.Printf("=== DNS Exercise for %s ===\n", domain)

	ips, err := net.LookupHost(domain)
	if err != nil {
		fmt.Printf("Ошибка DNS (IPv4/IPv6): %v\n", err)
	} else {
		fmt.Println("IP-адреса:")
		for _, ip := range ips {
			fmt.Printf("  %s\n", ip)
		}
	}

	mxs, err := net.LookupMX(domain)
	if err != nil {
		fmt.Printf("Ошибка DNS (MX): %v\n", err)
	} else {
		fmt.Println("\nMX-записи:")
		for _, mx := range mxs {
			fmt.Printf("  %s (приоритет: %d)\n", mx.Host, mx.Pref)
		}
	}

	nss, err := net.LookupNS(domain)
	if err != nil {
		fmt.Printf("Ошибка DNS (NS): %v\n", err)
	} else {
		fmt.Println("\nNS-записи:")
		for _, ns := range nss {
			fmt.Printf("  %s\n", ns.Host)
		}
	}

	txts, err := net.LookupTXT(domain)
	if err != nil {
		fmt.Printf("Ошибка DNS (TXT): %v\n", err)
	} else {
		fmt.Println("\nTXT-записи:")
		for _, txt := range txts {
			fmt.Printf("  %s\n", txt)
		}
	}
	fmt.Println()
}

func runInterfaceExercise() {
	fmt.Println("=== Network Interfaces ===")

	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("Ошибка получения интерфейсов: %v\n", err)
		return
	}

	for _, iface := range interfaces {
		fmt.Printf("Интерфейс: %s\n", iface.Name)
		fmt.Printf("  MAC: %s\n", iface.HardwareAddr)
		fmt.Printf("  MTU: %d\n", iface.MTU)

		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Printf("  Ошибка получения адресов: %v\n", err)
			continue
		}

		for _, addr := range addrs {
			fmt.Printf("  Адрес: %s\n", addr)
		}
		fmt.Println()
	}
}

func checkHost(host string, port string) {
	address := net.JoinHostPort(host, port)
	start := time.Now()
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("[FAIL] %s — %v (%v)\n", address, err, duration)
		return
	}
	defer conn.Close()

	fmt.Printf("[ OK ] %s — доступен (%v)\n", address, duration)
}

func runConnectivityExercise() {
	fmt.Println("=== Connectivity Check ===")

	hosts := []struct {
		host string
		port string
	}{
		{"google.com", "80"},
		{"github.com", "443"},
		{"localhost", "8080"},
	}

	for _, h := range hosts {
		checkHost(h.host, h.port)
	}
	fmt.Println()
}

func main() {
	domain := "google.com"
	if len(os.Args) > 1 {
		domain = os.Args[1]
	}

	runDNSExercise(domain)
	runInterfaceExercise()
	runConnectivityExercise()
}
