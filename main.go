package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"net/http"
	"os"
	"strings"
)

func main() {
	app := cli.NewApp()

	app.Name = "kss"
	app.Usage = "A command line tool to interact with kss"
	app.Version = "1.0"

	app.Commands = []cli.Command{
		{
			Name:  "build",
			Usage: "Builds styleguide for given config",
			Action: func(c *cli.Context) {
				if len(c.Args()) < 1 {
					fmt.Println("please provide a config file")
					return
				}

				config, err := NewConfig(c.Args()[0])
				if err != nil {
					fmt.Println(err)
					return
				}

				guide, err := NewGuide(config)
				if err != nil {
					fmt.Println(err)
					return
				}

				err = Build(config, guide)
				if err != nil {
					fmt.Println(err)
					return
				}
			},
		},
		{
			Name:  "create",
			Usage: "Creates a new project in the given output directory",
			Action: func(c *cli.Context) {
				if len(c.Args()) < 1 {
					fmt.Println("please provide a output directory")
					return
				}

				err := Create(c.Args()[0])
				if err != nil {
					fmt.Println(err)
					return
				}
			},
		},
		{
			Name:  "serve",
			Usage: "Serves styleguide for given config",
			Action: func(c *cli.Context) {
				if len(c.Args()) < 1 {
					fmt.Println("please provide a config file")
					return
				}

				config, err := NewConfig(c.Args()[0])
				if err != nil {
					fmt.Println(err)
					return
				}

				guide, err := NewGuide(config)
				if err != nil {
					fmt.Println(err)
					return
				}

				addr := "127.0.0.1"
				port := "8080"

				if len(c.Args()) > 1 {
					bits := strings.Split(c.Args()[1], ":")
					if len(bits) > 1 {
						addr = bits[0]
						port = bits[1]
					} else {
						port = bits[0]
					}
				}

				address := fmt.Sprintf("%s:%s", addr, port)

				fmt.Printf("Starting server at %s\n", address)

				http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
					err := Serve(guide, w, r)
					if err != nil {
						fmt.Println(err)
					}
				})

				err = http.ListenAndServe(address, nil)
				if err != nil {
					fmt.Println(err)
				}
			},
		},
	}

	app.Run(os.Args)
}
