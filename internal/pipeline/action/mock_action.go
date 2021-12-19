package action

type MockAction struct {
	Executed bool
}

func (m *MockAction) Run(args map[string]string) error {
	m.Executed = true
	return nil
}
