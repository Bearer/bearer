package types

type Stats struct {
	NumberOfLines        int32            `json:"number_of_lines" yaml:"number_of_lines"`
	NumberOfDataTypes    int              `json:"number_of_data_types" yaml:"number_of_data_types"`
	DataTypes            []DataType       `json:"data_types" yaml:"data_types"`
	NumberOfDatabases    int              `json:"-" yaml:"-"`
	NumberOfExternalAPIs int              `json:"-" yaml:"-"`
	NumberOfInternalAPIs int              `json:"-" yaml:"-"`
	Languages            map[string]int32 `json:"-" yaml:"-"`
	DataGroups           []string         `json:"-" yaml:"-"`
}

type DataType struct {
	Name         string `json:"name" yaml:"name"`
	CategoryUUID string `json:"-" yaml:"-"`
	Encrypted    bool   `json:"-" yaml:"-"`
	Occurrences  int    `json:"occurrences" yaml:"occurrences"`
}
