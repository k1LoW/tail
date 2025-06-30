// Package tail provides a writer that maintains the last N lines of written data,
// similar to the Unix tail command.
package tail

import (
	"bytes"
	"io"
	"strings"
	"sync"
)

// TailBuffer implements io.Writer and maintains the last N lines
// of written data.
type TailBuffer struct {
	mu       sync.Mutex
	maxLines int
	lines    []string
	buffer   bytes.Buffer
}

// New creates a new TailBuffer with the specified maximum number of lines.
func New(maxLines int) *TailBuffer {
	return &TailBuffer{
		maxLines: maxLines,
		lines:    make([]string, 0, maxLines),
	}
}

// Write implements the io.Writer interface.
// It writes data and maintains the last N lines.
func (tb *TailBuffer) Write(p []byte) (n int, err error) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	n = len(p)

	// Add to buffer
	tb.buffer.Write(p)

	// Split buffer content into lines
	content := tb.buffer.String()
	lines := strings.Split(content, "\n")

	// If the last element is not empty, it's not a complete line yet
	if len(lines) > 0 && lines[len(lines)-1] != "" {
		// Keep the last incomplete line in the buffer
		tb.buffer.Reset()
		tb.buffer.WriteString(lines[len(lines)-1])
		lines = lines[:len(lines)-1]
	} else {
		// Clear the buffer
		tb.buffer.Reset()
		// Remove the last empty element
		if len(lines) > 0 && lines[len(lines)-1] == "" {
			lines = lines[:len(lines)-1]
		}
	}

	// Don't keep any lines if maxLines is 0
	if tb.maxLines == 0 {
		tb.lines = []string{}
	} else {
		// Add new lines
		tb.lines = append(tb.lines, lines...)

		// Remove old lines if exceeding maxLines
		if len(tb.lines) > tb.maxLines {
			tb.lines = tb.lines[len(tb.lines)-tb.maxLines:]
		}
	}

	return n, nil
}

// Lines returns the maintained lines as a slice.
func (tb *TailBuffer) Lines() []string {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	// Consider any unprocessed data in the current buffer
	result := make([]string, len(tb.lines))
	copy(result, tb.lines)

	// Add any remaining data in the buffer as the last line
	if tb.buffer.Len() > 0 {
		result = append(result, tb.buffer.String())
		// Adjust if exceeding maxLines
		if tb.maxLines > 0 && len(result) > tb.maxLines {
			result = result[len(result)-tb.maxLines:]
		}
	}

	return result
}

// String returns the maintained lines joined with newlines as a string.
func (tb *TailBuffer) String() string {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	if len(tb.lines) == 0 && tb.buffer.Len() == 0 {
		return ""
	}

	// Create a copy of lines
	result := make([]string, len(tb.lines))
	copy(result, tb.lines)

	// Check if there's data in buffer
	hasTrailingNewline := false
	if tb.buffer.Len() > 0 {
		result = append(result, tb.buffer.String())
		// Adjust if exceeding maxLines
		if tb.maxLines > 0 && len(result) > tb.maxLines {
			result = result[len(result)-tb.maxLines:]
		}
	} else if len(tb.lines) > 0 {
		// If buffer is empty, it means the last write ended with a newline
		hasTrailingNewline = true
	}

	str := strings.Join(result, "\n")
	if hasTrailingNewline && len(result) > 0 {
		str += "\n"
	}
	return str
}

// Bytes returns the maintained lines joined with newlines as a byte slice.
func (tb *TailBuffer) Bytes() []byte {
	return []byte(tb.String())
}

// WriteTo implements the io.WriterTo interface.
// It writes the maintained lines to the specified Writer.
func (tb *TailBuffer) WriteTo(w io.Writer) (n int64, err error) {
	data := tb.Bytes()
	written, err := w.Write(data)
	return int64(written), err
}
