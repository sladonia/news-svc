package post

type storage struct {
}

func (s *storage) ByID(id string) (Post, error) {
	panic("implement me")
}

func (s *storage) ByFilter(filter Filter) ([]Post, error) {
	panic("implement me")
}

func (s *storage) Insert(post Post) error {
	panic("implement me")
}

func (s *storage) Replace(post Post) error {
	panic("implement me")
}

func (s *storage) Remove(post Post) error {
	panic("implement me")
}
