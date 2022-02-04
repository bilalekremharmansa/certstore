package lego

import (
	"github.com/go-acme/lego/v4/certificate"
)

type LegoAdapter interface {
	Obtain(req certificate.ObtainRequest) (*certificate.Resource, error)
}
