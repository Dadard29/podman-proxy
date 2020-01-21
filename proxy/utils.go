package proxy

func checkMethod(methods []string, requestMethod string) bool {
	for _, v := range methods {
		if v == requestMethod {
			return true
		}
	}
	return false
}
