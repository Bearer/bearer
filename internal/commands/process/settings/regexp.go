package settings

import "regexp"

type Regexp struct {
	*regexp.Regexp
}

func (r *Regexp) UnmarshalText(b []byte) error {
	pattern, err := regexp.Compile(string(b))
	if err != nil {
		return err
	}

	r.Regexp = pattern

	return nil
}

func (r *Regexp) MarshalText() ([]byte, error) {
	if r.Regexp != nil {
		return []byte(r.Regexp.String()), nil
	}

	return nil, nil
}
