package sse

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/calvernaz/things/encode"
)

type Broker struct {

	// Create a map of clients, the keys of the map are the channels
	// over which we can push messages to attached clients.  (The values
	// are just booleans and are meaningless.)
	//
	clients map[chan string]bool

	// Channel into which new clients can be pushed
	//
	newClients chan chan string

	// Channel into which disconnected clients should be pushed
	//
	defunctClients chan chan string

	// Channel into which messages are pushed to be broadcast out
	// to attached clients.
	//
	messages chan string
}

func (b *Broker) Start() {

	go func() {
		for {
			select {
			case s := <-b.newClients:
				// There is a new client attached and we
				// want to start sending them messages.
				b.clients[s] = true
				log.Println("Added new client")
			case s := <-b.defunctClients:
				// A client has dettached and we want to
				// stop sending them messages.
				delete(b.clients, s)
				close(s)
				log.Println("Removed client")
			case msg := <-b.messages:
				// There is a new message to send.  For each
				// attached client, push the new message
				// into the client's message channel.
				for s := range b.clients {
					s <- msg
				}
				log.Printf("Broadcast message to %d clients", len(b.clients))
			}
		}
	}()
}

func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Make sure that the writer supports flushing.
	//
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Create a new channel, over which the broker can
	// send this client messages.
	messageChan := make(chan string)

	// Add this client to the map of those that should
	// receive updates
	b.newClients <- messageChan

	// Set the headers related to event streaming.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Listen to the closing of the http connection via the CloseNotifier
	cn, ok := w.(http.CloseNotifier)
	if !ok {
		http.Error(w, "cannot stream", http.StatusInternalServerError)
		return
	}

	for {
		select {
		case <-cn.CloseNotify():
			b.defunctClients <- messageChan
			log.Println("HTTP connection just closed.")
			return
		case msg := <-messageChan:
			fmt.Fprint(w, "event: temp\n")
			fmt.Fprintf(w, "data: %s\n\n", msg)
			f.Flush()
		}
	}

	// Done.
	log.Println("Finished HTTP request at ", r.URL.Path)
}

func Main(sseChannel chan string) {
	log.Println("Starting work goroutine: Main")
	b := &Broker{
		make(map[chan string]bool),
		make(chan (chan string)),
		make(chan (chan string)),
		make(chan string),
	}

	// Start processing events
	b.Start()

	// Make b the HTTP handler for "/events/".  It can do
	// this because it has a ServeHTTP method.  That method
	// is called in a separate goroutine for each
	// request to "/events/".
	http.Handle("/events/", b)

	go func() {
		for {
			select {
			case msg := <-sseChannel:
				js, _ := json.Marshal(encode.Encode(msg, time.Now()))
				log.Printf("Pushing messages to clients: %v", string(js))
				b.messages <- string(js)
			}
		}
	}()

	// When we get a request at "/", call `MainPageHandler`
	// in a new goroutine.
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./html"))))
	//	http.Handle("/", http.FileServer(http.Dir("./html")))

	// Start the server and listen forever on port 8000.
	fmt.Println("Serving at :8000")
	http.ListenAndServeTLS("0.0.0.0:8000", "./server.crt", "./server.key", nil)
}
