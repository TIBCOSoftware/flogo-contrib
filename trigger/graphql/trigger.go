package graphql

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/TIBCOSoftware/flogo-contrib/trigger/rest/cors"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/graphql-go/graphql"
	"github.com/julienschmidt/httprouter"

	"net/http"
)

const (
	REST_CORS_PREFIX = "GRAPHQL_TRIGGER"
)

// log is the default package logger
var log = logger.GetLogger("trigger-flogo-graphql")

var gqlObjects map[string]*graphql.Object
var graphQlSchema *graphql.Schema

// GraphQLTrigger trigger struct
type GraphQLTrigger struct {
	metadata *trigger.Metadata
	server   *Server
	config   *trigger.Config
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &GraphQLFactory{metadata: md}
}

// GraphQLFactory Trigger factory
type GraphQLFactory struct {
	metadata *trigger.Metadata
}

//New Creates a new trigger instance for a given id
func (t *GraphQLFactory) New(config *trigger.Config) trigger.Trigger {
	return &GraphQLTrigger{metadata: t.metadata, config: config}
}

// Metadata implements trigger.Trigger.Metadata
func (t *GraphQLTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

func (t *GraphQLTrigger) Initialize(ctx trigger.InitContext) error {
	router := httprouter.New()

	if t.config.Settings == nil {
		return fmt.Errorf("no Settings found for trigger '%s'", t.config.Id)
	}

	if _, ok := t.config.Settings["port"]; !ok {
		return fmt.Errorf("no Port found for trigger '%s' in settings", t.config.Id)
	}

	addr := ":" + t.config.GetSetting("port")

	// Build the GraphQL Object Types & Schemas
	t.buildGraphQLObjects()
	graphQlSchema = t.buildGraphQLSchema(ctx.GetHandlers())

	// Setup routes for the path & verb
	router.Handle("GET", t.config.GetSetting("path"), newActionHandler(t))
	router.Handle("POST", t.config.GetSetting("path"), newActionHandler(t))

	log.Debugf("Configured on port %s", t.config.Settings["port"])
	t.server = NewServer(addr, router)

	return nil
}

func (t *GraphQLTrigger) buildGraphQLObjects() {
	gqlTypes := t.config.Settings["types"].([]interface{})

	// Create type objects
	gqlObjects = make(map[string]*graphql.Object)

	// Get the graphql types
	for _, typ := range gqlTypes {
		lTyp := lower(typ)
		typ := lTyp.(map[string]interface{})
		name := typ["name"].(string)
		fields := make(graphql.Fields)

		for k, f := range typ["fields"].(map[string]interface{}) {
			fTyp := f.(map[string]interface{})

			fields[k] = &graphql.Field{
				Type: coerceType(fTyp["type"].(string)),
			}
		}

		obj := graphql.NewObject(
			graphql.ObjectConfig{
				Name:   name,
				Fields: fields,
			})

		gqlObjects[name] = obj
	}
}

func (t *GraphQLTrigger) buildGraphQLSchema(handlers []*trigger.Handler) *graphql.Schema {
	fSchema := t.config.Settings["schema"].(map[string]interface{})
	fSchema = lower(fSchema).(map[string]interface{})

	// Build the graphql schema
	var schema graphql.Schema
	var queryType *graphql.Object

	if strings.EqualFold(t.config.Settings["operation"].(string), "query") {

		var objName string
		queryFields := make(graphql.Fields)

		// Get the object name
		for k, v := range fSchema["query"].(map[string]interface{}) {
			if strings.EqualFold(k, "name") {
				objName = v.(string)
			} else if strings.EqualFold(k, "fields") {
				qf := v.(map[string]interface{})

				for k, v := range qf {

					// Grab query args
					argObj := v.(map[string]interface{})
					args := make(graphql.FieldConfigArgument)

					for k, v := range argObj["args"].(map[string]interface{}) {

						argTyp := v.(map[string]interface{})
						args[k] = &graphql.ArgumentConfig{
							Type: coerceType(argTyp["type"].(string)),
						}
					}

					for _, handler := range handlers {
						if strings.EqualFold(handler.GetStringSetting("resolverFor"), k) {
							// Build the queryField
							queryFields[k] = &graphql.Field{
								Type:    gqlObjects[k],
								Args:    args,
								Resolve: fieldResolver(handler),
							}
						}
					}
				}
			}
		}

		queryType = graphql.NewObject(
			graphql.ObjectConfig{
				Name:   objName,
				Fields: queryFields,
			})
	}

	schema, _ = graphql.NewSchema(
		graphql.SchemaConfig{
			Query: queryType,
		})

	return &schema
}

func (t *GraphQLTrigger) Start() error {
	return t.server.Start()
}

// Stop implements util.Managed.Stop
func (t *GraphQLTrigger) Stop() error {
	return t.server.Stop()
}

func fieldResolver(handler *trigger.Handler) graphql.FieldResolveFn {

	return func(p graphql.ResolveParams) (interface{}, error) {

		triggerData := map[string]interface{}{
			"args": p.Args,
		}

		results, err := handler.Handle(context.Background(), triggerData)
		return results["data"].Value(), err
	}

}

// Handles the cors preflight request
func handleCorsPreflight(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	log.Infof("Received [OPTIONS] request to CorsPreFlight: %+v", r)

	c := cors.New(REST_CORS_PREFIX, log)
	c.HandlePreflight(w, r)
}

func newActionHandler(rt *GraphQLTrigger) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		log.Infof("Received request for id '%s'", rt.config.Id)

		c := cors.New(REST_CORS_PREFIX, log)
		c.WriteCorsActualRequestHeaders(w)

		queryValues := r.URL.Query()
		queryParams := make(map[string]string, len(queryValues))
		header := make(map[string]string, len(r.Header))

		for key, value := range r.Header {
			header[key] = strings.Join(value, ",")
		}

		for key, value := range queryValues {
			queryParams[key] = strings.Join(value, ",")
		}

		var query string

		httpVerb := strings.ToUpper(r.Method)
		if val, ok := queryParams["query"]; ok && strings.EqualFold(httpVerb, "GET") {
			query = val
		} else if strings.EqualFold(httpVerb, "POST") {
			// Check the HTTP Header Content-Type
			contentType := r.Header.Get("Content-Type")
			if !strings.EqualFold(contentType, "application/json") {
				err := fmt.Errorf("%v", "Invalid content type. Must be application/json for POST methods.")
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			var content interface{}
			err := json.NewDecoder(r.Body).Decode(&content)
			if err != nil {
				switch {
				case err == io.EOF:
					// empty body
					//todo should handler say if content is expected?
				case err != nil:
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}
			jsonContent := content.(map[string]interface{})
			query = jsonContent["query"].(string)
		} else {
			err := fmt.Errorf("%v", "HTTP GET and POST are the only supported verbs.")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Process the request
		result := graphql.Do(graphql.Params{
			Schema:        *graphQlSchema,
			RequestString: query,
		})

		if len(result.Errors) > 0 {
			log.Errorf("GraphQL Trigger Error: %#v", result.Errors)
		}

		if result != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)

			if err := json.NewEncoder(w).Encode(result); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Error(err)
			}
		}
	}
}

////////////////////////////////////////////////////////////////////////////////////////
// Utils

func coerceType(typ string) *graphql.Scalar {
	switch typ {
	case "graphql.String":
		return graphql.String
	case "graphql.Float":
		return graphql.Float
	case "graphql.Int":
		return graphql.Int
	case "graphql.Boolean":
		return graphql.Boolean
	}

	return nil
}

func lower(f interface{}) interface{} {
	switch f := f.(type) {
	case []interface{}:
		for i := range f {
			f[i] = lower(f[i])
		}
		return f
	case map[string]interface{}:
		lf := make(map[string]interface{}, len(f))
		for k, v := range f {
			lf[strings.ToLower(k)] = lower(v)
		}
		return lf
	default:
		return f
	}
}
