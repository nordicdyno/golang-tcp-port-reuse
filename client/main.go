package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const (
	tcpEndpointEnvName = "TCP_CONNECTION_ENDPOINT"
)

var (
	argAddr    = flag.String("e", "", "tcp connection endpoint")
	sourceAddr = flag.String("s", "", "client address to bind")
	sleepSecs  = flag.Int("w", 10, "wait N seconds before exit")
	notClose   = flag.Bool("noclose", false, "not close connection behaviour")
)

func main() {
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

	fmt.Println(" exit")
	os.Exit(0)
}

func connect(addr string) {
	dialer := net.Dialer{}
	if *sourceAddr != "" {
		tcpAddr, err := net.ResolveTCPAddr("tcp", *sourceAddr)
		if err != nil {
			panic(err)
		}
		dialer.LocalAddr = tcpAddr
	}
	conn, err := dialer.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("failed connect to %v: %v", addr, err)
	}
	fmt.Print("connected")

	fmt.Printf(" waiting %v seconds\n", *sleepSecs)
	time.Sleep(time.Duration(*sleepSecs * int(time.Second)))
	if !*notClose {
		fmt.Println("close conn")
		conn.Close()
	}
}
