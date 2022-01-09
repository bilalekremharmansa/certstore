package action

import "bilalekrem.com/certstore/internal/pipeline/context"

type Action interface {
	Run(*context.Context, map[string]string) error
}
