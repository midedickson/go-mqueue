package main

import (
	"fmt"
	"sync"
	"time"
)

// Message represents a simple message structure
type Message struct {
	ID      int
	Content string
}

// MessageQueue represents a simple in-memory message queue
type MessageQueue struct {
	messages []Message
	lock     sync.Mutex
}

// Produce adds a message to the message queue
func (mq *MessageQueue) Produce(wg *sync.WaitGroup, message Message) {
	defer wg.Done()
	mq.lock.Lock()
	defer mq.lock.Unlock()
	mq.messages = append(mq.messages, message)
	fmt.Printf("Produced: %+v\n", message)
}

// Consume removes and returns a message from the message queue
func (mq *MessageQueue) Consume(wg *sync.WaitGroup, result chan<- Message) {
	defer wg.Done()
	mq.lock.Lock()
	defer mq.lock.Unlock()

	if len(mq.messages) == 0 {
		fmt.Println("Queue is empty")
		return
	}

	message := mq.messages[0]
	mq.messages = mq.messages[1:]
	fmt.Printf("Consumed: %+v\n", message)
	result <- message
}

func main() {
	var wg sync.WaitGroup
	messageQueue := MessageQueue{}
	result := make(chan Message, 5)

	// Start the producer
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= 5; i++ {
			message := Message{ID: i, Content: fmt.Sprintf("Message %d", i)}
			wg.Add(1)
			go messageQueue.Produce(&wg, message)
			time.Sleep(time.Second)
		}
	}()

	// Start the consumer
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= 5; i++ {
			wg.Add(1)
			go messageQueue.Consume(&wg, result)
			time.Sleep(time.Second)
		}
	}()

	// Wait for all goroutines to finish
	wg.Wait()

	// Close the result channel after all consumers are done
	close(result)
}
