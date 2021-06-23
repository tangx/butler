package butler

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"
)

func Test_Butler(t *testing.T) {
	// b := Defualt()

	b := &Butler{}
	b.Init(WithJobs(20), WithWorkers(5))

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
		fmt.Printf("job %d: sleep %d \n", rand.Int(), t)
	}

}
