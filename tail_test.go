package tail

import (
	"bytes"
	"strings"
	"testing"
)

func TestTailBuffer_Write(t *testing.T) {
	tests := []struct {
		name     string
		limit    int
		writes   []string
		expected []string
	}{
		{
			name:     "basic write and read",
			limit:    3,
			writes:   []string{"line1\n"},
			expected: []string{"line1"},
		},
		{
			name:     "exceed line limit",
			limit:    3,
			writes:   []string{"line1\n", "line2\n", "line3\n", "line4\n", "line5\n"},
			expected: []string{"line3", "line4", "line5"},
		},
		{
			name:     "multiple writes",
			limit:    5,
			writes:   []string{"line1\nli", "ne2\nline3", "\nline4\n"},
			expected: []string{"line1", "line2", "line3", "line4"},
		},
		{
			name:     "empty lines",
			limit:    5,
			writes:   []string{"line1\n\nline3\n\n\nline6\n"},
			expected: []string{"", "line3", "", "", "line6"},
		},
		{
			name:     "long lines",
			limit:    3,
			writes:   []string{"short1\n", strings.Repeat("a", 1000) + "\n", "short2\n", "short3\n"},
			expected: []string{strings.Repeat("a", 1000), "short2", "short3"},
		},
		{
			name:     "no newline at end",
			limit:    3,
			writes:   []string{"line1\nline2\nline3"},
			expected: []string{"line1", "line2", "line3"},
		},
		{
			name:     "zero lines",
			limit:    0,
			writes:   []string{"line1\nline2\nline3\n"},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tw := New(tt.limit)

			// Write all data
			for _, data := range tt.writes {
				n, err := tw.Write([]byte(data))
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if n != len(data) {
					t.Errorf("expected %d bytes written, got %d", len(data), n)
				}
			}

			// Check result
			result := tw.Lines()
			if len(result) != len(tt.expected) {
				t.Errorf("expected %d lines, got %d", len(tt.expected), len(result))
			}

			for i, line := range result {
				if i < len(tt.expected) && line != tt.expected[i] {
					t.Errorf("line %d: expected '%s', got '%s'", i, tt.expected[i], line)
				}
			}
		})
	}
}

func TestTailBuffer_Output(t *testing.T) {
	tests := []struct {
		name           string
		limit          int
		input          string
		expectedString string
		expectedBytes  []byte
	}{
		{
			name:           "string method",
			limit:          3,
			input:          "line1\nline2\nline3\nline4\n",
			expectedString: "line2\nline3\nline4\n",
			expectedBytes:  []byte("line2\nline3\nline4\n"),
		},
		{
			name:           "no newline at end",
			limit:          3,
			input:          "line1\nline2\nline3\nli",
			expectedString: "line2\nline3\nli",
			expectedBytes:  []byte("line2\nline3\nli"),
		},
		{
			name:           "empty buffer",
			limit:          3,
			input:          "",
			expectedString: "",
			expectedBytes:  []byte{},
		},
		{
			name:           "single line without newline",
			limit:          3,
			input:          "single",
			expectedString: "single",
			expectedBytes:  []byte("single"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tw := New(tt.limit)
			if _, err := tw.Write([]byte(tt.input)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Test String()
			if result := tw.String(); result != tt.expectedString {
				t.Errorf("String(): expected '%s', got '%s'", tt.expectedString, result)
			}

			// Test Bytes()
			if result := tw.Bytes(); !bytes.Equal(result, tt.expectedBytes) {
				t.Errorf("Bytes(): expected %v, got %v", tt.expectedBytes, result)
			}
		})
	}
}

func TestTailBuffer_WriteTo(t *testing.T) {
	tests := []struct {
		name          string
		limit         int
		input         string
		expected      string
		expectedBytes int64
	}{
		{
			name:          "basic writeto",
			limit:         3,
			input:         "line1\nline2\nline3\nline4\n",
			expected:      "line2\nline3\nline4\n",
			expectedBytes: 18,
		},
		{
			name:          "empty buffer writeto",
			limit:         3,
			input:         "",
			expected:      "",
			expectedBytes: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tw := New(tt.limit)
			if _, err := tw.Write([]byte(tt.input)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			var buf bytes.Buffer
			n, err := tw.WriteTo(&buf)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if buf.String() != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, buf.String())
			}

			if n != tt.expectedBytes {
				t.Errorf("expected %d bytes written, got %d", tt.expectedBytes, n)
			}
		})
	}
}

func TestTailBuffer_ConcurrentWrites(t *testing.T) {
	tw := New(100)
	done := make(chan bool)

	// Write concurrently from 10 goroutines
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 10; j++ {
				_, _ = tw.Write([]byte(strings.Repeat("a", id+1) + "\n"))
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	result := tw.Lines()
	if len(result) != 100 {
		t.Errorf("expected 100 lines, got %d", len(result))
	}
}

// Benchmark tests
func BenchmarkTailBuffer_Write(b *testing.B) {
	tw := New(1000)
	data := []byte("This is a benchmark test line\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = tw.Write(data)
	}
}

func BenchmarkTailBuffer_WriteLongLines(b *testing.B) {
	tw := New(100)
	data := []byte(strings.Repeat("x", 1000) + "\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = tw.Write(data)
	}
}
