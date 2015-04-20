package main

import (
	"log"
	"time"

	_ "github.com/lib/pq"

	"database/sql"
	"reflect"
	"testing"

	"github.com/kyokomi/goma"
	"github.com/kyokomi/goma/test/simple/postgres/dao"
	"github.com/kyokomi/goma/test/simple/postgres/entity"
)

const testID = int64(1234567892)

func TestNumeric(t *testing.T) {
	db, err := goma.Open("config.json")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	id := testID

	// numeric
	d := dao.GomaNumericTypes(db)

	insertData := entity.GomaNumericTypesEntity{
		ID:              id,
		BoolColumns:     true,
		SmallintColumns: int(123),
		IntColumns:      int(11111111),
		IntegerColumns:  int(22222222),
		SerialColumns:   1234567890,
		DecimalColumns:  "1234567890",
		NumericColumns:  "1234567890",
		FloatColumns:    float64(1.234),
	}

	if _, err := d.Insert(insertData); err != nil {
		t.Errorf("ERROR: %s", err)
	}

	if e, err := d.SelectByID(id); err != nil {
		t.Errorf("ERROR: %s", err)
	} else if !reflect.DeepEqual(e, insertData) {
		t.Errorf("ERROR: %+v != %+v", e, insertData)
	}

	if _, err := d.Delete(id); err != nil {
		t.Errorf("ERROR: %s", err)
	}

	if _, err := d.SelectByID(id); err != sql.ErrNoRows {
		t.Errorf("ERROR: %s", "Deleteしたのにnilじゃない")
	}
}

func TestString(t *testing.T) {
	db, err := goma.Open("config.json")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	id := testID

	// string
	d := dao.GomaStringTypes(db)

	insertData := entity.GomaStringTypesEntity{
		ID:             id,
		TextColumns:    "あいうえおかきくけこ",
		CharColumns:    "a       ",
		VarcharColumns: "1234567890abcdefghijkelmnopqrstuvwxyz",
	}

	if _, err := d.Insert(insertData); err != nil {
		t.Errorf("ERROR: %s", err)
	}

	if e, err := d.SelectByID(id); err != nil {
		t.Errorf("ERROR: %s", err)
	} else if !reflect.DeepEqual(e, insertData) {
		t.Errorf("ERROR: %+v \n!= \n%+v", e, insertData)
	}

	if _, err := d.Delete(id); err != nil {
		t.Errorf("ERROR: %s", err)
	}

	if _, err := d.SelectByID(id); err != sql.ErrNoRows {
		t.Errorf("ERROR: %s", "Deleteしたのにnilじゃない")
	}
}

func TestDate(t *testing.T) {
	db, err := goma.Open("config.json")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	id := testID

	// date
	d := dao.GomaDateTypes(db)

	dateColumnsTime, _ := time.ParseInLocation(
		"2006-01-02",
		"2015-04-18",
		time.FixedZone("", 0),
	)

	timeStampColumnsTime, _ := time.ParseInLocation(
		"2006-01-02 15:04:05.999999",
		"2015-04-18 14:06:33.456791",
		time.FixedZone("", 0),
	)

	insertData := entity.GomaDateTypesEntity{
		ID:               id,
		DateColumns:      dateColumnsTime,
		TimestampColumns: timeStampColumnsTime,
	}

	if _, err := d.Insert(insertData); err != nil {
		t.Errorf("ERROR: %s", err)
	}

	if e, err := d.SelectByID(id); err != nil {
		t.Errorf("ERROR: %s", err)
	} else if !reflect.DeepEqual(e, insertData) {
		t.Errorf("ERROR: \n%+v \n!= \n%+v", e, insertData)
	}

	if _, err := d.Delete(id); err != nil {
		t.Errorf("ERROR: %s", err)
	}

	if _, err := d.SelectByID(id); err != sql.ErrNoRows {
		t.Errorf("ERROR: %s", "Deleteしたのにnilじゃない")
	}
}

func TestTx(t *testing.T) {
	db, err := goma.Open("config.json")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	id := testID

	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}

	// string
	dtx := dao.TxGomaStringTypes(tx)

	e := entity.GomaStringTypesEntity{
		ID:             id,
		TextColumns:    "あいうえおかきくけこ",
		CharColumns:    "a",
		VarcharColumns: "1234567890abcdefghijkelmnopqrstuvwxyz",
	}

	if _, err := dtx.Insert(e); err != nil {
		t.Errorf("ERROR: %s", err)
	}

	// Rollback（insertを無効にする）
	if err := tx.Rollback(); err != nil {
		t.Errorf("ERROR: %s", err)
	}

	// rollbackしたのでトランザクション生成しなおす
	tx, err = db.Begin()
	if err != nil {
		t.Errorf("ERROR: %s", err)
	}

	// string
	dtx = dao.TxGomaStringTypes(tx)

	// Rollbackでnilのはず
	if _, err := dtx.SelectByID(id); err != sql.ErrNoRows {
		t.Errorf("ERROR: %s", "Deleteしたのにnilじゃない")
	}

	// Insertする
	if _, err := dtx.Insert(e); err != nil {
		t.Errorf("ERROR: %s", err)
	}

	// Commit
	if err := tx.Commit(); err != nil {
		t.Errorf("ERROR: %s", err)
	}

	// commitしたのでtxじゃないdao
	d := dao.GomaStringTypes(db)

	// Commitしたのでnilじゃない
	if _, err := d.SelectByID(id); err == sql.ErrNoRows {
		t.Errorf("ERROR: %s", "Commitしたのにnil")
	}

	_, err = d.Delete(id)
	if err != nil {
		t.Errorf("ERROR: %s", err)
	}

	if _, err := d.SelectByID(id); err != sql.ErrNoRows {
		t.Errorf("ERROR: %s", "Deleteしたのにnilじゃない")
	}
}