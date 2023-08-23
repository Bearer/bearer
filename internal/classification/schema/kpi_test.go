package schema_test

import (
	"testing"

	"github.com/bearer/bearer/internal/classification/db"
	"github.com/bearer/bearer/internal/classification/schema"
	"github.com/bearer/bearer/internal/classification/schema/internal/testhelper"
	"github.com/bradleyjkemp/cupaloy"
)

func TestRubyKPI(t *testing.T) {
	classifier := schema.New(
		schema.Config{
			DataTypes:                      db.Default().DataTypes,
			DataTypeClassificationPatterns: db.Default().DataTypeClassificationPatterns,
			KnownPersonObjectPatterns:      db.Default().KnownPersonObjectPatterns,
		},
	)
	output := testhelper.ExtractExpectedOutput(t, "ruby", classifier)

	cupaloy.SnapshotT(t, output.KPI)
}
func TestRuby(t *testing.T) {
	classifier := schema.New(
		schema.Config{
			DataTypes:                      db.Default().DataTypes,
			DataTypeClassificationPatterns: db.Default().DataTypeClassificationPatterns,
			KnownPersonObjectPatterns:      db.Default().KnownPersonObjectPatterns,
		},
	)
	output := testhelper.ExtractExpectedOutput(t, "ruby", classifier)

	cupaloy.SnapshotT(t, output.ValidClassifications)
}

func TestCSharpKPI(t *testing.T) {
	classifier := schema.New(
		schema.Config{
			DataTypes:                      db.Default().DataTypes,
			DataTypeClassificationPatterns: db.Default().DataTypeClassificationPatterns,
			KnownPersonObjectPatterns:      db.Default().KnownPersonObjectPatterns,
		},
	)
	output := testhelper.ExtractExpectedOutput(t, "csharp", classifier)

	cupaloy.SnapshotT(t, output.KPI)
}

func TestCSharp(t *testing.T) {
	classifier := schema.New(
		schema.Config{
			DataTypes:                      db.Default().DataTypes,
			DataTypeClassificationPatterns: db.Default().DataTypeClassificationPatterns,
			KnownPersonObjectPatterns:      db.Default().KnownPersonObjectPatterns,
		},
	)
	output := testhelper.ExtractExpectedOutput(t, "csharp", classifier)

	cupaloy.SnapshotT(t, output.ValidClassifications)
}

// The following test are ignored because the classification output does
// not match the expectations set by the previous classification (Rails app)

// Include output.ExpectedClassifications and compare to output.ValidClassifications
// to investigate discrepancies

func TestGoKPI(t *testing.T) {
	classifier := schema.New(
		schema.Config{
			DataTypes:                      db.Default().DataTypes,
			DataTypeClassificationPatterns: db.Default().DataTypeClassificationPatterns,
			KnownPersonObjectPatterns:      db.Default().KnownPersonObjectPatterns,
		},
	)
	output := testhelper.ExtractExpectedOutput(t, "go", classifier)

	cupaloy.SnapshotT(t, output.KPI)
}

func TestGo(t *testing.T) {
	classifier := schema.New(
		schema.Config{
			DataTypes:                      db.Default().DataTypes,
			DataTypeClassificationPatterns: db.Default().DataTypeClassificationPatterns,
			KnownPersonObjectPatterns:      db.Default().KnownPersonObjectPatterns,
		},
	)
	output := testhelper.ExtractExpectedOutput(t, "go", classifier)

	cupaloy.SnapshotT(t, output.ValidClassifications)
}

func TestJavaKPI(t *testing.T) {
	classifier := schema.New(
		schema.Config{
			DataTypes:                      db.Default().DataTypes,
			DataTypeClassificationPatterns: db.Default().DataTypeClassificationPatterns,
			KnownPersonObjectPatterns:      db.Default().KnownPersonObjectPatterns,
		},
	)
	output := testhelper.ExtractExpectedOutput(t, "java", classifier)

	cupaloy.SnapshotT(t, output.KPI)
}

