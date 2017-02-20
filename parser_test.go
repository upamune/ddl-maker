package ddlmaker

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/upamune/ddl-maker/dialect"
	"github.com/upamune/ddl-maker/dialect/mysql"
)

type T1 struct {
	ID          uint64 `ddl:"auto"`
	Name        string
	Description sql.NullString `ddl:"null,text"`
	CreatedAt   time.Time
}

func (t1 T1) Table() string {
	return "test1"
}

func (t1 T1) Indexes() dialect.Indexes {
	return dialect.Indexes{
		mysql.AddUniqueIndex("token_idx", "token"),
	}
}

func (t1 T1) PrimaryKey() dialect.PrimaryKey {
	return mysql.AddPrimaryKey("id", "created_at")
}

func TestParseField(t *testing.T) {
	t1 := T1{}
	idColumn := column{
		name:     "id",
		tag:      "auto",
		typeName: "uint64",
		dialect:  mysql.MySQL{},
	}
	nameColumn := column{
		name:     "name",
		typeName: "string",
		dialect:  mysql.MySQL{},
	}
	descColumn := column{
		name:     "description",
		typeName: "sql.NullString",
		tag:      "null,text",
		dialect:  mysql.MySQL{},
	}
	createdAtColumn := column{
		name:     "created_at",
		typeName: "time.Time",
		dialect:  mysql.MySQL{},
	}
	columns := []dialect.Column{idColumn, nameColumn, descColumn, createdAtColumn}

	rt := reflect.TypeOf(t1)

	if rt.NumField() == 0 {
		t.Fatal("T1 field is 0")
	}

	for i := 0; i < rt.NumField(); i++ {
		column := parseField(rt.Field(i), mysql.MySQL{})

		if !reflect.DeepEqual(columns[i], column) {
			t.Fatalf("parsed %s: %v is different \n %s: %v", column.Name(), column, columns[i].Name(), columns[i])
		}
	}
}

func TestParseTable(t *testing.T) {
	t1 := T1{}
	d := mysql.MySQL{}

	var columns []dialect.Column
	table := parseTable(t1, columns, d)
	if table.Name() != d.Quote(t1.Table()) {
		t.Fatal("error parse table name", table.Name)
	}

	if len(table.Indexes()) != len(t1.Indexes()) {
		t.Fatal("error parse index ", table.Indexes)
	}

	if table.PrimaryKey().ToSQL() != "PRIMARY KEY (`id`, `created_at`)" {
		t.Fatal("error parse pk: ", table.PrimaryKey)
	}
}
