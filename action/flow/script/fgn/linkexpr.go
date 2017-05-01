package fgn

import (
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
)

type ExprFunc func(map[string]interface{}) bool

type GoExpr struct {
	attrs    []string
	exprFunc ExprFunc
}

type GoLinkExprManager struct {
	exprs  map[int]ExprFunc
	values map[int][]string
}

func NewGoLinkExprManager(exprs map[int]ExprFunc) *GoLinkExprManager {

	mgr := &GoLinkExprManager{exprs: exprs}
	return mgr
}

// EvalLinkExpr implements LinkExprManager.EvalLinkExpr
func (em *GoLinkExprManager) EvalLinkExpr(link *definition.Link, scope data.Scope) bool {

	if link.Type() == definition.LtDependency {
		// dependency links are always true
		return true
	}

	attrs, ok := em.values[link.ID()]

	if !ok {
		return false
	}

	vals := make(map[string]interface{})

	for _, origAttrName := range attrs {

		var attrValue interface{}

		attrName, attrPath, _ := data.GetAttrPath(origAttrName)

		attr, exists := scope.GetAttr(attrName)
		attrValue = attr.Value

		if exists && len(attrPath) > 0 {
			//for now assume if we have a path, attr is "object" and only one level
			valMap := attrValue.(map[string]interface{})
			//todo what if the value does not exists
			attrValue, _ = valMap[attrPath]
		}

		vals[origAttrName] = attrValue
	}

	expr := em.exprs[link.ID()]
	return expr(vals)
}
