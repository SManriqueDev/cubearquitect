package orchestrator

func getStringParam(params map[string]string, key, defaultValue string) string {
	if v, ok := params[key]; ok && v != "" {
		return v
	}
	return defaultValue
}

func getBoolParam(params map[string]string, key string, _ bool) bool {
	if v, ok := params[key]; ok {
		return v == paramTrue
	}
	return false
}
