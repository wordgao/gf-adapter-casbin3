package adapter

import (
	"context"
	"fmt"

	"github.com/casbin/casbin/v3/model"
)

// Filter 定义过滤规则
type Filter struct {
	P []string // p 策略过滤条件
	G []string // g 策略过滤条件
}

// LoadFilteredPolicy 加载过滤后的策略
func (a *Adapter) LoadFilteredPolicy(model model.Model, filter interface{}) error {
	return a.LoadFilteredPolicyCtx(context.Background(), model, filter)
}

// LoadFilteredPolicyCtx 带 context 的加载过滤策略
func (a *Adapter) LoadFilteredPolicyCtx(ctx context.Context, model model.Model, filter interface{}) error {
	f, ok := filter.(*Filter)
	if !ok {
		return fmt.Errorf("gf-adapter-casbin3: invalid filter type, expected *Filter")
	}

	// 加载 p 策略
	if err := a.loadFilteredPolicies(ctx, model, "p", f.P); err != nil {
		return err
	}

	// 加载 g 策略
	if err := a.loadFilteredPolicies(ctx, model, "g", f.G); err != nil {
		return err
	}

	a.isFiltered = true
	return nil
}

// loadFilteredPolicies 加载指定类型的过滤策略
func (a *Adapter) loadFilteredPolicies(ctx context.Context, model model.Model, sec string, filter []string) error {
	var rules []CasbinRule

	query := a.db.Ctx(ctx).
		Model(a.tableName)

	// 添加 ptype 条件
	if sec == "p" {
		query = query.WhereLike(a.fieldNames.PType, "p%")
	} else if sec == "g" {
		query = query.WhereLike(a.fieldNames.PType, "g%")
	}

	// 添加过滤条件
	columns := []string{
		a.fieldNames.V0,
		a.fieldNames.V1,
		a.fieldNames.V2,
		a.fieldNames.V3,
		a.fieldNames.V4,
		a.fieldNames.V5,
	}

	for i, filterValue := range filter {
		if filterValue != "" && i < len(columns) {
			query = query.Where(columns[i], filterValue)
		}
	}

	err := query.OrderAsc("id").Scan(&rules)
	if err != nil {
		return fmt.Errorf("gf-adapter-casbin3: failed to load filtered policies: %w", err)
	}

	for _, rule := range rules {
		loadPolicyLine(rule, model)
	}

	return nil
}
