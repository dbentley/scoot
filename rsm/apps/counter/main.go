package main

import (
	"log"
	"time"

	"github.com/scootdev/scoot/rsm"
	"github.com/scootdev/scoot/rsm/logs"
)

func main() {
	l := logs.NewMemoryTopicLog()
	mach := &QueueMachine{}

	manager := rsm.NewManagerImpl(l, mach)

	q, err := NewQueue(manager)
	if err != nil {
		log.Fatal(err)
	}

	go client(q, "client1", 1, 2)
	go client(q, "client2", 2, 1)
	go client(q, "client3", 3, 36)
	go client(q, "client4", 3, 36)
	go processor(q, "agent1")

	log.Println("Started")

	time.Sleep(365 * 24 * 60 * time.Minute)
}
