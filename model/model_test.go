package model

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"sync"
	"testing"
)

func Test_parseModel(t *testing.T) {
	testCases := []struct {
		name      string
		entity    any
		wantModel *Model
		wantErr   error
		opts      []Option
	}{
		{
			name:    "test Model",
			entity:  TestModel{},
			wantErr: errors.New("orm: 只支持指向结构体的一级指针"),
		},
		{
			name:   "pointer",
			entity: &TestModel{},
			wantModel: &Model{
				TableName: "test_model",
				FieldMap: map[string]*Field{
					"Id": {
						ColName: "id",
					},
					"FirstName": {
						ColName: "first_name",
					},
					"LastName": {
						ColName: "last_name",
					},
					"Age": {
						ColName: "age",
					},
				},
			},
		},
		{
			name: "table",
		},
	}
	r := &registry{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.Registry(tc.entity, tc.opts...)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantModel, m)
		})
	}
}

func TestRegistry_get(t *testing.T) {
	testCases := []struct {
		name      string
		entity    any
		wantModel *Model
		wantErr   error
		cacheSize int
	}{
		{
			name:   "pointer",
			entity: &TestModel{},
			wantModel: &Model{
				TableName: "test_model",
				FieldMap: map[string]*Field{
					"Id": {
						ColName: "id",
					},
					"FirstName": {
						ColName: "first_name",
					},
					"LastName": {
						ColName: "last_name",
					},
					"Age": {
						ColName: "age",
					},
				},
			},
			cacheSize: 1,
		},
		{
			name: "tag",
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"column=first_name_t"`
				}
				return &TagTable{}
			}(),
			wantModel: &Model{
				TableName: "tag_table",
				FieldMap: map[string]*Field{
					"FirstName": {
						ColName: "first_name_t",
					},
				},
			},
		},
		{
			name:   "table name",
			entity: &CustomTableName{},
			wantModel: &Model{
				TableName: "custom_table_name_t",
				FieldMap: map[string]*Field{
					"FirstName": {
						ColName: "first_name",
					},
				},
			},
		},
	}
	r := InitRegistry()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.Get(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantModel, m)
			//assert.Equal(t, tc.cacheSize, len(r.models))

			//typ := reflect.TypeOf(tc.entity)
			////cache, ok := r..Load(typ)
			//assert.True(t, ok)
			//assert.Equal(t, tc.wantModel, cache)
		})
	}
}

type CustomTableName struct {
	FirstName string
}

func (c CustomTableName) TableName() string {
	return "custom_table_name_t"
}

func TestModelWithTableName(t *testing.T) {
	r := InitRegistry()
	m, err := r.Registry(&TestModel{}, WithTableName("abc"))
	require.NoError(t, err)
	assert.Equal(t, "abc", m.TableName)

}

func TestModelWithColumnName(t *testing.T) {
	testCases := []struct {
		name    string
		field   string
		colName string

		wantColName string
		wantErr     error
	}{
		{
			name:        "column",
			field:       "FirstName",
			colName:     "first_name_ccc",
			wantColName: "first_name",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := InitRegistry()
			m, err := r.Registry(&TestModel{}, WithColumnName(tc.field, tc.colName))
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			fd, ok := m.FieldMap[tc.field]
			require.True(t, ok)
			assert.Equal(t, tc.wantColName, fd.ColName)
		})
	}
}

func Test_registry_parseTag(t *testing.T) {
	type fields struct {
		models sync.Map
	}
	type args struct {
		tag reflect.StructTag
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]string
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &registry{
				models: tt.fields.models,
			}
			got, err := r.parseTag(tt.args.tag)
			if !tt.wantErr(t, err, fmt.Sprintf("parseTag(%v)", tt.args.tag)) {
				return
			}
			assert.Equalf(t, tt.want, got, "parseTag(%v)", tt.args.tag)
		})
	}
}

func Test_underscoreName(t *testing.T) {
	type args struct {
		tableName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, underscoreName(tt.args.tableName), "underscoreName(%v)", tt.args.tableName)
		})
	}
}

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}
