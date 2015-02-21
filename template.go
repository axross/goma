package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/kyokomi/goma/debuglog"
)

// DaoTemplateData dao data.
type DaoTemplateData struct {
	Name       string
	MemberName string
	EntityName string
	Table      TableTemplateData
	Imports    []string
}

// TableTemplateData table data.
type TableTemplateData struct {
	Name      string
	TitleName string
	Columns   []ColumnTemplateData
}

// ColumnTemplateData column data.
type ColumnTemplateData struct {
	Name         string
	TitleName    string
	TypeName     string
	IsPrimaryKey bool
	Sample       string
}

// HelperTemplateData helper data.
type HelperTemplateData struct {
	PkgName      string
	DriverImport string
	Options      map[string]interface{}
}

func (d HelperTemplateData) execHelperTemplate(rootDir string) error {
	var buf bytes.Buffer
	if err := HelperTemplate(&buf, d); err != nil {
		return err
	}
	return formatFileWrite(rootDir, "helper_gen.go", buf.Bytes())
}

func (d DaoTemplateData) execDaoTemplate(daoRootDir string) error {
	var buf bytes.Buffer
	if err := DaoTemplate(&buf, d); err != nil {
		return err
	}
	return formatFileWrite(daoRootDir, d.Table.Name+"_gen.go", buf.Bytes())
}

func (t TableTemplateData) execTableTemplate(sqlRootDir string) error {
	sqlTableDir := filepath.Join(sqlRootDir, t.Name)
	if err := os.MkdirAll(sqlTableDir, 0755); err != nil {
		return err
	}

	var err error
	var buf bytes.Buffer
	fileWriteFunc := func(templateFunc func(io.Writer, TableTemplateData) error, fileName string) {
		if err != nil {
			return
		}

		buf.Reset()
		err = templateFunc(&buf, t)
		if err != nil {
			return
		}

		filePath := filepath.Join(sqlTableDir, fileName)
		err = ioutil.WriteFile(filePath, buf.Bytes(), 0644)
		if err != nil {
			err = fmt.Errorf("file write error: %s \n%s", err, filePath)
			return
		}

		debuglog.Println("generate file:", filePath)
	}

	fileWriteFunc(SelectAllTemplate, "selectAll.sql")
	fileWriteFunc(SelectByIDTemplate, "selectByID.sql")
	fileWriteFunc(InsertTemplate, "insert.sql")
	fileWriteFunc(UpdateTemplate, "update.sql")
	fileWriteFunc(DeleteTemplate, "delete.sql")

	return err
}

func formatFileWrite(path, fileName string, data []byte) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	bts, err := format.Source(data)
	if err != nil {
		return fmt.Errorf("go format error: %s \n%s", err, string(data))
	}

	filePath := filepath.Join(path, fileName)
	if err := ioutil.WriteFile(filePath, bts, 0644); err != nil {
		return fmt.Errorf("file write error: %s \n%s", err, filePath)
	}

	debuglog.Println("generate file:", filePath)
	return nil
}
