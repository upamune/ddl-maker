package ddlmaker

import (
	"bytes"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/upamune/ddl-maker/dialect"
	"github.com/upamune/ddl-maker/dialect/mysql"
)

type Test1 struct {
	ID        uint64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t1 Test1) PrimaryKey() dialect.PrimaryKey {
	return mysql.AddPrimaryKey("id")
}

type Test2 struct {
	ID        uint64
	Test1ID   uint64
	Comment   sql.NullString `ddl:"null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t2 *Test2) PrimaryKey() dialect.PrimaryKey {
	return mysql.AddPrimaryKey("id", "created_at")
}

func TestNew(t *testing.T) {
	conf := Config{}
	_, err := New(conf)
	if err == nil {
		t.Fatal("Not set driver name", err)
	}

	conf = Config{
		DB: DBConfig{Driver: "dummy"},
	}
	_, err = New(conf)
	if err == nil {
		t.Fatal("Set unsupport driver name", err)
	}

	conf = Config{
		DB: DBConfig{Driver: "mysql"},
	}
	_, err = New(conf)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAddStruct(t *testing.T) {
	dm, err := New(Config{
		DB: DBConfig{Driver: "mysql"},
	})
	if err != nil {
		t.Fatal(err)
	}
	err = dm.AddStruct(nil)
	if err == nil {
		t.Fatal("nil is not support")
	}

	dm.AddStruct(Test1{}, Test2{})
	if len(dm.Structs) != 2 {
		t.Fatal("[error] add stuct")
	}

	err = dm.AddStruct(Test1{})
	if err != nil {
		t.Fatal("[error] add duplicate struct")
	}
}

func TestGenerate(t *testing.T) {
	m := mysql.MySQL{}
	generatedDDL := fmt.Sprintf(`%s
DROP TABLE IF EXISTS %s;

CREATE TABLE %s (
    %s BIGINT unsigned NOT NULL,
    %s VARCHAR(191) NOT NULL,
    %s DATETIME NOT NULL,
    %s DATETIME NOT NULL,
    PRIMARY KEY (%s)
) ENGINE=InnoDB DEFAULT CHARACTER SET utf8mb4;

%s`, m.HeaderTemplate(), m.Quote("test1"), m.Quote("test1"), m.Quote("id"), m.Quote("name"), m.Quote("created_at"), m.Quote("updated_at"), m.Quote("id"), m.FooterTemplate())

	generatedDDL2 := fmt.Sprintf(`%s
DROP TABLE IF EXISTS %s;

CREATE TABLE %s (
    %s BIGINT unsigned NOT NULL,
    %s BIGINT unsigned NOT NULL,
    %s VARCHAR(191) NULL,
    %s DATETIME NOT NULL,
    %s DATETIME NOT NULL,
    PRIMARY KEY (%s, %s)
) ENGINE=InnoDB DEFAULT CHARACTER SET utf8mb4;

%s`, m.HeaderTemplate(), m.Quote("test2"), m.Quote("test2"), m.Quote("id"), m.Quote("test1_id"), m.Quote("comment"), m.Quote("created_at"), m.Quote("updated_at"), m.Quote("id"), m.Quote("created_at"), m.FooterTemplate())

	dm, err := New(Config{
		DB: DBConfig{
			Driver:  "mysql",
			Engine:  "InnoDB",
			Charset: "utf8mb4",
		},
	})
	if err != nil {
		t.Fatal("error new maker", err)
	}

	err = dm.AddStruct(&Test1{})
	if err != nil {
		t.Fatal("error add struct", err)
	}
	dm.parse()

	var ddl1 bytes.Buffer
	err = dm.generate(&ddl1)
	if err != nil {
		t.Fatal("error generate ddl", err)
	}

	if ddl1.String() != generatedDDL {
		t.Log(ddl1.String())
		t.Fatalf("generatedDDL: %s \n checkDDLL: %s \n", ddl1.String(), generatedDDL)
	}

	dm2, err := New(Config{
		DB: DBConfig{
			Driver:  "mysql",
			Engine:  "InnoDB",
			Charset: "utf8mb4",
		},
	})

	err = dm2.AddStruct(&Test2{})
	if err != nil {
		t.Fatal("error add pointer struct", err)
	}
	dm2.parse()

	var ddl2 bytes.Buffer
	err = dm2.generate(&ddl2)
	if err != nil {
		t.Fatal("error generate ddl", err)
	}

	if ddl2.String() != generatedDDL2 {
		t.Fatalf("generatedDDL: %s \n checkDDLL: %s \n", ddl2.String(), generatedDDL2)
	}
}
