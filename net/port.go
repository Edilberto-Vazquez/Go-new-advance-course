package main

import (
	"fmt"
	"net"
)

func main() {
	for i := 0; i < 100; i++ {
		conexion, err := net.Dial("tcp", fmt.Sprintf("%s:%d", "scanme.nmap.org", i))
		if err != nil {
			continue
		}
		conexion.Close()
		fmt.Printf("Port %d is open", i)
	}
}
