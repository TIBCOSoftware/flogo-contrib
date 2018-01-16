package util

var aliasToInt = make(map[string]int)

func RegisterIntAlias(alias string, intVal int) {
	aliasToInt[alias] = intVal
}

func GetIntFromAlias(alias string) (int, bool) {
	intVal, found := aliasToInt[alias]
	return intVal, found
}
