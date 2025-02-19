package auth

type Service struct {
	Repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{Repo: r}
}

func (s *Service) Create(u *User) error {
	return s.Repo.Save(u)
}

func (s *Service) UserExists(u string, e string) (bool, error) {
	return s.Repo.UserAlreadyExists(u, e)
}

func (s *Service) GetByUsername(u string) (*User, error) {
	return s.Repo.GetByUsername(u)
}

func (s *Service) SetSessionToken(t string, u string) error {
	return s.Repo.SetSessionToken(t, u)
}

func (s *Service) SetCSRFToken(t string, u string) error {
	return s.Repo.SetCSRFToken(t, u)
}

func (s *Service) ResetTokens(username string) error {
	return s.Repo.ResetTokensForUser(username)
}
