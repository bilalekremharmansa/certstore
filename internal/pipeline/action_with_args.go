package pipeline

import "bilalekrem.com/certstore/internal/pipeline/action"

type ActionWithArgs struct {
	action action.Action
	args   map[string]string
}
