package model

type Metadata struct {
	Name   string
	Inputs struct {
		Params   map[string]OperationParam
		Features map[string]Feature
	}
	Outputs map[string]OperationParam
	Method  string
	Tag     string
	SigDef  string
}

type Feature struct {
	Shape []int64
	Type  string
}

type OperationParam struct {
	Name  string
	Type  string
	Shape []int64
}
