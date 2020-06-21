package shortener

type RedirectRepository interface {
	Find(hash string) (*Redirect, error)
	Store(redirect *Redirect) error
}
