package adapter

import (
	"context"
	"fmt"

	"github.com/casbin/casbin/v3/model"
	"github.com/gogf/gf/v2/database/gdb"
)

// TransactionalAdapter 事务适配器
type TransactionalAdapter struct {
	*Adapter
	tx gdb.TX
}

// NewTransactionalAdapter 创建事务适配器
func NewTransactionalAdapter(adapter *Adapter, tx gdb.TX) *TransactionalAdapter {
	return &TransactionalAdapter{
		Adapter: adapter,
		tx:      tx,
	}
}

// BeginTransaction 开始事务
func (a *Adapter) BeginTransaction(ctx context.Context) (*TransactionalAdapter, error) {
	tx, err := a.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("gf-adapter-casbin3: failed to begin transaction: %w", err)
	}

	return NewTransactionalAdapter(a, tx), nil
}

// Commit 提交事务
func (ta *TransactionalAdapter) Commit() error {
	return ta.tx.Commit()
}

// Rollback 回滚事务
func (ta *TransactionalAdapter) Rollback() error {
	return ta.tx.Rollback()
}

// LoadPolicyCtx 带事务的加载策略
func (ta *TransactionalAdapter) LoadPolicyCtx(ctx context.Context, model model.Model) error {
	var rules []CasbinRule

	err := ta.tx.Ctx(ctx).
		Model(ta.tableName).
		OrderAsc("id").
		Scan(&rules)
	if err != nil {
		return fmt.Errorf("gf-adapter-casbin3: failed to load policy in transaction: %w", err)
	}

	for _, rule := range rules {
		loadPolicyLine(rule, model)
	}

	return nil
}

// SavePolicyCtx 带事务的保存策略
func (ta *TransactionalAdapter) SavePolicyCtx(ctx context.Context, model model.Model) error {
	// 先删除所有现有策略
	_, err := ta.tx.Ctx(ctx).
		Model(ta.tableName).
		Delete()
	if err != nil {
		return fmt.Errorf("gf-adapter-casbin3: failed to clear policy in transaction: %w", err)
	}

	// 收集所有策略
	var rules []CasbinRule

	// p 策略
	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			r := CasbinRule{PType: ptype}
			for i, v := range rule {
				switch i {
				case 0:
					r.V0 = v
				case 1:
					r.V1 = v
				case 2:
					r.V2 = v
				case 3:
					r.V3 = v
				case 4:
					r.V4 = v
				case 5:
					r.V5 = v
				}
			}
			rules = append(rules, r)
		}
	}

	// g 策略
	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			r := CasbinRule{PType: ptype}
			for i, v := range rule {
				switch i {
				case 0:
					r.V0 = v
				case 1:
					r.V1 = v
				case 2:
					r.V2 = v
				case 3:
					r.V3 = v
				case 4:
					r.V4 = v
				case 5:
					r.V5 = v
				}
			}
			rules = append(rules, r)
		}
	}

	// 批量插入
	if len(rules) > 0 {
		_, err = ta.tx.Ctx(ctx).
			Model(ta.tableName).
			Data(rules).
			Batch(100).
			Insert()
		if err != nil {
			return fmt.Errorf("gf-adapter-casbin3: failed to save policy in transaction: %w", err)
		}
	}

	return nil
}

// AddPolicyCtx 带事务的添加策略
func (ta *TransactionalAdapter) AddPolicyCtx(ctx context.Context, sec, ptype string, rule []string) error {
	r := ta.buildCasbinRule(ptype, rule)

	_, err := ta.tx.Ctx(ctx).
		Model(ta.tableName).
		Data(r).
		Insert()
	if err != nil {
		return fmt.Errorf("gf-adapter-casbin3: failed to add policy in transaction: %w", err)
	}

	return nil
}

// RemovePolicyCtx 带事务的删除策略
func (ta *TransactionalAdapter) RemovePolicyCtx(ctx context.Context, sec, ptype string, rule []string) error {
	query := ta.tx.Ctx(ctx).
		Model(ta.tableName).
		Where(ta.fieldNames.PType, ptype)

	for i, v := range rule {
		switch i {
		case 0:
			query = query.Where(ta.fieldNames.V0, v)
		case 1:
			query = query.Where(ta.fieldNames.V1, v)
		case 2:
			query = query.Where(ta.fieldNames.V2, v)
		case 3:
			query = query.Where(ta.fieldNames.V3, v)
		case 4:
			query = query.Where(ta.fieldNames.V4, v)
		case 5:
			query = query.Where(ta.fieldNames.V5, v)
		}
	}

	_, err := query.Delete()
	if err != nil {
		return fmt.Errorf("gf-adapter-casbin3: failed to remove policy in transaction: %w", err)
	}

	return nil
}
