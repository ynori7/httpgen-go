package curl

const methodHeader = "-X"

func GetMethod(params map[string][]string) string {
	method := "GET"
	if len(params[methodHeader]) == 1 {
		method = params[methodHeader][0]
	}
	return method
}
