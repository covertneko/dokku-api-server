package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/nikelmwann/dokku-api/dokku"
)

const (
	DEFAULT_SOCKET_DIR = "/tmp/dokku-api"

	// Socket locations relative to the socket directory
	CACHE_SOCKET = "cache.sock"
	API_SOCKET   = "api/api.sock"
)

func getCacheSocket(c *cli.Context) string {
	dir := c.String("socket-dir")
	if len(dir) == 0 {
		dir = DEFAULT_SOCKET_DIR
	}

	return path.Join(dir, CACHE_SOCKET)
}

func getApiSocket(c *cli.Context) string {
	dir := c.String("socket-dir")
	if len(dir) == 0 {
		dir = DEFAULT_SOCKET_DIR
	}

	return path.Join(dir, API_SOCKET)
}

func cache(action string, c *cli.Context) {
	socket := getCacheSocket(c)

	conn, err := net.Dial("unix", socket)
	if err != nil {
		err = fmt.Errorf("Error connecting to cache server - caching may not be enabled: %s", err)
		panic(err)
	}
	defer conn.Close()

	go func() {
		command := action + " " + strings.Join(c.Args(), " ") + "\n"

		_, err = conn.Write([]byte(command))
		if err != nil {
			panic(err)
		}
	}()

	// Read one line for response
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		panic(err)
	}

	res := line[:len(line)-1]

	// Panic with error if command isn't successful
	if res != CACHE_COMMAND_RESPONSE_SUCCESS {
		panic(res)
	}

	// Otherwise print success
	fmt.Println(res)
}

func serve(c *cli.Context) {
	// Ensure api socket directory exists; if not, create it
	apiSocket := getApiSocket(c)
	if _, err := os.Stat(apiSocket); os.IsNotExist(err) {
		os.MkdirAll(path.Dir(apiSocket), 0777)
	}

	d := dokku.New()

	// If caching is enabled, create cache socket if necessary and launch cache
	// server in the background
	if !c.BoolT("enable-cache") {
		cacheSocket := getCacheSocket(c)
		if _, err := os.Stat(cacheSocket); os.IsNotExist(err) {
			os.MkdirAll(path.Dir(cacheSocket), 0777)
		}

		go func() {
			err := listenCache(d, cacheSocket)
			if err != nil {
				panic(err)
			}
		}()
	}

	// Start api server
	err := listenAPI(d, apiSocket)
	if err != nil {
		panic(err)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "dokku-api"
	app.Usage = "Serve a REST API to interact with the Dokku instance running on this host."

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "socket-dir, d",
			Value: DEFAULT_SOCKET_DIR,
			Usage: "Specify the directory in which to create the API socket.",
		},
		cli.BoolTFlag{
			Name:  "enable-cache, c",
			Usage: "Specify whether to cache API resources (requires external cache invalidation via the cache:purge command)",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "serve",
			Usage: "Serve the API on a socket in the directory specified by socket-dir",
			Action: func(c *cli.Context) {
				serve(c)
			},
		},
		{
			Name:  "cache:purge",
			Usage: "Purge a particular cache entry for a resource",
			Action: func(c *cli.Context) {
				cache("purge", c)
			},
			ArgsUsage: "<type> <id> - where <type> is one of 'app' or 'container'",
		},
		{
			Name:  "cache:fetch",
			Usage: "Find a resource and store it in the cache, replacing any previous entry",
			Action: func(c *cli.Context) {
				cache("fetch", c)
			},
			ArgsUsage: "<type> <id> - where <type> is one of 'app' or 'container'",
		},
	}

	// Handle panics in commands
	// TODO: do something smarter than this for errors
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error:", r)
			os.Exit(1)
		}
	}()

	app.Run(os.Args)
}
