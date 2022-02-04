package lego

import (
	"github.com/go-acme/lego/certificate"
)

type LegoAdapter interface {
	Obtain(req certificate.ObtainRequest) (*certificate.Resource, error)
}
