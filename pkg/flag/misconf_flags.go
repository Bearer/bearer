package flag

// e.g. config yaml:
//
//	misconfiguration:
//	  trace: true
//	  config-policy: "custom-policy/policy"
//	  policy-namespaces: "user"
var (
	ConfigPolicyFlag = Flag{
		Name:       "config-policy",
		ConfigName: "misconfiguration.policy",
		Value:      []string{},
		Usage:      "specify paths to the Rego policy files directory, applying config files",
	}
	ConfigDataFlag = Flag{
		Name:       "config-data",
		ConfigName: "misconfiguration.data",
		Value:      []string{},
		Usage:      "specify paths from which data for the Rego policies will be recursively loaded",
	}
	PolicyNamespaceFlag = Flag{
		Name:       "policy-namespaces",
		ConfigName: "misconfiguration.namespaces",
		Value:      []string{},
		Usage:      "rego namespaces",
	}
)

// MisconfFlagGroup composes common printer flag structs used for commands providing misconfinguration scanning.
type MisconfFlagGroup struct {
	// Rego
	PolicyPaths      *Flag
	DataPaths        *Flag
	PolicyNamespaces *Flag
}

type MisconfOptions struct {
	// Rego
	PolicyPaths      []string
	DataPaths        []string
	PolicyNamespaces []string
}

func NewMisconfFlagGroup() *MisconfFlagGroup {
	return &MisconfFlagGroup{
		PolicyPaths:      &ConfigPolicyFlag,
		DataPaths:        &ConfigDataFlag,
		PolicyNamespaces: &PolicyNamespaceFlag,
	}
}

func (f *MisconfFlagGroup) Name() string {
	return "Misconfiguration"
}

func (f *MisconfFlagGroup) Flags() []*Flag {
	return []*Flag{
		f.PolicyPaths,
		f.DataPaths,
		f.PolicyNamespaces,
	}
}

func (f *MisconfFlagGroup) ToOptions() (MisconfOptions, error) {
	return MisconfOptions{
		PolicyPaths:      getStringSlice(f.PolicyPaths),
		DataPaths:        getStringSlice(f.DataPaths),
		PolicyNamespaces: getStringSlice(f.PolicyNamespaces),
	}, nil
}
