package views

func Lookup(name string) View { return cache[name] }

func Available() (names []string) {
	for name := range cache {
		names = append(names, name)
	}

	if names == nil {
		names = make([]string, 0, 0)
	}

	return names
}

func Render(name string, data interface{}) ([]byte, error) {
	v := Lookup(name)
	if v == nil {
		return []byte{}, nil
	}

	buf := getbuf()
	defer putbuf(buf)

	err := v.Execute(buf, data)

	return buf.Bytes(), err
}
