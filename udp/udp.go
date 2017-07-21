package udp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/calvernaz/things/encode"
)

func ListenUdp(shutdownChannel chan bool, waitGroup *sync.WaitGroup) {
	log.Println("Starting work goroutine: ListenUdp")
	defer waitGroup.Done()

	timeoutDuration := 3 * time.Second

	addr := net.UDPAddr{
		Port: 50141,
		IP:   net.ParseIP("127.0.0.1"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 100)
	r := bufio.NewReader(conn)
	for {
		/*
		 * Listen on channels for message.
		 */
		conn.SetReadDeadline(time.Now().Add(timeoutDuration))
		select {
		case <-shutdownChannel:
			log.Printf("ListenUdp: Received shutdown on goroutine\n")
			break
		default:

		}

		var n int
		if n, err = r.Read(buf); nil != err {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			log.Println(err)
			break
		}
		fmt.Println(string(buf[:n]))
		r.Reset(conn)
	}
}

func SendUdp(shutdownUdpBroadcastCh chan bool, msgChannel chan string, waitGroup *sync.WaitGroup) error {
	log.Println("Starting work goroutine: SendUdp")
	defer waitGroup.Done()

	BROADCAST_IPv4 := net.IPv4(255, 255, 255, 255)
	socket, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   BROADCAST_IPv4,
		Port: 50140,
	})
	defer socket.Close()
	if err != nil {
		return err
	}

	for msgChannel != nil {
		select {
		case v, ok := <-msgChannel:
			if !ok {
				msgChannel = nil
				continue
			}
			jsonMsg, err := json.Marshal(encode.Encode(v, time.Now()))
			if err != nil {
				continue
			}
			socket.Write(jsonMsg)
		case <-shutdownUdpBroadcastCh:
			log.Println("SendUdp: Received shutdown on goroutine")
			return nil
		}
	}
	return nil
}
