package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"spotcapture/internal/buffer"
	"spotcapture/internal/capture"
	"spotcapture/internal/storage"
	"spotcapture/internal/upload"
	"syscall"
	"time"

	"github.com/spf13/viper"
)

func LoadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	return viper.ReadInConfig()
}

func main() {
	log.Print("Spotcapture: loading config...")

	LoadConfig()

	stopChan := make(chan struct{})
	ring := buffer.NewRingBuffer(1000)

	log.Print("Spotcapture: starting capturing...")
	go capture.StartCapture(viper.GetString("interface"), ring, stopChan)

	go listenTCP(viper.GetInt("port"), stopChan)
	go listenUDP(viper.GetInt("port"), stopChan)

	log.Print("Spotcapture: capturing in progress")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigChan:
		log.Print("Spotcapture: Stopping due to signal")
	case <-stopChan:
		log.Print("Spotcapture: Stopping due to TCP/UDP trigger")
	}

	log.Print("Spotcapture: Sapturing stopped")

	packets := ring.Packets()
	pcap, err := storage.SaveRingBufferToPcapBuffer(packets)
	if err != nil {
		log.Fatalf("Failed to save PCAP: %v", err)
	}

	filename := fmt.Sprintf("file_%s.pcap", time.Now().Format("20060102_150405"))
	log.Println("Spotcapture: starting upload...")
	err = upload.Upload(context.TODO(), pcap, filename)
	if err != nil {
		log.Fatalf("Failed to upload to S3: %v", err)
	}

	log.Println("Spotcapture: Upload complete!")
}

func listenTCP(port int, stopChan chan struct{}) {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("TCP listener error: %v", err)
		return
	}
	defer listener.Close()
	log.Printf("Spotcapture: Listening for TCP on %s", addr)

	conn, err := listener.Accept()
	if err == nil {
		log.Printf("Spotcapture: TCP connection detected from %s", conn.RemoteAddr())
		close(stopChan)
		conn.Close()
	}
}

func listenUDP(port int, stopChan chan struct{}) {
	addr := fmt.Sprintf(":%d", port)
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		log.Printf("UDP listener error: %v", err)
		return
	}
	defer conn.Close()
	log.Printf("Spotcapture: Listening for UDP on %s", addr)

	buf := make([]byte, 1024)
	_, remoteAddr, err := conn.ReadFrom(buf)
	if err == nil {
		log.Printf("Spotcapture: UDP packet detected from %s", remoteAddr)
		close(stopChan)
	}
}
