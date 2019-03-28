package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
	"log"
)

const (
	tcpEndpointEnvName = "TCP_CONNECTION_ENDPOINT"
)


func main() {
	argAddr := flag.String("e", "", "tcp connection endpoint")
	sleepSecs := flag.Int("s", 1, "sleep before stop")
	flag.Parse()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Println()
		fmt.Println("  supported environment variables:\n")
		fmt.Printf("    %v overwrites -e flag.\n", tcpEndpointEnvName)
		fmt.Println()
	}
	// _ = addr
	addr := os.Getenv(tcpEndpointEnvName)
	if len(addr) == 0 {
		addr = *argAddr
	}
	// fmt.Fprintf(os.Stderr, "Error: at least one port should be used\n\n")
	if len(addr) == 0 {
		flag.Usage()
		os.Exit(1)

	}
	fmt.Println("connect to tcp endpoint:", addr)

	connect(addr)
	fmt.Printf(" waiting %v seconds\n", *sleepSecs)
	time.Sleep(time.Duration(*sleepSecs*int(time.Second)))
	fmt.Println(" exit")
	os.Exit(0)
}

func connect(addr string) {
	_, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("failed connect to %v: %v", addr, err)
	}
	fmt.Print("connected")
}