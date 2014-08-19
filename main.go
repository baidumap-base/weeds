package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	content, e := ioutil.ReadFile("/home/clark/桌面/map_android_crash.201408170000.map_android_crash")
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}

	buf := bytes.NewBuffer(content)
	var s string
	weeds := make([]*Weed, 0, 128)
	var w *Weed
	for {
		s, e = buf.ReadString(byte('\n'))
		if e != nil {
			break
		}

		w, e = parseWeed(s[:len(s)-1])
		if e != nil {
			continue
		}

		weeds = append(weeds, w)
	}

	for _, ele := range weeds {
		fmt.Println(ele.Sv)
	}
}
