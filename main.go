package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/isbadawi/octranspo/api"
)

func stopCommandHandler(stop string) {
	routes, err := api.GetNextTripsForStop(stop)
	if err != nil {
		fmt.Printf("error: %v\n", err.Error())
		return
	}

	for _, route := range routes {
		fmt.Printf("%v %v (%v): ", route.Number, route.Heading, route.Direction)
		if len(route.Trips) == 0 {
			fmt.Printf("no more trips today.")
		}
		fmt.Println()
		for _, trip := range route.Trips {
			fmt.Printf("  %v (%v)\n", trip.StartTime, trip.Destination)
		}
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "octranspo"
	app.Usage = "OC Transpo CLI"
	app.Version = "1.0.0"
	app.Commands = []cli.Command{
		{
			Name:  "stop",
			Usage: "get upcoming trips for a given stop number",
			Action: func(c *cli.Context) {
				if len(c.Args()) == 0 {
					fmt.Println("usage: octranspo stop <number>")
					return
				}

				stopCommandHandler(c.Args()[0])

			},
		},
	}
	app.Run(os.Args)
}
