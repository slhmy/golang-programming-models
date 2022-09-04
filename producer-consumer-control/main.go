package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type Message struct {
	Content string
}

func producer(dataCh chan<- Message, countCh chan int, ctx context.Context) {
	// Create a unique id for the producer
	id := rand.Int()
	for {
		select {
		// Handle <-ctx.Done()
		case <-ctx.Done():
			fmt.Printf("Producer_%d quit\n", id)
			return
		default:
			// Get message count
			Count := <-countCh
			Count += 1
			countCh <- Count

			// Produce
			dataCh <- Message{
				Content: fmt.Sprintf("Producer_%d sending message_%d", id, Count),
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func consumer(dataCh <-chan Message) {
	// Create a unique id for the consumer
	id := rand.Int()
	for {
		// Consume
		data := <-dataCh
		fmt.Printf("Consumer_%d get message: %s\n", id, data.Content)
	}
}

func main() {
	// Prepare
	var ctx, cancel = context.WithCancel(context.Background())
	dataChannel := make(chan Message)
	countChannel := make(chan int, 1)
	countChannel <- 0

	// Start goroutines
	go producer(dataChannel, countChannel, ctx)
	go producer(dataChannel, countChannel, ctx)
	go consumer(dataChannel)
	go consumer(dataChannel)

	time.Sleep(5 * time.Second)

	fmt.Println("Canceling producer..")
	cancel()
	time.Sleep(1 * time.Second)
}
