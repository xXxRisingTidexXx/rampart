package mining

type number string

func (n *number) UnmarshalJSON(bytes []byte) error {
	s := string(bytes)
	if s[0] == '"' {
		s = s[1:]
	}
	if i := len(s) - 1; s[i] == '"' {
		s = s[:i]
	}
	*n = number(s)
	return nil
}
