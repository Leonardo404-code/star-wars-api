package services

import (
	"github.com/peterhellberg/swapi"
)

func SetupService() *swapi.Client {
	sw := swapi.DefaultClient

	return sw
}
