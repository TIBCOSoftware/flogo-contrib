package fggos

import (
	"strings"

	"encoding/json"
	"strconv"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/japm/goScript"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"fmt"
)

// GosLinkExprManager is the Lua Implementation of a Link Expression Manager
type GosLinkExprManager struct {
	values map[int][]*varInfo
	exprs  map[int]*goScript.Expr
}

type varInfo struct {
	name string
	isd  int
}

type GosLinkExprManagerFactory struct {
}

// NewGosLinkExprManager creates a new LuaLinkExprManager
func (f *GosLinkExprManagerFactory) NewLinkExprManager(def *definition.Definition) definition.LinkExprManager {

	mgr := &GosLinkExprManager{}
	mgr.values = make(map[int][]*varInfo)
	mgr.exprs = make(map[int]*goScript.Expr)

	links := definition.GetExpressionLinks(def)

	for _, link := range links {

		if len(strings.TrimSpace(link.Value())) > 0 {
			vars, exprStr := transExpr(link.Value())

			mgr.values[link.ID()] = vars

			logger.Debugf("expr: %v\n", exprStr)

			expr := &goScript.Expr{}
			err := expr.Prepare(exprStr)

			if err == nil {
				mgr.exprs[link.ID()] = expr
			} else {
				logger.Errorf("Error preparing expression: %s - %v", link.Value(), err)
			}
		}
	}

	return mgr
}

func transExpr(s string) ([]*varInfo, string) {

	var vars []*varInfo
	var rvars []string

	strLen := len(s)

	isd := 0

	for i := 0; i < strLen; i++ {
		if s[i] == '$' {

			isdefcheck := false

			if strings.HasSuffix(s[0:i], "isDefined(") {
				isdefcheck = true
			}

			ignoreBraces := s[i+1] == '{'
			var partOfName bool

			var j int
			for j = i + 1; j < strLen; j++ {

				partOfName, ignoreBraces = isPartOfName(s[j], ignoreBraces)

				if !partOfName {
					break
				}
			}

			if isdefcheck {
				isd++
				vars = append(vars, &varInfo{isd: isd, name: s[i+1: j]})
				rvars = append(rvars, s[i-10:j+1])
				rvars = append(rvars, "isd"+strconv.Itoa(isd))
				i = j + 1
			} else {
				vars = append(vars, &varInfo{name: s[i+1: j]})
				rvars = append(rvars, s[i:j])
				rvars = append(rvars, `v["`+s[i+1:j]+`"]`)
				i = j
			}
		}
	}

	replacer := strings.NewReplacer(rvars...)

	return vars, replacer.Replace(s)
}

func isPartOfName(char byte, ignoreBraces bool) (bool, bool) {

	if (char < '0' || char > '9') && (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') && char != '.' && char != '_' {

		if ignoreBraces && char == '{' {
			return true, true
		} else if ignoreBraces && char == '}' {
			return true, false
		}

		return false, ignoreBraces

	}

	return true, ignoreBraces
}

// EvalLinkExpr implements LinkExprManager.EvalLinkExpr
func (em *GosLinkExprManager) EvalLinkExpr(link *definition.Link, scope data.Scope) (ret bool, err error) {

	defer func() {
		if r := recover(); r != nil {
			ret = false
			err = definition.NewLinkExprError(fmt.Sprintf("Error evaluting expression: %s", r))
		}
	}()

	if link.Type() == definition.LtDependency {
		// dependency links are always true
		return true, nil
	}

	vars, attrsOK := em.values[link.ID()]
	expr, exprOK := em.exprs[link.ID()]

	if !attrsOK || !exprOK {

		return false, fmt.Errorf("Unable to evaluate expression '%s', did not compile properly\n", link.Value())
	}

	ctxt := make(map[string]interface{})
	vals := make(map[string]interface{})

	for _, varInfo := range vars {

		var attrValue interface{}

		attrName, attrPath, _ := data.GetAttrPath(varInfo.name)
		attr, exists := scope.GetAttr(attrName)

		attrValue = attr.Value

		if varInfo.isd > 0 {

			if exists && len(attrPath) > 0 {

				//for now assume if we have a path, attr is "object" and only one level
				valMap, ok := attrValue.(map[string]interface{})

				if ok {
					_, exists = valMap[attrPath]
				} else {
					//assume its a map[string]string
					strMap, ok := attrValue.(map[string]string)

					if ok {
						_, exists = strMap[attrPath]
					}
				}
				//todo what if the value does not exists
			}

			ctxt["isd"+strconv.Itoa(varInfo.isd)] = exists

		} else {

			if exists && len(attrPath) > 0 {

				valMap, ok := attrValue.(map[string]interface{})

				var val interface{}

				if ok {
					val = data.GetMapValue(valMap, attrPath)
				} else {
					//assume its a map[string]string
					strMap, ok := attrValue.(map[string]string)
					if ok {
						val = strMap[attrPath]
					}
				}

				attrValue = FixUpValue(val)
			}

			vals[varInfo.name] = attrValue
		}
	}

	ctxt["v"] = vals

	logger.Debugf("Vals: %v", vals)

	val, err := expr.Eval(ctxt)

	if err != nil {
		return false, definition.NewLinkExprError(fmt.Sprintf("Error evaluting expression: %s", err.Error()))
	}

	return val.(bool), nil
}

// FixUpValue fixes json numbers
func FixUpValue(val interface{}) interface{} {

	ret := val
	var err error

	switch t := val.(type) {
	case json.Number:
		if strings.Index(t.String(), ".") > -1 {
			ret, err = t.Float64()
		} else {
			ret, err = t.Int64()
		}
	}

	if err != nil {
		ret = val
	}

	return ret
}
