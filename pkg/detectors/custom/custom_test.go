package custom_test

import (
	_ "embed"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"gopkg.in/yaml.v2"

	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/detectors"
	"github.com/bearer/curio/pkg/detectors/custom"
	"github.com/bearer/curio/pkg/parser/nodeid"

	"github.com/bearer/curio/pkg/detectors/internal/testhelper"
	detectortypes "github.com/bearer/curio/pkg/report/detectors"
)

const detectorType = detectortypes.DetectorCustom

//go:embed testdata/config/ruby_loggers.yml
var configRubyLoggers []byte

//go:embed testdata/config/rails_sessions.yml
var configRailsSessions []byte

//go:embed testdata/config/ruby_password_length.yml
var configPasswordLength []byte

//go:embed testdata/config/rails_encrypts.yml
var configRailsEncrypts []byte

//go:embed testdata/config/sql_create_function.yml
var configSQLCreateFunction []byte

//go:embed testdata/config/sql_create_table.yml
var configSQLCreateTable []byte

//go:embed testdata/config/sql_create_table_mysql.yml
var configSQLCreateTableMySQL []byte

//go:embed testdata/config/sql_create_trigger.yml
var configSQLCreateTrigger []byte

//go:embed testdata/config/insecure_smtp.yml
var configInsecureSMTP []byte

//go:embed testdata/config/insecure_communication.yml
var configInsecureCommunication []byte

//go:embed testdata/config/insecure_ftp.yml
var configInsecureFTP []byte

//go:embed testdata/config/ruby_third_party_data_send.yml
var configRubyThirdPartyDataSend []byte

func TestRubyPasswordLengthClassMinMax(t *testing.T) {
	result := runTest(configPasswordLength, filepath.Join("testdata", "ruby", "password_length", "class_min_max"), t)
	cupaloy.SnapshotT(t, result)
}

func TestRubyPasswordLengthClass(t *testing.T) {
	result := runTest(configPasswordLength, filepath.Join("testdata", "ruby", "password_length", "class"), t)
	cupaloy.SnapshotT(t, result)
}

func TestRubyPasswordLengthIteration(t *testing.T) {
	result := runTest(configPasswordLength, filepath.Join("testdata", "ruby", "password_length", "iteration"), t)
	cupaloy.SnapshotT(t, result)
}

func TestRailsSessionsCallJSON(t *testing.T) {
	result := runTest(configRailsSessions, filepath.Join("testdata", "ruby", "sessions", "call"), t)
	cupaloy.SnapshotT(t, result)
}

func TestRailsSessionsIdentifierJSON(t *testing.T) {
	result := runTest(configRailsSessions, filepath.Join("testdata", "ruby", "sessions", "identifier"), t)
	cupaloy.SnapshotT(t, result)
}

func TestRailSessionHashJSON(t *testing.T) {
	result := runTest(configRailsSessions, filepath.Join("testdata", "ruby", "sessions", "hash"), t)
	cupaloy.SnapshotT(t, result)
}

func TestRailSessionHashAssigmentJSON(t *testing.T) {
	result := runTest(configRailsSessions, filepath.Join("testdata", "ruby", "sessions", "hash_assigment"), t)
	cupaloy.SnapshotT(t, result)
}

func TestRubyLoggersJSON(t *testing.T) {
	result := runTest(configRubyLoggers, filepath.Join("testdata", "ruby", "loggers"), t)
	cupaloy.SnapshotT(t, result)
}

func TestRubyLoggersVariableReconciliation(t *testing.T) {
	result := runTest(configRubyLoggers, filepath.Join("testdata", "ruby", "variable_reconciliation"), t)
	cupaloy.SnapshotT(t, result)
}

func TestRailsEncryptsJSON(t *testing.T) {
	result := runTest(configRailsEncrypts, filepath.Join("testdata", "ruby", "class", "encrypts"), t)
	cupaloy.SnapshotT(t, result)
}

func TestSQLCreateFunctionJSON(t *testing.T) {
	result := runTest(configSQLCreateFunction, filepath.Join("testdata", "sql", "create_function"), t)
	cupaloy.SnapshotT(t, result)
}
func TestSQLCreateTableJSON(t *testing.T) {
	result := runTest(configSQLCreateTable, filepath.Join("testdata", "sql", "create_table"), t)
	cupaloy.SnapshotT(t, result)
}
func TestSQLCreateTableMySQLJSON(t *testing.T) {
	result := runTest(configSQLCreateTableMySQL, filepath.Join("testdata", "sql", "create_table_mysql"), t)
	cupaloy.SnapshotT(t, result)
}

func TestSQLCreateTriggerJSON(t *testing.T) {
	result := runTest(configSQLCreateTrigger, filepath.Join("testdata", "sql", "create_trigger"), t)
	cupaloy.SnapshotT(t, result)
}

func TestInsecureSMTPJSON(t *testing.T) {
	result := runTest(configInsecureSMTP, filepath.Join("testdata", "ruby", "insecure_smtp"), t)
	cupaloy.SnapshotT(t, result)
}

func TestInsecureCommunicationJSON(t *testing.T) {
	result := runTest(configInsecureCommunication, filepath.Join("testdata", "ruby", "insecure_communication"), t)
	cupaloy.SnapshotT(t, result)
}

func TestInsecureFTPJSON(t *testing.T) {
	result := runTest(configInsecureFTP, filepath.Join("testdata", "ruby", "insecure_ftp"), t)
	cupaloy.SnapshotT(t, result)
}

func TestRubyThirdPartyDataSendJSON(t *testing.T) {
	result := runTest(configRubyThirdPartyDataSend, filepath.Join("testdata", "ruby", "third_party_data_send"), t)
	cupaloy.SnapshotT(t, result)
}

func runTest(config []byte, path string, t *testing.T) string {
	var rulesConfig map[string]settings.Rule

	detector := custom.New(&nodeid.IntGenerator{Counter: 0})
	err := yaml.Unmarshal(config, &rulesConfig)
	if err != nil {
		t.Fatal(err)
	}
	customDetector := detector.(*custom.Detector)
	err = customDetector.CompileRules(rulesConfig)
	if err != nil {
		t.Fatal(err)
	}

	var registrations = []detectors.InitializedDetector{{
		Type:     detectorType,
		Detector: detector}}
	detectorReport := testhelper.Extract(t, path, registrations, detectorType)

	bytes, err := json.MarshalIndent(detectorReport.CustomDetections, "", "\t")
	if err != nil {
		t.Fatal(err)
	}

	return string(bytes)
}
