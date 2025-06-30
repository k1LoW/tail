# tail

`tail` provides a writer maintaining the last N lines of written data, similar to the Unix `tail` command.

## Usage

```go
package main

import (
    "fmt"
    "github.com/k1LoW/tail"
)

func main() {
    tb := tail.New(3)

    // Write in multiple chunks
    tb.Write([]byte("5\n"))
    tb.Write([]byte("4\n"))
    tb.Write([]byte("3\n"))
    tb.Write([]byte("2\n"))
    tb.Write([]byte("1\n"))
    tb.Write([]byte("Hello "))
    tb.Write([]byte("World\n"))
    tb.Write([]byte("Foo"))
    tb.Write([]byte("Bar\n"))

    fmt.Println(tb.String())
    // Output:
    // 1
    // Hello World
    // FooBar
}
```
