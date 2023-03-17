package main

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"

	"github.com/iugstav/sse_queue/queue"
)

type IncomingLogs struct {
	Data []queue.PriorityQueueElement `json:"data"`
}

type SSEBroker struct {
	Clients  map[chan string]string
	Notifier chan string
	clientWG *sync.WaitGroup
}

const letterBytes string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomToken() string {
	b := make([]byte, 15)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func main() {
	var fila queue.PriorityQueue
	heap.Init(&fila)

	sse := SSEBroker{
		Clients:  make(map[chan string]string),
		Notifier: make(chan string, 1),
		clientWG: &sync.WaitGroup{},
	}

	done := make(chan struct{})
	defer close(done)

	go logsBroadcaster(&sse, &fila, done)

	mux := http.NewServeMux()
	mux.HandleFunc("/income", handleIncomingLogs(&fila, &sse))
	mux.HandleFunc("/logs", handleOutputingLogs(&fila, &sse))
	http.ListenAndServe(":8080", mux)
}

func logsBroadcaster(broker *SSEBroker, q *queue.PriorityQueue, done chan struct{}) {
	for {
		select {
		case data := <-broker.Notifier:
			switch data {
			case "incoming.logs.notify":
				for c := range broker.Clients {
					c <- data
				}

			case "output.logs.clean":
				broker.clientWG.Wait()
				go func() {
					for q.Len() != 0 {
						heap.Pop(q)
					}
				}()
			}

		case <-done:
			return
		}
	}
}

func handleIncomingLogs(q *queue.PriorityQueue, b *SSEBroker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data IncomingLogs

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for _, el := range data.Data {
			heap.Push(q, el)
		}

		b.Notifier <- "incoming.logs.notify"
	}
}

func handleOutputingLogs(q *queue.PriorityQueue, b *SSEBroker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestCtx := r.Context()

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Connection doesnot support streaming", http.StatusBadRequest)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		clientChan := make(chan string)
		clientToken := randomToken()

		b.Clients[clientChan] = clientToken
		defer func() {
			close(clientChan)
			delete(b.Clients, clientChan)
		}()

		for {
			select {
			case <-clientChan:
				fmt.Printf("[LOG]: message to %s\n", clientToken)
				b.clientWG.Add(1)

				for _, elem := range *q {
					fmt.Println("queue element")

					data, err := json.Marshal(elem)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						fmt.Println("deu ruim")
						return
					}

					fmt.Fprintf(w, "data: %v\n\n", string(data))

				}

				b.clientWG.Done()
				flusher.Flush()
				b.Notifier <- "output.logs.clean"

			case <-requestCtx.Done():
				return
			}
		}
	}
}
