package main

import (
	"fmt"
	"os"
)

func main() {
	f, err := os.Open("/dev/input/event0")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b := make([]byte, 200)
	for {
		n, _ := f.Read(b)
		fmt.Println(n)
		// sec := binary.LittleEndian.Uint64(b[0:8])
		// usec := binary.LittleEndian.Uint64(b[8:16])
		// t := time.Unix(int64(sec), int64(usec))
		// fmt.Println(t)
		// var value int32
		// typ := binary.LittleEndian.Uint16(b[16:18])
		// code := binary.LittleEndian.Uint16(b[18:20])
		// binary.Read(bytes.NewReader(b[20:]), binary.LittleEndian, &value)
		// fmt.Printf("type: %x\ncode: %d\nvalue: %d\n", typ, code, value)
	}
}
