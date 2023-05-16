package flogocontrib

import _ "embed"

//go:embed VERSION
var version string

// Version will return the release version
func Version() string {
	return version
}
