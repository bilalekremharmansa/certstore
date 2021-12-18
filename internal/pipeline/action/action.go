package action

type action interface {
	run(map[string]string) error
}
