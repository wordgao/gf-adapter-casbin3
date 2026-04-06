package adapter

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
)

// UpdatePolicy 更新单条策略
func (a *Adapter) UpdatePolicy(sec, ptype string, oldRule, newRule []string) error {
	return a.UpdatePolicyCtx(context.Background(), sec, ptype, oldRule, newRule)
}

// UpdatePolicyCtx 带 context 的更新策略
func (a *Adapter) UpdatePolicyCtx(ctx context.Context, sec, ptype string, oldRule, newRule []string) error {
	// 构建更新数据
	newData := a.buildCasbinRuleMap(newRule)

	// 构建查询条件
	query := a.db.Ctx(ctx).
		Model(a.tableName).
		Where(a.fieldNames.PType, ptype)

	for i, v := range oldRule {
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

	_, err := query.Data(newData).Update()
	if err != nil {
		return fmt.Errorf("gf-adapter-casbin3: failed to update policy: %w", err)
	}

	return nil
}

// UpdatePolicies 批量更新策略
func (a *Adapter) UpdatePolicies(sec, ptype string, oldRules, newRules [][]string) error {
	return a.UpdatePoliciesCtx(context.Background(), sec, ptype, oldRules, newRules)
}

// UpdatePoliciesCtx 带 context 的批量更新
func (a *Adapter) UpdatePoliciesCtx(ctx context.Context, sec, ptype string, oldRules, newRules [][]string) error {
	if len(oldRules) != len(newRules) {
		return fmt.Errorf("gf-adapter-casbin3: oldRules and newRules length mismatch")
	}

	// 使用事务
	return a.db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for i := range oldRules {
			newData := a.buildCasbinRuleMap(newRules[i])

			query := a.db.Ctx(ctx).
				Model(a.tableName).
				Where(a.fieldNames.PType, ptype)

			for j, v := range oldRules[i] {
				switch j {
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

			_, err := query.Data(newData).Update()
			if err != nil {
				return fmt.Errorf("gf-adapter-casbin3: failed to update policy at index %d: %w", i, err)
			}
		}

		return nil
	})
}

// UpdateFilteredPolicies 按条件更新策略
func (a *Adapter) UpdateFilteredPolicies(sec, ptype string, newRules [][]string, fieldIndex int, fieldValues ...string) ([][]string, error) {
	return a.UpdateFilteredPoliciesCtx(context.Background(), sec, ptype, newRules, fieldIndex, fieldValues...)
}

// UpdateFilteredPoliciesCtx 带 context 的条件更新
func (a *Adapter) UpdateFilteredPoliciesCtx(ctx context.Context, sec, ptype string, newRules [][]string, fieldIndex int, fieldValues ...string) ([][]string, error) {
	// 先查询旧策略
	oldRules, err := a.queryFilteredPolicies(ctx, ptype, fieldIndex, fieldValues...)
	if err != nil {
		return nil, err
	}

	// 删除旧策略
	_, err = a.db.Ctx(ctx).
		Model(a.tableName).
		Where(buildFilterConditions(a, ptype, fieldIndex, fieldValues...)).
		Delete()
	if err != nil {
		return nil, fmt.Errorf("gf-adapter-casbin3: failed to delete old policies: %w", err)
	}

	// 插入新策略
	if len(newRules) > 0 {
		var rules []CasbinRule
		for _, rule := range newRules {
			rules = append(rules, a.buildCasbinRule(ptype, rule))
		}

		_, err = a.db.Ctx(ctx).
			Model(a.tableName).
			Data(rules).
			Batch(100).
			Insert()
		if err != nil {
			return nil, fmt.Errorf("gf-adapter-casbin3: failed to insert new policies: %w", err)
		}
	}

	return oldRules, nil
}

// queryFilteredPolicies 查询符合条件的策略
func (a *Adapter) queryFilteredPolicies(ctx context.Context, ptype string, fieldIndex int, fieldValues ...string) ([][]string, error) {
	var rules []CasbinRule

	query := a.db.Ctx(ctx).
		Model(a.tableName).
		Where(a.fieldNames.PType, ptype)

	// 构建过滤条件
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

	err := query.OrderAsc("id").Scan(&rules)
	if err != nil {
		return nil, fmt.Errorf("gf-adapter-casbin3: failed to query filtered policies: %w", err)
	}

	// 转换为规则数组
	var result [][]string
	for _, rule := range rules {
		r := []string{}
		if rule.V0 != "" {
			r = append(r, rule.V0)
		}
		if rule.V1 != "" {
			r = append(r, rule.V1)
		}
		if rule.V2 != "" {
			r = append(r, rule.V2)
		}
		if rule.V3 != "" {
			r = append(r, rule.V3)
		}
		if rule.V4 != "" {
			r = append(r, rule.V4)
		}
		if rule.V5 != "" {
			r = append(r, rule.V5)
		}
		if len(r) > 0 {
			result = append(result, r)
		}
	}

	return result, nil
}

// buildCasbinRuleMap 从规则数组构建 map
func (a *Adapter) buildCasbinRuleMap(rule []string) map[string]interface{} {
	m := make(map[string]interface{})
	for i, v := range rule {
		switch i {
		case 0:
			m[a.fieldNames.V0] = v
		case 1:
			m[a.fieldNames.V1] = v
		case 2:
			m[a.fieldNames.V2] = v
		case 3:
			m[a.fieldNames.V3] = v
		case 4:
			m[a.fieldNames.V4] = v
		case 5:
			m[a.fieldNames.V5] = v
		}
	}
	return m
}

// buildFilterConditions 构建过滤条件
func buildFilterConditions(a *Adapter, ptype string, fieldIndex int, fieldValues ...string) map[string]interface{} {
	conditions := make(map[string]interface{})
	conditions[a.fieldNames.PType] = ptype

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
			conditions[columns[fieldIndex+i]] = fieldValue
		}
	}

	return conditions
}
