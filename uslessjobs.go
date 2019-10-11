package main

import (
	"errors"
	"time"
)

func job1s() error {
	time.Sleep(time.Second * 1)
	return nil
}

// func job3s() error {
// 	time.Sleep(time.Second * 3)
// 	return nil
// }

func job2sErr() error {
	time.Sleep(time.Second * 2)
	return errors.New("job2sErr error")
}
