package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
	"time"
)

const (
	tcpListenAddrEnvName = "TCP_CONNECTION_LISTEN_ADDR"
)

var (
	argAddr     = flag.String("l", "", "listen tcp address")
	sleepSecs   = flag.Int("w", 1, "wait before exit after first incoming connection")
	noClose     = flag.Bool("no-close", false, "not close connection behaviour")
	noReuseAddr = flag.Bool("no-reuse-addr", false, "not close connection behaviour")
)

func main() {
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
}

func listen(address string) {
	addr, err := net.ResolveTCPAddr("tcp4", address)
	// TODO: return nil, error and decide how to handle it in the calling function
	if err != nil {
		fmt.Println("Failed to resolve address", err.Error())
		os.Exit(1)
	}

	lCfg := &net.ListenConfig{}
	lCfg.Control = func(network, address string, c syscall.RawConn) error {
		fmt.Println("ListenConfig")
		fmt.Println("network:", network)
		fmt.Println("address:", address)
		var fn = func(s uintptr) {
			if *noReuseAddr {
				setErr := syscall.SetsockoptInt(int(s), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 0)
				if setErr != nil {
					log.Fatal(setErr)
				}
			}

			val, getErr := syscall.GetsockoptInt(int(s), syscall.SOL_SOCKET, syscall.SO_REUSEADDR)
			if getErr != nil {
				log.Fatal(getErr)
			}
			log.Printf("value of SO_REUSEADDR option is: %d", int(val))
		}
		if err := c.Control(fn); err != nil {
			return err
		}
		return nil
	}

	listener, err := lCfg.Listen(context.Background(), "tcp", addr.String())
	// listener, err := net.Listen("tcp", addr.String())
	if err != nil {
		fmt.Println("Failed to", err.Error())
		os.Exit(1)
	}

	go func() {
		for {
			_, err := listener.Accept()
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					fmt.Println("got error:", opErr)
				}
				fmt.Println("Failed to accept connection:", err.Error())
				continue
			}
			fmt.Println("got connection")
		}
	}()

	fmt.Printf(" waiting %v seconds\n", *sleepSecs)
	time.Sleep(time.Duration(*sleepSecs * int(time.Second)))

	if !*noClose {
		fmt.Println("close conn")
		listener.Close()
	}

	fmt.Println(" exit")
	os.Exit(0)
}
