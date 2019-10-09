package main

import (
	"fmt"
	"log"
)

func beeper() error {
	fmt.Println("beep-beep")
	return nil
}

func pinger() error {
	fmt.Println("ping-ping")
	return nil
}

func jobExecutor(jobs []func() error) error {
	for _, fn := range jobs {
		err := fn()
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	var myJobs []func() error
	myJobs = append(myJobs, beeper)
	myJobs = append(myJobs, pinger)

	err := jobExecutor(myJobs)
	if err != nil {
		log.Fatalln(err)
	}
}
