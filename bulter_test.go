package bulter

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"
)

func Test_Bulter(t *testing.T) {
	bulter := New(WithJobs(20), WithWorkers(5))

	go func() {
		for i := 0; i < 30; i++ {
			fn := newJob()
			bulter.AddJobs(fn)
			time.Sleep(time.Second * 1)
		}
	}()

	bulter.Work()

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
