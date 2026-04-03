package adapter

import "errors"

var (
	// ErrDBRequired 数据库实例未提供
	ErrDBRequired = errors.New("gf-adapter-casbin3: database instance is required")
	// ErrTableExists 表已存在
	ErrTableExists = errors.New("gf-adapter-casbin3: table already exists")
	// ErrNotImplemented 方法未实现
	ErrNotImplemented = errors.New("gf-adapter-casbin3: not implemented")
	// ErrInvalidRule 无效的策略规则
	ErrInvalidRule = errors.New("gf-adapter-casbin3: invalid policy rule")
	// ErrEmptyTableName 表名为空
	ErrEmptyTableName = errors.New("gf-adapter-casbin3: table name cannot be empty")
)
