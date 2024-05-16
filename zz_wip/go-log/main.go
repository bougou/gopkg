package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// goID return the ID of the goroutine
func goID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}

func init() {
	log.SetFlags(log.Lmicroseconds)
}

func DebugPrint(format string, args ...interface{}) {
	if os.Getenv("DEBUG") == "true" {
		v := []interface{}{goID()}
		v = append(v, args...)
		log.Printf("debug_goroutine:%02d: "+format, v...)
	}
}

func main() {
	item := "apples"
	capacity := 10
	size := 100
	DebugPrint("adding (%s), capacity=%d, size=%d", item, capacity, size)

}
