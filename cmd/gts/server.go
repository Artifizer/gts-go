/*
Copyright © 2025 Global Type System
Released under Apache License 2.0
*/

package main

import (
	"fmt"

	"github.com/GlobalTypeSystem/gts-go/server"
)

var cmdServer = &Command{
	UsageLine: "server [-host address] [-port number]",
	Short:     "start the GTS HTTP server",
	Long: `
Server starts the GTS HTTP server for REST API access.

The -host flag specifies the host address (default: 127.0.0.1).
The -port flag specifies the port number (default: 8000).

Example:

	gts -path ./examples server -host 127.0.0.1 -port 8000
	`,
}

var (
	serverHost string
	serverPort int
)

func init() {
	cmdServer.Run = runServer
	cmdServer.Flag.StringVar(&serverHost, "host", "127.0.0.1", "host address")
	cmdServer.Flag.IntVar(&serverPort, "port", 8000, "port number")
}

func runServer(cmd *Command, args []string) {
	store := newStore()

	fmt.Printf("starting server at http://%s:%d\n", serverHost, serverPort)
	if verbose == 0 {
		fmt.Println("use -v for verbose logging")
	}

	srv := server.NewServer(store, serverHost, serverPort, verbose)
	if err := srv.Start(); err != nil {
		fatalf("server failed: %v", err)
	}
}
