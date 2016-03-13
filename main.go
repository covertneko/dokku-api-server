package main

import (
	"os"
	"path"

	"github.com/codegangsta/cli"
	"github.com/nikelmwann/dokku-api/dokku"
)

const (
	DOKKU_API_DEFAULT_SOCKDIR = "/tmp/dokku-api"
)

func serve(sockDir string) {
	sockPath := path.Join(sockDir, "plugin.sock")
	apiSockPath := path.Join(sockDir, "api", "api.sock")

	// Ensure socket directories exist
	// If not, create them
	if _, err := os.Stat(sockPath); os.IsNotExist(err) {
		os.MkdirAll(path.Dir(sockPath), 0777)
	}

	if _, err := os.Stat(apiSockPath); os.IsNotExist(err) {
		os.MkdirAll(path.Dir(apiSockPath), 0777)
	}

	d := dokku.New()

	// Start plugin server
	go func() {
		err := listenPlugin(d, sockPath)
		if err != nil {
			panic(err)
		}
	}()

	// Start api server and block
	err := listenAPI(d, apiSockPath)
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
			Name: "socket-dir, d",
			Value: DOKKU_API_DEFAULT_SOCKDIR,
			Usage: "Specify the directory in which to create the API socket.",
		},
	}

	app.Commands = []cli.Command{
		{
			Name: "serve",
			Usage: "Serve the API on a socket in the directory specified by socket-dir",
			Action: func(c *cli.Context) {
				dir := c.String("socket-dir")
				if len(dir) == 0 {
					dir = DOKKU_API_DEFAULT_SOCKDIR
				}
				serve(dir)
			},
		},
	}

	app.Run(os.Args)
}
