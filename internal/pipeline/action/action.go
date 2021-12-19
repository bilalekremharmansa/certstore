package action

type Action interface {
	Run(map[string]string) error
}
