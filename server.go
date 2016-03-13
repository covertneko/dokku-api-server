package main

import (
	"os"
	"bufio"
	"strings"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/nikelmwann/dokku-api/dokku"
	r "github.com/nikelmwann/dokku-api/resources"
)

const (
	PLUGIN_COMMAND_RESPONSE_SUCCESS = "ok"
)

func handlePluginCommand(d *dokku.Dokku, conn net.Conn) {
	defer conn.Close()

	for {
		_command, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			conn.Write([]byte("Error reading command\n"))
			continue
		}

		// Trim trailing newline from command
		command := _command[:len(_command) - 1]

		cols := strings.Split(command, " ")

		if len(cols) < 3 {
			conn.Write([]byte(fmt.Sprintf("Unrecognized command %q", command) + "\n"))
			continue
		}

		action := cols[0]
		_type := cols[1]
		id := cols[2]

		switch action {
		case "purge":
			switch _type {
			case "app":
				d.Apps.Invalidate(id)
				break
			case "container":
				d.Containers.Invalidate(id)
				break
			default:
				conn.Write([]byte(fmt.Sprintf("Unrecognized type %q", _type) + "\n"))
				continue
			}
			conn.Write([]byte(PLUGIN_COMMAND_RESPONSE_SUCCESS + "\n"))
			break
		case "update":
			switch _type {
			case "app":
				d.Apps.Invalidate(id)
				d.Apps.Find(id)
				break
			case "container":
				d.Containers.Invalidate(id)
				d.Containers.Find(id)
				break
			default:
				conn.Write([]byte(fmt.Sprintf("Unrecognized type %q", _type) + "\n"))
				continue
			}
			conn.Write([]byte(PLUGIN_COMMAND_RESPONSE_SUCCESS + "\n"))
			break
		default:
			conn.Write([]byte(fmt.Sprintf("Unrecognized action %q", action) + "\n"))
		}
	}
}

func listenPlugin(d *dokku.Dokku, socketPath string) error {
	ln, err := ListenSocket(socketPath, 0666)
	if err != nil {
		return err
	}
	defer ln.Close()
	defer os.Remove(socketPath)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err.Error())
			continue;
		}

		go handlePluginCommand(d, conn)
	}
}

func listenAPI(d *dokku.Dokku, socketPath string) error {
	e := echo.New()
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	e.Get("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "TODO: add version and stuff")
	})
	e.Get("/apps", r.HandlerFor(r.Apps{}, d))
	e.Get("/apps/:name", r.HandlerFor(r.App{}, d))

	err := (*Echo)(e).RunDomainSocket(socketPath, os.FileMode(0666))

	if err != nil {
		return err
	}

	return nil
}
