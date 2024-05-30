package detectors

type Type string
type Language string

const (
	DetectorDependencies Type = "dependencies"
	DetectorBeego        Type = "beego"
	DetectorCSharp       Type = "csharp"
	DetectorDjango       Type = "django"
	DetectorDotnet       Type = "dotnet"
	DetectorEnvFile      Type = "env_file"
	DetectorGo           Type = "golang"
	DetectorJava         Type = "java"
	DetectorJavascript   Type = "javascript"
	DetectorTypescript   Type = "typescript"
	DetectorTsx          Type = "tsx"
	DetectorOpenAPI      Type = "openapi"
	DetectorRails        Type = "rails"
	DetectorRuby         Type = "ruby"
	DetectorPHP          Type = "php"
	DetectorPython       Type = "python"
	DetectorSimple       Type = "simple"
	DetectorSpring       Type = "spring"
	DetectorSymfony      Type = "symfony"
	DetectorYamlConfig   Type = "yaml_config"
	DetectorSQL          Type = "sql"
	DetectorProto        Type = "proto"
	DetectorGraphQL      Type = "graphql"
	DetectorHTML         Type = "html"
	DetectorIPYNB        Type = "ipynb"
	DetectorGitleaks     Type = "gitleaks"
	DetectorCustom       Type = "custom"
	DetectorSchemaRb     Type = "schema_rb"
)
