package db

import (
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func Test_dbHandler_connect(t *testing.T) {
	type args struct {
		openConnection func(dialect string, args ...interface{}) (db *gorm.DB, err error)
	}
	tests := []struct {
		name      string
		value     string
		dbHandler *dbHandler
		args      args
	}{
		struct {
			name      string
			value     string
			dbHandler *dbHandler
			args      args
		}{
			name:      "Connected with DB.Values",
			value:     "dummy value data",
			dbHandler: &dbHandler{},
			args: args{openConnection: func(dialect string, args ...interface{}) (db *gorm.DB, err error) {
				return &gorm.DB{Value: "dummy value data"}, nil
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.dbHandler.connect(tt.args.openConnection)
			if tt.dbHandler.database.Value != tt.value {
				t.Errorf("db.Value:%v should be equal to mock value:%v", tt.dbHandler.database.Value, tt.value)
			}
		})
	}
}