func TestJava(t *testing.T) {
	classifier := schema.New(
		schema.Config{
			DataTypes:                      db.Default().DataTypes,
			DataTypeClassificationPatterns: db.Default().DataTypeClassificationPatterns,
			KnownPersonObjectPatterns:      db.Default().KnownPersonObjectPatterns,
		},
	)
	output := testhelper.ExtractExpectedOutput(t, "java", classifier)

	cupaloy.SnapshotT(t, output.ValidClassifications)
}

func TestJavascriptKPI(t *testing.T) {
	classifier := schema.New(
		schema.Config{
			DataTypes:                      db.Default().DataTypes,
			DataTypeClassificationPatterns: db.Default().DataTypeClassificationPatterns,
			KnownPersonObjectPatterns:      db.Default().KnownPersonObjectPatterns,
		},
	)
	output := testhelper.ExtractExpectedOutput(t, "javascript", classifier)

	cupaloy.SnapshotT(t, output.KPI)
}

func TestJavascript(t *testing.T) {
	classifier := schema.New(
		schema.Config{
			DataTypes:                      db.Default().DataTypes,
			DataTypeClassificationPatterns: db.Default().DataTypeClassificationPatterns,
			KnownPersonObjectPatterns:      db.Default().KnownPersonObjectPatterns,
		},
	)
	output := testhelper.ExtractExpectedOutput(t, "javascript", classifier)

	cupaloy.SnapshotT(t, output.ValidClassifications)
}

func TestTypescriptKPI(t *testing.T) {
	classifier := schema.New(
		schema.Config{
			DataTypes:                      db.Default().DataTypes,
			DataTypeClassificationPatterns: db.Default().DataTypeClassificationPatterns,
			KnownPersonObjectPatterns:      db.Default().KnownPersonObjectPatterns,
		},
	)
	output := testhelper.ExtractExpectedOutput(t, "typescript", classifier)

	cupaloy.SnapshotT(t, output.KPI)
}

func TestTypescript(t *testing.T) {
	classifier := schema.New(
		schema.Config{
			DataTypes:                      db.Default().DataTypes,
			DataTypeClassificationPatterns: db.Default().DataTypeClassificationPatterns,
			KnownPersonObjectPatterns:      db.Default().KnownPersonObjectPatterns,
		},
	)
	output := testhelper.ExtractExpectedOutput(t, "typescript", classifier)

	cupaloy.SnapshotT(t, output.ValidClassifications)
}

func TestPHPKPI(t *testing.T) {
	classifier := schema.New(
		schema.Config{
			DataTypes:                      db.Default().DataTypes,
			DataTypeClassificationPatterns: db.Default().DataTypeClassificationPatterns,
			KnownPersonObjectPatterns:      db.Default().KnownPersonObjectPatterns,
		},
	)
	output := testhelper.ExtractExpectedOutput(t, "php", classifier)

	cupaloy.SnapshotT(t, output.KPI)
}

func TestPHP(t *testing.T) {
	classifier := schema.New(
		schema.Config{
			DataTypes:                      db.Default().DataTypes,
			DataTypeClassificationPatterns: db.Default().DataTypeClassificationPatterns,
			KnownPersonObjectPatterns:      db.Default().KnownPersonObjectPatterns,
		},
	)
	output := testhelper.ExtractExpectedOutput(t, "php", classifier)

	cupaloy.SnapshotT(t, output.ValidClassifications)
}

func TestPythonKPI(t *testing.T) {
	classifier := schema.New(
		schema.Config{
			DataTypes:                      db.Default().DataTypes,
			DataTypeClassificationPatterns: db.Default().DataTypeClassificationPatterns,
			KnownPersonObjectPatterns:      db.Default().KnownPersonObjectPatterns,
		},
	)
	output := testhelper.ExtractExpectedOutput(t, "python", classifier)

	cupaloy.SnapshotT(t, output.KPI)
}

func TestPython(t *testing.T) {
	classifier := schema.New(
		schema.Config{
			DataTypes:                      db.Default().DataTypes,
			DataTypeClassificationPatterns: db.Default().DataTypeClassificationPatterns,
			KnownPersonObjectPatterns:      db.Default().KnownPersonObjectPatterns,
		},
	)
	output := testhelper.ExtractExpectedOutput(t, "python", classifier)

	cupaloy.SnapshotT(t, output.ValidClassifications)
}
