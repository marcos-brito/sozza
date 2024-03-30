package connector

import "database/sql"

type Connector interface {
	Connect(url string) *sql.DB
}

func PickConnector(databaseName string) Connector {
	switch databaseName {
	case "postgres":
		return &Postgres{}
	case "sqlite":
		return &Sqlite{}
	// case "mysql":
	// 	return &Mysql{}
	// case "oracle":
	// 	return &Oracle{}
	default:
		// should be oracle by default
		return &Postgres{}
	}
}
