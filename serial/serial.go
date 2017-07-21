package serial

import (
	"fmt"
	"log"
	"net"
	"sync"

	"bufio"
	"strings"

	"go.bug.st/serial.v1"
)

func ReadSerial(shutdownSerialBroadcastCh chan bool, msg2UdpChannel chan string, msg2SseCh chan string, waitGroup *sync.WaitGroup) {
	log.Println("Starting work goroutine: ReadSerial")
	defer waitGroup.Done()

	port, err := serial.Open("/dev/ttyAMA0", &serial.Mode{BaudRate: 9600})
	if err != nil {
		log.Fatal(err)
	}

	var n byte
	r := bufio.NewReader(port)
	for {

		if n, err = r.ReadByte(); nil != err {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			log.Fatal(err)
			break
		}

		if n == 'a' {
			lot := "a"
			count := 0
			for count < 11 {
				n, err = r.ReadByte()
				if err != nil {
					return
				}
				if n == 'a' {
					count = 0
					lot = "a"
				} else if count == 0 || count == 1 {
					lot += string(n)
					count += 1
				} else if count >= 2 {
					lot += string(n)
					count += 1
				} else {
					fmt.Printf("RX: %v\n", lot[1:]+string(n))
					break
				}
			}

			fmt.Printf("RX:%v\n", lot[1:])

			if len(lot) == 12 {
				if lot[1:3] == "??" {
					lot = strings.Replace(lot[3:], "-", "", -1)
				} else {
					select {
					case <-shutdownSerialBroadcastCh:
						log.Println("ReadSerial: Received shutdown on goroutine")
						return
					default:

						msg := lot[1:]
						//msg2UdpChannel <- msg
						msg2SseCh <- msg
					}
				}
			}
		}
	}
}
