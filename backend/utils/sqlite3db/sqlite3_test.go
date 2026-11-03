package sqlite3db

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/NeverStopDreamingWang/goi/v3"
	"github.com/NeverStopDreamingWang/goi/v3/db"
	"github.com/NeverStopDreamingWang/goi/v3/db/sqlite3"
)

type ormCompatibilityModel struct {
	ID   *int64  `field_name:"id" field_type:"INTEGER PRIMARY KEY AUTOINCREMENT"`
	Name *string `field_name:"name" field_type:"TEXT NOT NULL"`
	Note *string `field_name:"note" field_type:"TEXT"`
}

func (ormCompatibilityModel) ModelSet() *sqlite3.Settings {
	return &sqlite3.Settings{TableName: "orm_compatibility"}
}

// TestConnectGoiV3ORM 只使用临时 SQLite 文件，不加载项目配置或应用初始化逻辑。
func TestConnectGoiV3ORM(t *testing.T) {
	server := goi.NewHTTPServer()
	oldConfig, oldDB := Config, SQLite3DB
	oldBaseDir, oldDebug, oldLog := server.Settings.BaseDir, server.Settings.Debug, goi.Log
	root := t.TempDir()
	var database *goi.Database
	var logger *goi.Logger
	t.Cleanup(func() {
		if database != nil {
			if err := database.Close(); err != nil {
				t.Errorf("close temporary database: %v", err)
			}
		}
		if logger != nil {
			if err := logger.Close(); err != nil {
				t.Errorf("close temporary logger: %v", err)
			}
		}
		Config, SQLite3DB = oldConfig, oldDB
		server.Settings.BaseDir, server.Settings.Debug, goi.Log = oldBaseDir, oldDebug, oldLog
	})

	server.Settings.BaseDir = root
	server.Settings.Debug = false
	logger = goi.NewLogger(filepath.Join(root, "orm.log"))
	goi.Log = logger
	Config = &ConfigModel{Uri: "v3-compatibility.sqlite3"}
	database = &goi.Database{Engine: "sqlite3", Connect: Connect}
	engine := db.ConnectDatabase[*sqlite3.Engine]("orm-compatibility", database)
	ctx := context.Background()
	if err := engine.DB.PingContext(ctx); err != nil {
		t.Fatalf("open temporary SQLite database (requires CGO and a C compiler): %v", err)
	}
	if engine.DB != SQLite3DB || engine.DB != database.DB() {
		t.Fatal("ORM did not reuse the connection created by Connect")
	}
	engine.Migrate(ctx, ormCompatibilityModel{})

	name, note := "before", "keep"
	result, err := engine.SetModel(ormCompatibilityModel{}).Insert(ctx, ormCompatibilityModel{Name: &name, Note: &note}, true)
	if err != nil {
		t.Fatal(err)
	}
	id, err := result.LastInsertId()
	if err != nil || id <= 0 {
		t.Fatalf("insert ID=%d, error=%v", id, err)
	}
	readRecord := func() ormCompatibilityModel {
		t.Helper()
		var record ormCompatibilityModel
		if err := engine.SetModel(ormCompatibilityModel{}).Where("id = ?", id).First(ctx, &record); err != nil {
			t.Fatal(err)
		}
		return record
	}

	updatedName := "after"
	_, err = engine.SetModel(ormCompatibilityModel{}).Fields("Name", "Note").Where("id = ?", id).
		Update(ctx, ormCompatibilityModel{Name: &updatedName}, true)
	if err != nil {
		t.Fatal(err)
	}
	got := readRecord()
	if got.Name == nil || *got.Name != updatedName || got.Note == nil || *got.Note != note {
		t.Fatalf("skipNil=true did not preserve Note: %+v", got)
	}
	_, err = engine.SetModel(ormCompatibilityModel{}).Fields("Note").Where("id = ?", id).
		Update(ctx, ormCompatibilityModel{}, false)
	if err != nil {
		t.Fatal(err)
	}
	got = readRecord()
	if got.Name == nil || *got.Name != updatedName || got.Note != nil {
		t.Fatalf("selected nil field was not written as NULL: %+v", got)
	}

	var records []ormCompatibilityModel
	if err := engine.SetModel(ormCompatibilityModel{}).Select(ctx, &records); err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 || records[0].ID == nil || *records[0].ID != id {
		t.Fatalf("unexpected slice result: %+v", records)
	}
	assertMap := func(row map[string]any) {
		t.Helper()
		if len(row) != 3 {
			t.Fatalf("unexpected map fields: %#v", row)
		}
		for _, key := range []string{"ID", "Name", "Note"} {
			if _, ok := row[key]; !ok {
				t.Fatalf("map must use Go field name %q: %#v", key, row)
			}
		}
		if value, ok := row["Name"].(*string); !ok || value == nil || *value != updatedName {
			t.Fatalf("unexpected map Name: %#v", row["Name"])
		}
		if value, ok := row["Note"].(*string); !ok || value != nil {
			t.Fatalf("unexpected nullable map Note: %#v", row["Note"])
		}
	}
	row := make(map[string]any)
	if err := engine.SetModel(ormCompatibilityModel{}).Where("id = ?", id).First(ctx, row); err != nil {
		t.Fatal(err)
	}
	assertMap(row)
	var rows []map[string]any
	if err := engine.SetModel(ormCompatibilityModel{}).Select(ctx, &rows); err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("unexpected map slice length: %d", len(rows))
	}
	assertMap(rows[0])

	rollback := errors.New("rollback compatibility test")
	err = engine.SetModel(ormCompatibilityModel{}).WithTransaction(ctx, func(tx *sqlite3.Engine, _ ...any) error {
		transactionName := "rolled-back"
		if _, err := tx.Fields("Name").Insert(ctx, ormCompatibilityModel{Name: &transactionName}, true); err != nil {
			return err
		}
		return rollback
	})
	if !errors.Is(err, rollback) {
		t.Fatalf("transaction error=%v, want rollback sentinel", err)
	}
	count, err := engine.SetModel(ormCompatibilityModel{}).Count(ctx)
	if err != nil || count != 1 {
		t.Fatalf("transaction did not roll back: count=%d, error=%v", count, err)
	}

	canceled, cancel := context.WithCancel(ctx)
	cancel()
	if _, err := engine.SetModel(ormCompatibilityModel{}).Count(canceled); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled count error=%v", err)
	}
	if _, err := engine.SetModel(ormCompatibilityModel{}).Fields("Name").Where("id = ?", id).
		Update(canceled, ormCompatibilityModel{Name: &name}, true); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled update error=%v", err)
	}
}
