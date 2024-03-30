package main

import (
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/marcos-brito/sozza/internal"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "dbname",
				Aliases: []string{"d"},
				Value:   "oracle",
				Usage:   "The database to be used",
			},
			&cli.StringFlag{
				Name:     "url",
				Aliases:  []string{"u"},
				Usage:    "The connection url",
				Required: true,
			},
		},
		Commands: []*cli.Command{
			{
				Action: internal.Create,
				Name:   "create",
				Usage:  "Create the tables from a sql schema",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "schema",
						Aliases:  []string{"s"},
						Usage:    "The schema path",
						Required: true,
					},
				},
			},
			{
				Action: internal.Insert,
				Name:   "insert",
				Usage:  "Insert the content from a csv file in the database",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "mapping",
						Aliases:  []string{"m"},
						Usage:    "A .yml file with the mapping",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "number-of-lines",
						Aliases:  []string{"n"},
						Usage:    "The number of lines to be inserted",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "csv",
						Aliases:  []string{"c"},
						Usage:    "The path to the csv file",
						Required: true,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
