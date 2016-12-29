package main

import (
	"log"
)

func leaseOnCh(leaseCh chan struct{}, releaseCh chan struct{}) {
	avail := 100
	for {
		outCh := leaseCh
		if avail <= 0 {
			outCh = nil
		}
		select {
		case _, ok := <-releaseCh:
			if !ok {
				return
			}
			log.Println("Release", avail)
			avail++
		case outCh <- struct{}{}:
			log.Println("Lease", avail)
			avail--
		}
	}
}
