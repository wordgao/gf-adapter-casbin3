package adapter

import (
	"context"
	"fmt"
	"strings"

	"github.com/casbin/casbin/v3/model"
	"github.com/casbin/casbin/v3/persist"
	"github.com/gogf/gf/v2/database/gdb"
)

// 确保 Adapter 实现所有接口
var (
	_ persist.Adapter              = (*Adapter)(nil)
	_ persist.ContextAdapter       = (*Adapter)(nil)
	_ persist.UpdatableAdapter     = (*Adapter)(nil)
	_ persist.BatchAdapter         = (*Adapter)(nil)
	_ persist.FilteredAdapter      = (*Adapter)(nil)
)

// Adapter Casbin v3 适配器
type Adapter struct {
	db         gdb.DB       // GoFrame 数据库实例
	tableName  string       // 表名，默认 casbin_rule
	fieldNames *FieldName   // 字段名配置
	isFiltered bool         // 是否支持过滤
}

// NewAdapter 创建适配器实例
func NewAdapter(opts Options) (*Adapter, error) {
	if opts.GDB == nil {
		return nil, ErrDBRequired
	}

	// 设置默认值
	if opts.TableName == "" {
		opts.TableName = "casbin_rule"
	}
	if opts.FieldName == nil {
		opts.FieldName = DefaultFieldName()
	}
	if opts.AutoCreate {
		// 自动建表
		if err := createTable(opts.GDB, opts.TableName, opts.FieldName); err != nil {
			return nil, err
		}
	}

	return &Adapter{
		db:         opts.GDB,
		tableName:  opts.TableName,
		fieldNames: opts.FieldName,
		isFiltered: false,
	}, nil
}

// NewAdapterByDB 通过 GoFrame 数据库实例创建适配器
func NewAdapterByDB(db gdb.DB, opts ...Option) (*Adapter, error) {
	if db == nil {
		return nil, ErrDBRequired
	}

	options := Options{
		GDB:        db,
		TableName:  "casbin_rule",
		FieldName:  DefaultFieldName(),
		AutoCreate: true,
	}

	for _, opt := range opts {
		opt(&options)
	}

	return NewAdapter(options)
}

// NewAdapterByGroup 通过 GoFrame 配置分组创建适配器
func NewAdapterByGroup(group string, opts ...Option) (*Adapter, error) {
	db, err := gdb.NewByGroup(group)
	if err != nil {
		return nil, fmt.Errorf("gf-adapter-casbin3: failed to create db from group %s: %w", group, err)
	}
	return NewAdapterByDB(db, opts...)
}

// createTable 自动创建表
func createTable(db gdb.DB, tableName string, fieldNames *FieldName) error {
	ctx := context.Background()

	// 检查表是否存在
	exists, err := db.GetCore().Tables(ctx)
	if err != nil {
		return err
	}

	for _, table := range exists {
		if strings.EqualFold(table, tableName) {
			// 表已存在，无需创建
			return nil
		}
	}

	// 获取建表 SQL
	sql := GetCreateTableSQL(tableName, fieldNames)

	// 执行建表
	_, err = db.Exec(ctx, sql)
	return err
}

// buildColumnList 构建字段列表
func (a *Adapter) buildColumnList() string {
	return fmt.Sprintf("%s, %s, %s, %s, %s, %s, %s",
		a.fieldNames.PType,
		a.fieldNames.V0,
		a.fieldNames.V1,
		a.fieldNames.V2,
		a.fieldNames.V3,
		a.fieldNames.V4,
		a.fieldNames.V5)
}

// ruleToLine 将 CasbinRule 转换为策略行
func ruleToLine(rule CasbinRule) string {
	line := rule.PType
	if rule.V0 != "" {
		line += ", " + rule.V0
	}
	if rule.V1 != "" {
		line += ", " + rule.V1
	}
	if rule.V2 != "" {
		line += ", " + rule.V2
	}
	if rule.V3 != "" {
		line += ", " + rule.V3
	}
	if rule.V4 != "" {
		line += ", " + rule.V4
	}
	if rule.V5 != "" {
		line += ", " + rule.V5
	}
	return line
}

// loadPolicyLine 加载单条策略到模型
func loadPolicyLine(rule CasbinRule, model model.Model) {
	line := ruleToLine(rule)
	persist.LoadPolicyLine(line, model)
}

// lineToRule 将策略行转换为 CasbinRule
func lineToRule(line string) CasbinRule {
	if line == "" {
		return CasbinRule{}
	}

	tokens := strings.Split(line, ",")
	for i := range tokens {
		tokens[i] = strings.TrimSpace(tokens[i])
	}

	rule := CasbinRule{}
	if len(tokens) > 0 {
		rule.PType = tokens[0]
	}
	if len(tokens) > 1 {
		rule.V0 = tokens[1]
	}
	if len(tokens) > 2 {
		rule.V1 = tokens[2]
	}
	if len(tokens) > 3 {
		rule.V2 = tokens[3]
	}
	if len(tokens) > 4 {
		rule.V3 = tokens[4]
	}
	if len(tokens) > 5 {
		rule.V4 = tokens[5]
	}
	if len(tokens) > 6 {
		rule.V5 = tokens[6]
	}

	return rule
}

// LoadPolicy 加载所有策略
func (a *Adapter) LoadPolicy(model model.Model) error {
	return a.LoadPolicyCtx(context.Background(), model)
}

// SavePolicy 保存所有策略
func (a *Adapter) SavePolicy(model model.Model) error {
	return a.SavePolicyCtx(context.Background(), model)
}

// AddPolicy 添加单条策略
func (a *Adapter) AddPolicy(sec, ptype string, rule []string) error {
	return a.AddPolicyCtx(context.Background(), sec, ptype, rule)
}

// RemovePolicy 删除单条策略
func (a *Adapter) RemovePolicy(sec, ptype string, rule []string) error {
	return a.RemovePolicyCtx(context.Background(), sec, ptype, rule)
}

// RemoveFilteredPolicy 按条件删除策略
func (a *Adapter) RemoveFilteredPolicy(sec, ptype string, fieldIndex int, fieldValues ...string) error {
	return a.RemoveFilteredPolicyCtx(context.Background(), sec, ptype, fieldIndex, fieldValues...)
}

// IsFiltered 是否支持过滤
func (a *Adapter) IsFiltered() bool {
	return a.isFiltered
}

// SetFiltered 设置过滤状态
func (a *Adapter) SetFiltered(filtered bool) {
	a.isFiltered = filtered
}

// GetTableName 获取表名
func (a *Adapter) GetTableName() string {
	return a.tableName
}

// GetDB 获取数据库实例
func (a *Adapter) GetDB() gdb.DB {
	return a.db
}
