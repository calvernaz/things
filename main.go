package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/calvernaz/things/serial"
	"github.com/calvernaz/things/sse"
)

func main() {

	log.Println("Starting application...")

	/*
	 * When SIGINT or SIGTERM is caught write to the quitChannel
	 */
	quitChannel := make(chan os.Signal)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)

	shutdownUdpBroadcastCh := make(chan bool)
	shutdownSerialBroadcastCh := make(chan bool)
	msgChannel := make(chan string)
	sseChannel := make(chan string)

	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(1)

	// go udp.ListenUdp(shutdownChannel, waitGroup)
	//go udp.SendUdp(shutdownUdpBroadcastCh, msgChannel, waitGroup)
	go serial.ReadSerial(shutdownSerialBroadcastCh, msgChannel, sseChannel, waitGroup)
	go sse.Main(sseChannel)

	/*
	 * Wait until we get the quit message
	 */
	<-quitChannel

	log.Println("Received quit. Sending shutdown and waiting on goroutines...")
	shutdownUdpBroadcastCh <- true
	shutdownSerialBroadcastCh <- true

	/*
	 * Block until wait group counter gets to zero
	 */
	waitGroup.Wait()
	log.Println("Done.")
}
