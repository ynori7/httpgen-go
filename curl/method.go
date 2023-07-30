package curl

const methodHeader = "-X"

func GetMethod(params map[string][]string) string {
	method := ""
	if len(params[methodHeader]) == 1 {
		method = params[methodHeader][0]
	}

	// infer default method
	if method == "" {
		if _, ok := GetBody(params); ok {
			method = "POST"
		} else {
			method = "GET"
		}
	}
	return method
}
