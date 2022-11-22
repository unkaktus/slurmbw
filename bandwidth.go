package slurmbw

import (
	"crypto/rand"
	"fmt"
	"log"
	"net"
	"time"
)

func handleConnection(conn net.Conn) {
	buffer := make([]byte, 1024*1024*10)
	_, _ = rand.Reader.Read(buffer)
	for {
		start := time.Now()
		n, err := conn.Write(buffer)
		elapsed := time.Since(start)
		if err != nil {
			log.Printf("write: %v", err)
			return
		}
		speed := float64(n) / elapsed.Seconds()
		speedGbps := speed * 8 / (1024 * 1024 * 1024)
		log.Printf("speed: %v Gbps", speedGbps)
	}
}

func Listen(address string) error {
	listener, err := net.Listen("tcp4", address)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}
