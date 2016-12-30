package main

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
				close(leaseCh)
				return
			}
			avail++
		case outCh <- struct{}{}:
			avail--
		}
	}
}
