package buffer

import "sync"

type RingBuffer struct {
	packets [][]byte
	size    int
	index   int
	mu      sync.Mutex
}

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		packets: make([][]byte, size),
		size:    size,
	}
}

func (rb *RingBuffer) Add(packet []byte) {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	rb.packets[rb.index] = packet
	rb.index = (rb.index + 1) % rb.size
}

func (rb *RingBuffer) Packets() [][]byte {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	return append([][]byte(nil), rb.packets...)
}

func (rb *RingBuffer) Dump() []byte {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	var out []byte
	for _, pkt := range rb.packets {
		if pkt != nil {
			out = append(out, pkt...)
			out = append(out, '\n')
		}
	}
	return out
}
