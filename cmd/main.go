/*
Developer: Terminal

Discord : casanova4
Telegram : https://t.me/icebergs

If u need good proxies for cf hmu
*/

package main

import (
	"fmt"
	"streamlabsuwu/internal/streamlabs"

	"github.com/zenthangplus/goccm"
)

func main() {
	var err error
	var threads int
	fmt.Print("Enter the number of threads: ")
	_, err = fmt.Scanln(&threads)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	c := goccm.New(threads)
	for {
		c.Wait()
		go func() {
			defer c.Done()
			streamlabs.Create()
		}()
	}
	c.WaitAllDone()

}
