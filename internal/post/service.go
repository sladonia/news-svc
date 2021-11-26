package post

type Service interface {
	GetPost(id string) (Post, error)
	CreatePost(title, content string) (Post, error)
	UpsertPost(id, title, content string) error
	DeletePost(id string) error
	FindPosts(f Filter) ([]Post, error)
}

func NewService(storage Storage) Service {
	return &service{storage: storage}
}

type service struct {
	storage Storage
}

func (s *service) GetPost(id string) (Post, error) {
	return s.storage.ByID(id)
}

func (s *service) CreatePost(title, content string) (Post, error) {
	p := NewPost(title, content)

	err := s.storage.Insert(p)

	return p, err
}

func (s *service) UpsertPost(id, title, content string) error {
	err := s.storage.Update(id, title, content)
	if err == nil {
		return nil
	}

	newPost := NewPost(title, content)
	newPost.ID = id

	return s.storage.Insert(newPost)
}

func (s *service) DeletePost(id string) error {
	return s.storage.Remove(id)
}

func (s *service) FindPosts(f Filter) ([]Post, error) {
	return s.storage.ByFilter(f)
}
