package rest

//Settings are the global settings for the rest trigger
type Settings struct {
	Port string `json:"port"`
}

type HandlerSettings struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}
