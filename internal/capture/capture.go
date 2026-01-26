package capture

import (
	"log"
	"net"
	"os"
	"spotcapture/internal/buffer"
	"time"

	"github.com/mdlayher/packet"
	"golang.org/x/sys/unix"
)

func StartCapture(device string, ring *buffer.RingBuffer, stopChan <-chan struct{}) {
	iface, err := net.InterfaceByName(device)
	if err != nil {
		log.Fatalf("Error getting interface: %v", err)
	}

	conn, err := packet.Listen(iface, packet.Raw, unix.ETH_P_ALL, nil)
	if err != nil {
		log.Fatalf("Error creating packet socket: %v", err)
	}
	defer conn.Close()

	buf := make([]byte, 65535)

	log.Println("Starting capture on", device)
	for {
		select {
		case <-stopChan:
			log.Println("Stopping capture...")
			return
		default:
			// Set a short read deadline to avoid blocking forever
			conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))

			n, _, err := conn.ReadFrom(buf)
			if err != nil {
				if os.IsTimeout(err) {
					continue // Check stopChan again
				}
				log.Printf("Read error: %v", err)
				continue
			}
			ring.Add(buf[:n])
		}
	}
}
