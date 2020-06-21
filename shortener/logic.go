package shortener

import (
	"errors"
	"time"

	errs "github.com/pkg/errors"
	"github.com/teris-io/shortid"
	"gopkg.in/dealancer/validate.v2"
)

var (
	ErrRedirectNotFound = errors.New("Redirect Not Found!")
	ErrRedirectInvalid  = errors.New("Redirect Invalid!")
)

type redirectService struct {
	redirectRepo RedirectRepository
}

func NewRedirectService(redirectRepo RedirectRepository) RedirectService {
	return &redirectService{redirectRepo}
}

func (s *redirectService) Find(hash string) (*Redirect, error) {
	return s.redirectRepo.Find(hash)
}

func (s *redirectService) Store(redirect *Redirect) error {
	if err := validate.Validate(redirect); err != nil {
		return errs.Wrap(ErrRedirectInvalid, "service.Redirect")
	}
	redirect.Hash = shortid.MustGenerate()
	redirect.CreatedAt = time.Now().UTC().Unix()

	return s.redirectRepo.Store(redirect)
}
