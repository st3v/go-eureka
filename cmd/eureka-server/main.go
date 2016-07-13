package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/st3v/jolt/test"
)

const warning = `
    !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
    !!!                                            !!!
    !!! THIS EUREKA SERVER IS FOR TESTING PURPOSES !!!
    !!!    ONLY. DO NOT USE THIS IN PRODUCTION.    !!!
    !!!                                            !!!
    !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
`

var (
	host  string
	port  int
	debug bool
)

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	flag.StringVar(&host, "host", "0.0.0.0", "Host")
	flag.IntVar(&port, "port", 8080, "Port")
	flag.BoolVar(&debug, "debug", false, "Turn on debug logging")
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", host, port)
	server := test.NewServer(addr, debug)

	log.Println(warning)

	log.Printf("Listening on %s...\n", addr)
	log.Fatal(server.ListenAndServe())
}
