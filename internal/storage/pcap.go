package storage

import (
	"bytes"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
)

// SaveRingBufferToPcapBuffer writes packets to a PCAP buffer and returns it as []byte
func SaveRingBufferToPcapBuffer(packets [][]byte) ([]byte, error) {
	buf := &bytes.Buffer{}
	w := pcapgo.NewWriter(buf)

	// Write PCAP file header: Ethernet link type
	if err := w.WriteFileHeader(65535, layers.LinkTypeEthernet); err != nil {
		return nil, err
	}

	for _, pkt := range packets {
		if pkt != nil {
			ci := gopacket.CaptureInfo{
				Timestamp:     time.Now(),
				CaptureLength: len(pkt),
				Length:        len(pkt),
			}
			if err := w.WritePacket(ci, pkt); err != nil {
				return nil, err
			}
		}
	}

	return buf.Bytes(), nil
}

// SaveRingBufferToPcapFile writes packets to a PCAP file using the buffer function
func SaveRingBufferToPcapFile(filename string, packets [][]byte) error {
	data, err := SaveRingBufferToPcapBuffer(packets)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}
