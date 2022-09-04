package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Message struct {
	Content string
}

func producer(ch chan<- Message, countCh chan int, ctx context.Context) {
	id := rand.Int()
	for {
		select {
		case <-ctx.Done():
			break
		default:
			Counter := <-countCh
			Counter += 1
			countCh <- Counter
			ch <- Message{
				Content: fmt.Sprintf("Producer_%d sending message_%d", id, Counter),
			}
			time.Sleep(1000 * time.Millisecond)
		}
	}
}

func consumer(ch <-chan Message) {
	id := rand.Int()
	for {
		d := <-ch
		fmt.Printf("Consumer_%d get message: %s\n", id, d.Content)
	}
}

func main() {
	var ctx, cancel = context.WithCancel(context.Background())
	dataChannel := make(chan Message, 1)
	countChannel := make(chan int, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	countChannel <- 0
	go producer(dataChannel, countChannel, ctx)
	go producer(dataChannel, countChannel, ctx)
	go consumer(dataChannel)
	go consumer(dataChannel)
	time.Sleep(10 * time.Second)
	fmt.Println("Canceling producer..")
	cancel()
	wg.Done()
	wg.Wait()
}
