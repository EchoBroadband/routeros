package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/EchoBroadband/routeros"
)

var (
	command  = flag.String("command", "/system/resource/print", "RouterOS command")
	address  = flag.String("address", "127.0.0.1:8728", "RouterOS address and port")
	username = flag.String("username", "admin", "User name")
	password = flag.String("password", "admin", "Password")
	async    = flag.Bool("async", false, "Use async code")
	useTLS   = flag.Bool("tls", false, "Use TLS")
	verbose  = flag.Bool("v", false, "Use verbose mode")
)

func dial() (*routeros.Client, error) {
	if *useTLS {
		return routeros.DialTLS(*address, *username, *password, nil)
	}
	return routeros.Dial(*address, *username, *password)
}

func main() {
	flag.Parse()

	c, err := dial()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	if *async {
		c.Async()
	}

	cmdArray := []string{}

	tmpCmd := *command
	quoteIndex := strings.Index(tmpCmd, `"`)

	for quoteIndex != -1 {
		// Extract up to the first instance of quote in string
		ss := tmpCmd[:quoteIndex]
		// Split by spaces to get the command arguments
		split := strings.Split(ss, " ")

		// Get everything after the first quote
		ss = tmpCmd[quoteIndex+1:]

		// Get the next instance of quote
		quoteIndex = strings.Index(ss, `"`)

		// Append everything up to the second quotation
		split[len(split)-1] += ss[:quoteIndex]

		// Append to cmdArray
		cmdArray = append(cmdArray, split...)

		// tmCmd is everything after the second quotation mark
		tmpCmd = ss[quoteIndex+1:]

		// Get the next quotation mark index
		quoteIndex = strings.Index(tmpCmd, `"`)
	}

	cmdArray = append(cmdArray, strings.Split(strings.Trim(tmpCmd, " "), " ")...)

	if *verbose {
		for _, cmd := range cmdArray {
			fmt.Println(cmd)
		}
	}

	r, err := c.RunArgs(cmdArray)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(r)
}
