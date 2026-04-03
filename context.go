package adapter

import (
	"context"
	"fmt"

	"github.com/casbin/casbin/v3/model"
)

// LoadPolicyCtx 带 context 的加载策略
func (a *Adapter) LoadPolicyCtx(ctx context.Context, model model.Model) error {
	var rules []CasbinRule

	err := a.db.Ctx(ctx).
		Model(a.tableName).
		OrderAsc("id").
		Scan(&rules)
	if err != nil {
		return fmt.Errorf("gf-adapter-casbin3: failed to load policy: %w", err)
	}

	for _, rule := range rules {
		loadPolicyLine(rule, model)
	}

	return nil
}

// SavePolicyCtx 带 context 的保存策略
func (a *Adapter) SavePolicyCtx(ctx context.Context, model model.Model) error {
	// 先删除所有现有策略
	_, err := a.db.Ctx(ctx).
		Model(a.tableName).
		Delete()
	if err != nil {
		return fmt.Errorf("gf-adapter-casbin3: failed to clear policy: %w", err)
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

	// g 策略 (角色/组)
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
		_, err = a.db.Ctx(ctx).
			Model(a.tableName).
			Data(rules).
			Batch(100).
			Insert()
		if err != nil {
			return fmt.Errorf("gf-adapter-casbin3: failed to save policy: %w", err)
		}
	}

	return nil
}

// AddPolicyCtx 带 context 的添加策略
func (a *Adapter) AddPolicyCtx(ctx context.Context, sec, ptype string, rule []string) error {
	r := a.buildCasbinRule(ptype, rule)

	_, err := a.db.Ctx(ctx).
		Model(a.tableName).
		Data(r).
		Insert()
	if err != nil {
		return fmt.Errorf("gf-adapter-casbin3: failed to add policy: %w", err)
	}

	return nil
}

// RemovePolicyCtx 带 context 的删除策略
func (a *Adapter) RemovePolicyCtx(ctx context.Context, sec, ptype string, rule []string) error {
	query := a.db.Ctx(ctx).
		Model(a.tableName).
		Where(a.fieldNames.PType, ptype)

	// 构建删除条件
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

	return nil
}

// RemoveFilteredPolicyCtx 带 context 的条件删除
func (a *Adapter) RemoveFilteredPolicyCtx(ctx context.Context, sec, ptype string, fieldIndex int, fieldValues ...string) error {
	query := a.db.Ctx(ctx).
		Model(a.tableName)

	if ptype != "" {
		query = query.Where(a.fieldNames.PType, ptype)
	}

	// 根据 fieldIndex 添加条件
	columns := []string{
		a.fieldNames.V0,
		a.fieldNames.V1,
		a.fieldNames.V2,
		a.fieldNames.V3,
		a.fieldNames.V4,
		a.fieldNames.V5,
	}

	for i, fieldValue := range fieldValues {
		if fieldValue != "" && fieldIndex+i < len(columns) {
			query = query.Where(columns[fieldIndex+i], fieldValue)
		}
	}

	_, err := query.Delete()
	if err != nil {
		return fmt.Errorf("gf-adapter-casbin3: failed to remove filtered policy: %w", err)
	}

	return nil
}

// buildCasbinRule 从规则数组构建 CasbinRule
func (a *Adapter) buildCasbinRule(ptype string, rule []string) CasbinRule {
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
	return r
}
