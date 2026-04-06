package adapter

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
)

// AddPolicies 批量添加策略
func (a *Adapter) AddPolicies(sec, ptype string, rules [][]string) error {
	return a.AddPoliciesCtx(context.Background(), sec, ptype, rules)
}

// AddPoliciesCtx 带 context 的批量添加
func (a *Adapter) AddPoliciesCtx(ctx context.Context, sec, ptype string, rules [][]string) error {
	if len(rules) == 0 {
		return nil
	}

	// 构建批量插入数据
	var casbinRules []CasbinRule
	for _, rule := range rules {
		casbinRules = append(casbinRules, a.buildCasbinRule(ptype, rule))
	}

	// 批量插入
	_, err := a.db.Ctx(ctx).
		Model(a.tableName).
		Data(casbinRules).
		Batch(100).
		Insert()
	if err != nil {
		return fmt.Errorf("gf-adapter-casbin3: failed to add policies: %w", err)
	}

	return nil
}

// RemovePolicies 批量删除策略
func (a *Adapter) RemovePolicies(sec, ptype string, rules [][]string) error {
	return a.RemovePoliciesCtx(context.Background(), sec, ptype, rules)
}

// RemovePoliciesCtx 带 context 的批量删除
func (a *Adapter) RemovePoliciesCtx(ctx context.Context, sec, ptype string, rules [][]string) error {
	if len(rules) == 0 {
		return nil
	}

	// 使用事务
	return a.db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, rule := range rules {
			query := a.db.Ctx(ctx).
				Model(a.tableName).
				Where(a.fieldNames.PType, ptype)

			for i, v := range rule {
				switch i {
				case 0:
					query = query.Where(a.fieldNames.V0, v)
				case 1:
					query = query.Where(a.fieldNames.V1, v)
				case 2:
					query = query.Where(a.fieldNames.V2, v)
				case 3:
					query = query.Where(a.fieldNames.V3, v)
				case 4:
					query = query.Where(a.fieldNames.V4, v)
				case 5:
					query = query.Where(a.fieldNames.V5, v)
				}
			}

			_, err := query.Delete()
			if err != nil {
				return fmt.Errorf("gf-adapter-casbin3: failed to remove policy: %w", err)
			}
		}

		return nil
	})
}
