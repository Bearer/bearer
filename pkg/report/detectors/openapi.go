package detectors

type OpenAPIFileType string

const (
	OpenApi2JSONFile OpenAPIFileType = "2_json"
	OpenApi2YAMLFile OpenAPIFileType = "2_yaml"
	OpenApi3JSONFile OpenAPIFileType = "3_json"
	OpenApi3YAMLFile OpenAPIFileType = "3_yaml"
)
