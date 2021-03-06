package ts

import (
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx"

	"github.com/timescale/outflux/internal/idrf"
)

const (
	createTableQueryTemplate                = "CREATE TABLE \"%s\"(%s)"
	createTableWithSchemaQueryTemplate      = "CREATE TABLE \"%s\".\"%s\"(%s)"
	columnDefTemplate                       = "\"%s\" %s"
	createHypertableQueryTemplate           = "SELECT create_hypertable('\"%s\"', '%s');"
	createHypertableWithSchemaQueryTemplate = "SELECT create_hypertable('\"%s\".\"%s\"', '%s');"
	createTimescaleExtensionQuery           = "CREATE EXTENSION IF NOT EXISTS timescaledb"
)

type tableCreator interface {
	CreateTable(dbConn *pgx.Conn, info *idrf.DataSet) error
	CreateHypertable(dbConn *pgx.Conn, info *idrf.DataSet) error
	CreateTimescaleExtension(dbConn *pgx.Conn) error
}

func newTableCreator() tableCreator {
	return &defaultTableCreator{}
}

type defaultTableCreator struct{}

func (d *defaultTableCreator) CreateTable(dbConn *pgx.Conn, info *idrf.DataSet) error {
	query := dataSetToSQLTableDef(info)
	log.Printf("Creating table with:\n %s", query)

	_, err := dbConn.Exec(query)
	if err != nil {
		return err
	}

	log.Printf("Preparing TimescaleDB extension:\n%s", createTimescaleExtensionQuery)
	_, err = dbConn.Exec(createTimescaleExtensionQuery)
	if err != nil {
		return err
	}

	schema, table := info.SchemaAndTable()
	var hypertableQuery string
	if schema != "" {
		hypertableQuery = fmt.Sprintf(createHypertableWithSchemaQueryTemplate, schema, table, info.TimeColumn)
	} else {
		hypertableQuery = fmt.Sprintf(createHypertableQueryTemplate, info.DataSetName, info.TimeColumn)
	}

	log.Printf("Creating hypertable with: %s", hypertableQuery)
	_, err = dbConn.Exec(hypertableQuery)
	if err != nil {
		return err
	}

	return nil
}

func (d *defaultTableCreator) CreateHypertable(dbConn *pgx.Conn, info *idrf.DataSet) error {
	hypertableQuery := fmt.Sprintf(createHypertableQueryTemplate, info.DataSetName, info.TimeColumn)
	log.Printf("Creating hypertable with: %s", hypertableQuery)
	_, err := dbConn.Exec(hypertableQuery)
	return err
}

func (d *defaultTableCreator) CreateTimescaleExtension(dbConn *pgx.Conn) error {
	log.Printf("Preparing TimescaleDB extension:\n%s", createTimescaleExtensionQuery)
	_, err := dbConn.Exec(createTimescaleExtensionQuery)
	return err
}

func dataSetToSQLTableDef(dataSet *idrf.DataSet) string {
	columnDefinitions := make([]string, len(dataSet.Columns))
	for i, column := range dataSet.Columns {
		dataType := idrfToPgType(column.DataType)
		columnDefinitions[i] = fmt.Sprintf(columnDefTemplate, column.Name, dataType)
	}

	columnsString := strings.Join(columnDefinitions, ", ")

	schema, table := dataSet.SchemaAndTable()
	if schema != "" {
		return fmt.Sprintf(createTableWithSchemaQueryTemplate, schema, table, columnsString)
	}
	return fmt.Sprintf(createTableQueryTemplate, table, columnsString)
}
