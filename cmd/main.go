package main

import (
	"streamlabs/internal/streamlabs"

	"github.com/zenthangplus/goccm"
)

func main() {
	c := goccm.New(1)
	for {
		c.Wait()
		go func() {
			defer c.Done()
			streamlabs.Create()
		}()
	}
	c.WaitAllDone()

}
