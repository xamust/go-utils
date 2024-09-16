package rest

type StatusCode interface {
	SetStatus(s int)
	Status() int
}

type StatusNoop struct {
	statusCode int
}

func (s *StatusNoop) SetStatus(n int) {
	s.statusCode = n
}

func (s *StatusNoop) Status() int {
	return s.statusCode
}
