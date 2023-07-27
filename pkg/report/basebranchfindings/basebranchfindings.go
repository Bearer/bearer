package basebranchfindings

type key struct {
	RuleID   string
	Filename string
}

type lineRange struct {
	StartLine,
	EndLine int
}

type Findings map[key][]lineRange

func (findings Findings) Add(RuleID string, Filename string, StartLine, EndLine int) {
	key := key{
		RuleID:   RuleID,
		Filename: Filename,
	}

	findings[key] = append(findings[key], lineRange{
		StartLine: StartLine,
		EndLine:   EndLine,
	})
}

func (findings Findings) Has(RuleID string, Filename string, StartLine, EndLine int) bool {
	key := key{
		RuleID:   RuleID,
		Filename: Filename,
	}

	for _, finding := range findings[key] {
		if finding.StartLine <= EndLine && finding.EndLine >= StartLine {
			return true
		}
	}

	return false
}
