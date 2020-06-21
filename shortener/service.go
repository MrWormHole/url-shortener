package shortener

type RedirectService interface {
	Find(hash string) (*Redirect, error)
	Store(redirect *Redirect) error
}
