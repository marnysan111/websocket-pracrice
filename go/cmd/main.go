package main

import "github.com/marnysan111/websocket-pracrice/internal/router"

func main() {
	r := router.SetupRouter()
	r.Run()
}
