package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

const (
	tcpListenAddrEnvName = "TCP_CONNECTION_LISTEN_ADDR"
)


func main() {
	argAddr := flag.String("l", "", "listen tcp address")
	sleepSecs := flag.Int("s", 3, "wait before exit after first incoming connection")
	flag.Parse()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Println()
		fmt.Println("  supported environment variables:\n")
		fmt.Printf("    %v overwrites -l flag.\n", tcpListenAddrEnvName)
		fmt.Println()
	}
	// _ = addr
	addr := os.Getenv(tcpListenAddrEnvName)
	if len(addr) == 0 {
		addr = *argAddr
	}
	// fmt.Fprintf(os.Stderr, "Error: at least one port should be used\n\n")
	if len(addr) == 0 {
		flag.Usage()
		os.Exit(1)

	}
	fmt.Println("TCP bind to:", addr)

	listen(addr)
	fmt.Printf(" waiting %v seconds\n", *sleepSecs)
	time.Sleep(time.Duration(*sleepSecs*int(time.Second)))
	fmt.Println(" exit")
	os.Exit(0)
}

func listen(address string) {
	addr, err := net.ResolveTCPAddr("tcp4", address)
	// TODO: return nil, error and decide how to handle it in the calling function
	if err != nil {
		fmt.Println("Failed to resolve address", err.Error())
		os.Exit(1)
	}

	listener, err := net.Listen("tcp", addr.String())
	if err != nil {
		fmt.Println("Failed to", err.Error())
		os.Exit(1)
	}

	for {
		_, err := listener.Accept()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			fmt.Println("Failed to accept connection:", err.Error())
			continue
		}
		fmt.Println("got connection")
		return
	}
}