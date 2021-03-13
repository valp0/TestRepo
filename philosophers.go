// This program simulates the dining philosophers problem, a classical computer science
// problem (https://en.wikipedia.org/wiki/Dining_philosophers_problem).

package main

import (
	"fmt"
	"sync"
	"time"
)

type ChopS struct{ sync.Mutex }

type Philo struct {
	leftCS, rightCS *ChopS
	count, philoNum int
}

const eatTimes int = 3 // Controls the amount of times you want each philosopher to eat.

func main() {
	var wg sync.WaitGroup
	c1 := make(chan bool, 2)
	c2 := make(chan bool, 2)

	// Creating 5 chop sticks,
	CSticks := make([]*ChopS, 5) // Creating an array of 5 chopsticks
	for i := 0; i < 5; i++ {
		CSticks[i] = new(ChopS) // Filling the array with mutex type ChopS
	}
	philos := make([]*Philo, 5) // Creating an array of 5 philosophers
	for i := 0; i < 5; i++ {
		philos[i] = &Philo{CSticks[i], CSticks[(i+1)%5], 0, i + 1} // Filling philosohers array with type Philo.
	} /*-------------*/ // This chunk of code \_______/ creates an abstract "circular" table.

	wg.Add(eatTimes * 5)
	for i := 0; i < 5; i++ {
		go philos[i].eat(c1, c2, &wg) // Calling each philosopher to eat (<eatTimes> is in the eat() method).
	}
	go host(c1, c2)
	wg.Wait()
}

func (p Philo) eat(c1 chan bool, c2 chan bool, wg *sync.WaitGroup) {
	for {
		if p.count < eatTimes {
			c1 <- true // If c1 is full, it waits until host function cleans it up.
			p.leftCS.Lock()
			p.rightCS.Lock()

			// Eating and increasing count of eat times
			fmt.Printf("Starting to eat Philo #%d (%d%s time)\n", p.philoNum, p.count+1, firSecThi(p.count+1))
			p.count++
			time.Sleep(1 * time.Nanosecond) // This helps to notice concurrency in the output, but isn't critical.
			fmt.Printf("Finished eating Philo #%d\n", p.philoNum)

			p.rightCS.Unlock()
			p.leftCS.Unlock()
			wg.Done()
			c2 <- true // Sending c2, so host function knows one philosopher is done eating.
		}
	}
}

// The host function allows only two philosophers to eat simultaneously.
func host(c1 chan bool, c2 chan bool) {
	for {
		if len(c2) >= 1 {
			<-c1 // Will only free c1 (permission to eat) if c2 (finished eating) is one or more in length,
			<-c2 // c2 is received just to make sure its capacity doesn't get exceeded and stops the code.
		}
	}
}

// Just a function to prettify the output
func firSecThi(x int) string {
	switch {
	case x == 1:
		return "st"
	case x == 2:
		return "nd"
	case x == 3:
		return "rd"
	case x > 20:
		switch {
		case x%10 == 1:
			return "st"
		case x%10 == 2:
			return "nd"
		case x%10 == 3:
			return "rd"
		}
	}
	return "th"
}
