package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"path/filepath"

	"strings"

	"github.com/codegangsta/cli"
	"github.com/jmoiron/sqlx"
	"github.com/kyokomi/goma"
)

func genAction(c *cli.Context) {
	generate(c.GlobalString("pkg"), scanGenFlags(c))
}

func initConfigAction(c *cli.Context) {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	opt := goma.Options{}
	opt.CurrentDir = currentDir
	data, err := json.MarshalIndent(opt, "", "    ")
	if err != nil {
		log.Fatalln(err)
	}
	if err := ioutil.WriteFile(opt.ConfigPath(), data, 0644); err != nil {
		log.Fatalln(err)
	}
}

func genConfigAction(c *cli.Context) {
	opt, err := goma.NewOptions(c.String("path"))
	if err != nil {
		log.Fatalln(err)
	}
	generate(c.GlobalString("pkg"), opt)
}

func scanGenFlags(c *cli.Context) goma.Options {
	opt := goma.Options{}
	opt.Driver = c.String("driver")
	opt.UserName = c.String("user")
	opt.PassWord = c.String("password")
	opt.Host = c.String("host")
	opt.Port = c.Int("port")
	opt.DBName = c.String("db")
	opt.Location = c.String("location")
	opt.SSLMode = c.String("ssl")
	opt.SQLRootDir = c.String("sql")
	opt.DaoRootDir = c.String("dao")
	opt.EntityRootDir = c.String("entity")
	opt.IsConfig = c.Bool("config")
	opt.Debug = c.GlobalBool("debug")
	return opt
}

func explainAction(c *cli.Context) {
	opt := scanGenFlags(c)

	db, err := goma.OpenOptions(opt)
	if err != nil {
		log.Fatalf("goma.Open error : %s", err)
	}
	defer db.Close()

	infos, err := ioutil.ReadDir(opt.SQLRootDirPath())
	if err != nil {
		log.Fatalf("sql dir ReadDir paht = %s error : %s", opt.SQLRootDirPath(), err)
	}

	// TODO: あとで
	queryArgs := map[string]interface{}{
		"id":                int64(111111),
		"tinyint_columns":   int(8),
		"bool_columns":      false,
		"smallint_columns":  int(123),
		"mediumint_columns": int(256),
		"int_columns":       int(11111111),
		"integer_columns":   int(22222222),
		"decimal_columns":   int64(111111),
		"numeric_columns":   "1234567890",
		"float_columns":     "1234567890",
		"double_columns":    float32(1.234),
	}

	for _, info := range infos {
		if !info.IsDir() {
			continue
		}

		tableName := info.Name()

		tableDirPath := filepath.Join(opt.SQLRootDirPath(), info.Name())
		tableFileInfos, err := ioutil.ReadDir(tableDirPath)
		if err != nil {
			log.Fatalf("table file ReadDir paht = %s error : %s", tableDirPath, err)
		}

		for _, fInfo := range tableFileInfos {
			if fInfo.IsDir() {
				continue
			}

			queryName := fInfo.Name()

			sqlFilePath := filepath.Join(tableDirPath, fInfo.Name())
			sqlFileData, err := ioutil.ReadFile(sqlFilePath)
			if err != nil {
				log.Fatalf("sql file ReadFile paht = %s error : %s", sqlFilePath, err)
			}

			sqlQuery, args, err := sqlx.Named(string(sqlFileData), queryArgs)

			if strings.HasPrefix(queryName, "delete") ||
				strings.HasPrefix(queryName, "insert") ||
				strings.HasPrefix(queryName, "update") {
				continue
			}

			if err := PrintExplain(db, tableName, queryName, sqlQuery, args...); err != nil {
				log.Printf("printExplain error %s", err)
			}
		}
	}
}
