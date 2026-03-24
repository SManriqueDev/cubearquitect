package orchestrator

func getStringParam(params map[string]string, key, defaultValue string) string {
	if v, ok := params[key]; ok && v != "" {
		return v
	}
	return defaultValue
}

func getBoolParam(params map[string]string, key string, defaultValue bool) bool {
	if v, ok := params[key]; ok {
		return v == "true"
	}
	return defaultValue
}

func hasKey(params map[string]string, key string) bool {
	_, ok := params[key]
	return ok
}
