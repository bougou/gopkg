package main

import (
	"fmt"
	"time"

	"github.com/signcl/grsync"
)

func main() {
	task := grsync.NewTask(
		"datasets",
		"tests",
		grsync.RsyncOptions{
			Info: "progress2",
			// Progress: true,
			// Verbose:  true,
			// Quiet:    false,
		},
	)

	go func() {
		for {
			state := task.State()
			fmt.Printf(
				"progress: %.2f / rem. %d / tot. %d / sp. %s \n",
				state.Progress,
				state.Remain,
				state.Total,
				state.Speed,
			)
			<-time.After(time.Second)
		}
	}()

	if err := task.Run(); err != nil {
		panic(err)
	}

	fmt.Println("well done")
	fmt.Println(task.Log())
}
