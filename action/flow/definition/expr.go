package definition

import "github.com/TIBCOSoftware/flogo-lib/core/data"

// LinkExprManager interface that defines a Link Expression Manager
type LinkExprManager interface {

	// EvalLinkExpr evaluate the link expression
	EvalLinkExpr(link *Link, scope data.Scope) (bool, error)
}

type LinkExprManagerFactory interface {

	NewLinkExprManager(def *Definition) LinkExprManager
}

var	linkExprMangerFactory LinkExprManagerFactory

func SetLinkExprManagerFactory(factory LinkExprManagerFactory) {
	linkExprMangerFactory = factory
}

func GetLinkExprManagerFactory() LinkExprManagerFactory {
	return linkExprMangerFactory
}

// GetExpressionLinks gets the links of the definition that are of type LtExpression
func GetExpressionLinks(def *Definition) []*Link {

	var links []*Link

	getExpressionLinks(def.RootTask(), &links)

	if (def.ErrorHandlerTask() != nil) {
		getExpressionLinks(def.ErrorHandlerTask(), &links)
	}

	return links
}

// getExpressionLinks gets the links under the specified task that are of type LtExpression
func getExpressionLinks(task *Task, links *[]*Link) {

	for _, link := range task.ChildLinks() {

		if link.Type() == LtExpression {
			*links = append(*links, link)
		}
	}

	for _, childTask := range task.ChildTasks() {
		getExpressionLinks(childTask, links)
	}
}
