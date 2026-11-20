package main

import (
	"errors"
	"sync"
	"time"
)

// StreamBuffer manages chunked video streaming buffers with backpressure.
type StreamBuffer struct {
	mu          sync.RWMutex
	maxCapacity int
	chunks      [][]byte
	lastActive  time.Time
}

// NewStreamBuffer initializes a thread-safe stream buffer.
func NewStreamBuffer(maxCapacity int) *StreamBuffer {
	return &StreamBuffer{
		maxCapacity: maxCapacity,
		chunks:      make([][]byte, 0, maxCapacity),
		lastActive:  time.Now(),
	}
}

// Push appends a chunk or returns an error if high-watermark backpressure triggers.
func (sb *StreamBuffer) Push(chunk []byte) error {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	if len(sb.chunks) >= sb.maxCapacity {
		return errors.New("backpressure: buffer capacity exceeded")
	}

	sb.chunks = append(sb.chunks, chunk)
	sb.lastActive = time.Now()
	return nil
}

// Flush clears all accumulated chunks and returns the count of drained bytes.
func (sb *StreamBuffer) Flush() int {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	totalBytes := 0
	for _, c := range sb.chunks {
		totalBytes += len(c)
	}
	sb.chunks = sb.chunks[:0]
	return totalBytes
}

// Len returns the current buffered chunk count.
func (sb *StreamBuffer) Len() int {
	sb.mu.RLock()
	defer sb.mu.RUnlock()
	return len(sb.chunks)
}
