package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"github.com/mitchellh/colorstring"
)
type CheckExplain struct {
	TableName string
	QueryName string
	MySQLExplain
}
type MySQLExplain struct {
	ID           int
	SelectType   sql.NullString
	Table        sql.NullString
	Type         sql.NullString
	PossibleKeys sql.NullString
	Key          sql.NullString
	KeyLen       sql.NullInt64
	Ref          sql.NullString
	Rows         sql.NullInt64
	Extra        sql.NullString
}

/*
                            QUERY PLAN
-------------------------------------------------------------------
 Seq Scan on goma_string_types  (cost=0.00..1.00 rows=1 width=592)
   Filter: (id = 1)
(2 rows)
 */
var postgresExplainTemplate = "[red]%s : %s [default]=> %s"
var mysqlExplainTemplate = "[red]%s : %s [default]=> |  %d | %s      | %s  | %s | %s          | %s | %d    | %s | %d | %s |"

func main() {
	log.SetFlags(log.Llongfile)

	// TODO: driverを見て、mysqlとpostgresで分岐

	// TODO: cli引数
	db, err := sql.Open("mysql", "root:@/goma_test")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	// TODO: cli引数のsqlディレクトリ下を見る
	// dir1名: tableName
	// dir2名: queryName
	// file内容: query
	query := "select * from goma_string_types"
	err = PrintExplain(db, "goma_string_types", "selectAll", query)
	if err != nil {
		log.Fatalln(err)
	}
}

func PrintExplain(db *sql.DB, tableName, queryName, query string) error {
	e := CheckExplain{}
	e.TableName = tableName
	e.QueryName = queryName
	row := db.QueryRow(fmt.Sprintf("EXPLAIN %s", query))
	err := row.Scan(
		&e.ID,
		&e.SelectType,
		&e.Table,
		&e.Type,
		&e.PossibleKeys,
		&e.Key,
		&e.KeyLen,
		&e.Ref,
		&e.Rows,
		&e.Extra,
	)
	if err != nil {
		return err
	}

	colorstring.Println(fmt.Sprintf(mysqlExplainTemplate,
		e.TableName,
		e.QueryName,
		e.ID,
		e.SelectType.String,
		e.Table.String,
		e.Type.String,
		e.PossibleKeys.String,
		e.Key.String,
		e.KeyLen.Int64,
		e.Ref.String,
		e.Rows.Int64,
		e.Extra.String,
	))

	return nil
}

