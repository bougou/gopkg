package snippet

import (
	"fmt"
	"sync"
	"testing"
)

// wrong way
func test1(wg *sync.WaitGroup) {
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			fmt.Println(i) // always print 10
		}()
	}
}

// right way
func test2(wg *sync.WaitGroup) {

	for i := 0; i < 10; i++ {
		wg.Add(1)
		theI := i
		go func() {
			defer wg.Done()

			fmt.Println(theI)
		}()
	}
}

func Test_WaitGroup(t *testing.T) {
	var wg sync.WaitGroup

	test1(&wg)
	wg.Wait()

	test2(&wg)
	wg.Wait()

	fmt.Println("over")
}
