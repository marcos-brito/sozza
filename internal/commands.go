package internal

import (
	"os"
	"strconv"

	"github.com/marcos-brito/sozza/internal/connector"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func Create(ctx *cli.Context) error {
	db := connector.PickConnector(ctx.String("dbname")).Connect(ctx.String("url"))
	schema, err := os.ReadFile(ctx.String("schema"))

	if err != nil {
		log.Fatalf("Could not read the schema file: %s", err)
	}

	_, err = db.Exec(string(schema))

	if err != nil {
		log.Fatalf("Something went wrong executing the schema: %s", err)
	}

	return nil
}

func Insert(ctx *cli.Context) error {
	return nil
}
