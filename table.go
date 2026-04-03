package adapter

import "github.com/gogf/gf/v2/database/gdb"

// CasbinRule 数据库实体结构
type CasbinRule struct {
	Id    int    `json:"id"    orm:"id"    dc:"主键ID"`
	PType string `json:"pType" orm:"ptype"  dc:"策略类型"`
	V0    string `json:"v0"    orm:"v0"     dc:"字段0"`
	V1    string `json:"v1"    orm:"v1"     dc:"字段1"`
	V2    string `json:"v2"    orm:"v2"     dc:"字段2"`
	V3    string `json:"v3"    orm:"v3"     dc:"字段3"`
	V4    string `json:"v4"    orm:"v4"     dc:"字段4"`
	V5    string `json:"v5"    orm:"v5"     dc:"字段5"`
}

// TableFields 定义表字段信息
var TableFields = map[string]*gdb.TableField{
	"id": {
		Index: 0,
		Name:  "id",
		Type:  "int",
		Null:  false,
		Key:   "PRI",
		Extra: "auto_increment",
	},
	"ptype": {
		Index: 1,
		Name:  "ptype",
		Type:  "varchar(100)",
		Null:  false,
		Key:   "MUL",
	},
	"v0": {
		Index: 2,
		Name:  "v0",
		Type:  "varchar(100)",
		Null:  false,
	},
	"v1": {
		Index: 3,
		Name:  "v1",
		Type:  "varchar(100)",
		Null:  false,
	},
	"v2": {
		Index: 4,
		Name:  "v2",
		Type:  "varchar(100)",
		Null:  false,
	},
	"v3": {
		Index: 5,
		Name:  "v3",
		Type:  "varchar(100)",
		Null:  false,
	},
	"v4": {
		Index: 6,
		Name:  "v4",
		Type:  "varchar(100)",
		Null:  false,
	},
	"v5": {
		Index: 7,
		Name:  "v5",
		Type:  "varchar(100)",
		Null:  false,
	},
}

// GetCreateTableSQL 获取建表SQL（通用版本）
func GetCreateTableSQL(tableName string, fieldNames *FieldName) string {
	return `CREATE TABLE IF NOT EXISTS ` + tableName + ` (
	` + fieldNames.PType + ` VARCHAR(100) NOT NULL DEFAULT '',
	` + fieldNames.V0 + ` VARCHAR(100) NOT NULL DEFAULT '',
	` + fieldNames.V1 + ` VARCHAR(100) NOT NULL DEFAULT '',
	` + fieldNames.V2 + ` VARCHAR(100) NOT NULL DEFAULT '',
	` + fieldNames.V3 + ` VARCHAR(100) NOT NULL DEFAULT '',
	` + fieldNames.V4 + ` VARCHAR(100) NOT NULL DEFAULT '',
	` + fieldNames.V5 + ` VARCHAR(100) NOT NULL DEFAULT '',
	id INT AUTO_INCREMENT PRIMARY KEY,
	UNIQUE KEY unique_index (` + fieldNames.PType + `, ` + fieldNames.V0 + `, ` + fieldNames.V1 + `, ` + fieldNames.V2 + `, ` + fieldNames.V3 + `, ` + fieldNames.V4 + `, ` + fieldNames.V5 + `)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
}
