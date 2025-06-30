package tail_test

import (
	"fmt"
	"log"

	"github.com/k1LoW/tail"
)

func ExampleTailBuffer() {
	// Create a TailBuffer that keeps the last 3 lines
	tw := tail.New(3)

	// Write data
	if _, err := tw.Write([]byte("Line 1\n")); err != nil {
		log.Fatal(err)
	}
	if _, err := tw.Write([]byte("Line 2\n")); err != nil {
		log.Fatal(err)
	}
	if _, err := tw.Write([]byte("Line 3\n")); err != nil {
		log.Fatal(err)
	}
	if _, err := tw.Write([]byte("Line 4\n")); err != nil {
		log.Fatal(err)
	}
	if _, err := tw.Write([]byte("Line 5\n")); err != nil {
		log.Fatal(err)
	}

	// Get the last 3 lines
	lines := tw.Lines()
	for _, line := range lines {
		fmt.Println(line)
	}
	// Output:
	// Line 3
	// Line 4
	// Line 5
}

func ExampleTailBuffer_multipleWrites() {
	tw := tail.New(2)

	// Write in multiple chunks
	if _, err := tw.Write([]byte("Hello ")); err != nil {
		log.Fatal(err)
	}
	if _, err := tw.Write([]byte("World\n")); err != nil {
		log.Fatal(err)
	}
	if _, err := tw.Write([]byte("Foo")); err != nil {
		log.Fatal(err)
	}
	if _, err := tw.Write([]byte("Bar\n")); err != nil {
		log.Fatal(err)
	}

	fmt.Println(tw.String())
	// Output:
	// Hello World
	// FooBar
}

func ExampleTailBuffer_asWriter() {
	tw := tail.New(5)

	// Use as io.Writer
	logger := log.New(tw, "", 0)

	for i := 1; i <= 10; i++ {
		logger.Printf("Log entry %d", i)
	}

	// Only the last 5 entries are retained
	fmt.Println("Last 5 log entries:")
	for _, line := range tw.Lines() {
		fmt.Println(line)
	}
	// Output:
	// Last 5 log entries:
	// Log entry 6
	// Log entry 7
	// Log entry 8
	// Log entry 9
	// Log entry 10
}
