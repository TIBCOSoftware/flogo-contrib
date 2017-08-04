package definition

import (
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/util"
	"strconv"
)

// DefinitionRep is a serializable representation of a flow Definition
type DefinitionRep struct {
	ExplicitReply    bool               `json:"explicitReply"`
	Name             string             `json:"name"`
	ModelID          string             `json:"model"`
	Attributes       []*data.Attribute  `json:"attributes,omitempty"`
	InputMappings    []*data.MappingDef `json:"inputMappings,omitempty"`
	RootTask         *TaskRep           `json:"rootTask"`
	ErrorHandlerTask *TaskRep           `json:"errorHandlerTask"`
}

// TaskRep is a serializable representation of a flow Task
type TaskRep struct {
	// Using interface{} type to support backward compatibility changes since Id was
	// int before, change to string once BC is removed
	ID           interface{}       `json:"id"`
	TypeID       int               `json:"type"`
	ActivityType string            `json:"activityType"`
	ActivityRef  string            `json:"activityRef"`
	Name         string            `json:"name"`
	Attributes   []*data.Attribute `json:"attributes,omitempty"`

	InputAttrs  map[string]interface{} `json:"inputs,omitempty"`
	OutputAttrs map[string]interface{} `json:"outputs,omitempty"`

	InputMappings  []*data.MappingDef `json:"inputMappings,omitempty"`
	OutputMappings []*data.MappingDef `json:"ouputMappings,omitempty"`

	Tasks []*TaskRep `json:"tasks,omitempty"`
	Links []*LinkRep `json:"links,omitempty"`
}

// LinkRep is a serializable representation of a flow Link
type LinkRep struct {
	ID   int    `json:"id"`
	Type int    `json:"type"`
	Name string `json:"name"`
	// Using interface{} type to support backward compatibility changes since Id was
	// int before, change to string once BC is removed
	ToID interface{} `json:"to"`
	// Using interface{} type to support backward compatibility changes since Id was
	// int before, change to string once BC is removed
	FromID interface{} `json:"from"`
	Value  string      `json:"value"`
}

// NewDefinition creates a flow Definition from a serializable
// definition representation
func NewDefinition(rep *DefinitionRep) (def *Definition, err error) {

	defer util.HandlePanic("NewDefinition", &err)

	def = &Definition{}
	def.name = rep.Name
	def.modelID = rep.ModelID
	def.explicitReply = rep.ExplicitReply

	//todo is this used or needed?
	if rep.InputMappings != nil {
		def.inputMapper = GetMapperFactory().NewMapper(&MapperDef{Mappings: rep.InputMappings})
	}

	if len(rep.Attributes) > 0 {
		def.attrs = make(map[string]*data.Attribute, len(rep.Attributes))

		for _, value := range rep.Attributes {
			def.attrs[value.Name] = value
		}
	}

	def.rootTask = &Task{}

	def.tasks = make(map[string]*Task)
	def.links = make(map[int]*Link)

	addTask(def, def.rootTask, rep.RootTask)
	addLinks(def, def.rootTask, rep.RootTask)

	if rep.ErrorHandlerTask != nil {
		def.ehTask = &Task{}

		addTask(def, def.ehTask, rep.ErrorHandlerTask)
		addLinks(def, def.ehTask, rep.ErrorHandlerTask)
	}

	return def, nil
}

func addTask(def *Definition, task *Task, rep *TaskRep) {
	// Workaround to support Backwards compatibility
	// Remove once rep.ID is string
	task.id = convertInterfaceToString(rep.ID)
	task.activityType = rep.ActivityType
	task.activityRef = rep.ActivityRef
	task.typeID = rep.TypeID
	task.name = rep.Name
	//task.Definition = def

	if rep.InputMappings != nil {
		task.inputMapper = GetMapperFactory().NewTaskInputMapper(task, &MapperDef{Mappings: rep.InputMappings})
	}

	if rep.OutputMappings != nil {
		task.outputMapper = GetMapperFactory().NewTaskOutputMapper(task, &MapperDef{Mappings: rep.OutputMappings})
	}

	if task.outputMapper == nil {
		task.outputMapper = GetMapperFactory().GetDefaultTaskOutputMapper(task)
	}

	// Keep for now, DEPRECATE "attributes" section from flogo.json
	if len(rep.Attributes) > 0 {
		task.inputAttrs = make(map[string]*data.Attribute, len(rep.Attributes))

		for _, value := range rep.Attributes {
			task.inputAttrs[value.Name] = value
		}
	}

	act := activity.Get(task.activityRef)

	if act != nil {

		if len(rep.InputAttrs) > 0 {
			task.inputAttrs = make(map[string]*data.Attribute, len(rep.InputAttrs))

			for name, value := range rep.InputAttrs {

				attr := act.Metadata().Inputs[name]

				if attr != nil {
					newValue, err := data.CoerceToValue(value, attr.Type)
					if err != nil {
						//Todo handle error
						newValue = value
					}
					task.inputAttrs[name] = &data.Attribute{Name: name, Type: attr.Type, Value: newValue}
				}
			}
		}

		if len(rep.OutputAttrs) > 0 {

			task.outputAttrs = make(map[string]*data.Attribute, len(rep.OutputAttrs))

			for name, value := range rep.OutputAttrs {

				attr := act.Metadata().Outputs[name]

				if attr != nil {
					newValue, err := data.CoerceToValue(value, attr.Type)
					if err != nil {
						//Todo handle error
						newValue = value
					}
					task.outputAttrs[name] = &data.Attribute{Name: name, Type: attr.Type, Value: newValue}
				}
			}
		}
	}

	def.tasks[task.id] = task
	numTasks := len(rep.Tasks)

	// flow child tasks
	if numTasks > 0 {

		for _, childTaskRep := range rep.Tasks {

			childTask := &Task{}
			childTask.parent = task
			task.tasks = append(task.tasks, childTask)
			addTask(def, childTask, childTaskRep)
		}
	}
}

//convertInterfaceToString will identify whether the interface is int or string and return a string in any case
func convertInterfaceToString(m interface{}) string {
	if m == nil {
		panic("Invalid nil activity id found")
	}
	switch m.(type) {
	case string:
		return m.(string)
	case float64:
		return strconv.Itoa(int(m.(float64)))
	default:
		panic(fmt.Sprintf("Error parsing Task with Id '%v', invalid type '%T'", m, m))
	}
}

func addLinks(def *Definition, task *Task, rep *TaskRep) {

	numLinks := len(rep.Links)

	if numLinks > 0 {

		task.links = make([]*Link, numLinks)

		for i, linkRep := range rep.Links {

			link := &Link{}
			link.id = linkRep.ID
			//link.Parent = task
			//link.Definition = pd
			link.linkType = LinkType(linkRep.Type)
			link.value = linkRep.Value
			link.fromTask = def.tasks[convertInterfaceToString(linkRep.FromID)]
			link.toTask = def.tasks[convertInterfaceToString(linkRep.ToID)]

			// add this link as predecessor "fromLink" to the "toTask"
			link.toTask.fromLinks = append(link.toTask.fromLinks, link)

			// add this link as successor "toLink" to the "fromTask"
			link.fromTask.toLinks = append(link.fromTask.toLinks, link)

			task.links[i] = link
			def.links[link.id] = link
		}
	}
}
