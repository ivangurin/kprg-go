package grapher

import (
	"kprg/internal/enricher"
	"kprg/internal/repository"
)

type Grapher struct {
	repository repository.Repository
	enricher   *enricher.Enricher
}

func NewGrapher(repository repository.Repository, enricher *enricher.Enricher) *Grapher {

	grapher := &Grapher{
		repository: repository,
		enricher:   enricher,
	}

	return grapher

}

func (g *Grapher) Listen(port string) error {

	return nil

}

func (g *Grapher) Close() error {

	return nil

}
