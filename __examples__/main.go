package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/tangx/butler"
)

func main() {

	b := butler.Default()
	// or
	// b := &butler.Butler{}
	// b.WithOptions(butler.WithJobs(10), butler.WithWorkers(5))

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	b.WithOptions(butler.WithContext(ctx))

	b.Init()

	go func() {
		for i := 0; i < 30; i++ {
			fn := newJob()
			b.AddJobs(fn)
			time.Sleep(time.Second * 1)
		}
	}()

	b.Work()

}

func newJob() func() {
	log.Println("create a new job")
	return func() {
		rand.Seed(time.Now().UnixNano())
		t := rand.Intn(5)
		time.Sleep(time.Duration(t) * time.Second)
		jobid := rand.Int()
		fmt.Printf("job %d: sleep %d \n", jobid, t)
		if t%2 == 0 {
			log.Panic(jobid)
		}
	}

}
