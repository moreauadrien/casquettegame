package structures

type Set map[string]struct{}

func (s Set) Add(str string) {
	s[str] = struct{}{}
}

func (s Set) Has(str string) bool {
	_, ok := s[str]
	return ok
}

func (s Set) Remove(str string) {
	delete(s, str)
}
