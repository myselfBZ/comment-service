package main

import "log"

func main() {
	a := NewApp()
    log.Print("server started")
	a.run(":8080")
}
