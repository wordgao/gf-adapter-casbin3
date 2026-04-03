package adapter

import "github.com/gogf/gf/v2/database/gdb"

// FieldName 自定义字段名配置
type FieldName struct {
	PType string // ptype 字段名，默认 ptype
	V0    string // v0 字段名，默认 v0
	V1    string // v1 字段名，默认 v1
	V2    string // v2 字段名，默认 v2
	V3    string // v3 字段名，默认 v3
	V4    string // v4 字段名，默认 v4
	V5    string // v5 字段名，默认 v5
}

// DefaultFieldName 返回默认字段名配置
func DefaultFieldName() *FieldName {
	return &FieldName{
		PType: "ptype",
		V0:    "v0",
		V1:    "v1",
		V2:    "v2",
		V3:    "v3",
		V4:    "v4",
		V5:    "v5",
	}
}

// Options 适配器配置选项
type Options struct {
	GDB        gdb.DB     // GoFrame 数据库实例（必需）
	TableName  string     // 自定义表名，默认 casbin_rule
	FieldName  *FieldName // 自定义字段名
	AutoCreate bool       // 自动建表，默认 true
}

// Option 函数式选项
type Option func(*Options)

// WithTableName 设置表名
func WithTableName(tableName string) Option {
	return func(o *Options) {
		o.TableName = tableName
	}
}

// WithFieldName 设置字段名
func WithFieldName(fieldName *FieldName) Option {
	return func(o *Options) {
		o.FieldName = fieldName
	}
}

// WithAutoCreate 设置自动建表
func WithAutoCreate(autoCreate bool) Option {
	return func(o *Options) {
		o.AutoCreate = autoCreate
	}
}
