package jsonlinesreader

func ReadAllText(path string) ([]string, error) {
	var data []string
	reader, err := New(path)
	if err != nil {
		return nil, err
	}

	for reader.Next() {
		data = append(data, reader.Text())
	}

	return data, nil
}
