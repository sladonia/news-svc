package test

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/sladonia/news-svc/internal/handler"
	"github.com/sladonia/news-svc/internal/post"
)

func (s *Suite) TestGetPostByID() {
	s.Run("success", func() {
		res, err := http.Get(fmt.Sprintf("%s/posts/1", s.srv.URL))

		s.NoError(err)
		s.Equal(200, res.StatusCode)

		var p post.Post

		err = jsoniter.ConfigFastest.NewDecoder(res.Body).Decode(&p)
		s.NoError(err)
		s.Equal(post1, p)
	})

	s.Run("not_found", func() {
		res, err := http.Get(fmt.Sprintf("%s/posts/42", s.srv.URL))

		s.NoError(err)
		s.Equal(404, res.StatusCode)

		var apiError handler.ApiError

		err = jsoniter.ConfigFastest.NewDecoder(res.Body).Decode(&apiError)
		s.NoError(err)
		s.Equal("record not found", apiError.Error.Message)
	})
}

func (s *Suite) TestCreatePost() {
	requestBody := `{
	"title": "title1",
	"content": "content1"
}`

	r := strings.NewReader(requestBody)

	res, err := http.Post(fmt.Sprintf("%s/posts", s.srv.URL), "application/json", r)
	s.NoError(err)
	s.Equal(201, res.StatusCode)

	var p post.Post

	err = jsoniter.ConfigFastest.NewDecoder(res.Body).Decode(&p)

	s.NoError(err)
	s.Equal("title1", p.Title)
	s.Equal("content1", p.Content)

	createdID := p.ID

	fromStorage, err := s.storage.ByID(createdID)

	s.NoError(err)
	s.Equal(p, fromStorage)

	// create the same second time

	r = strings.NewReader(requestBody)

	res, err = http.Post(fmt.Sprintf("%s/posts", s.srv.URL), "application/json", r)
	s.NoError(err)
	s.Equal(201, res.StatusCode)

	err = jsoniter.ConfigFastest.NewDecoder(res.Body).Decode(&p)

	s.NoError(err)
	s.NotEqual(createdID, p.ID)
}

func (s *Suite) TestReplacePost() {
	requestBody := `{
	"title": "title1",
	"content": "content1"
}`

	r := strings.NewReader(requestBody)

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/posts/1", s.srv.URL), r)
	s.NoError(err)

	res, err := http.DefaultClient.Do(req)
	s.NoError(err)
	s.Equal(200, res.StatusCode)

	fromStorage, err := s.storage.ByID("1")

	s.NoError(err)
	s.Equal("title1", fromStorage.Title)
	s.Equal("content1", fromStorage.Content)

	// upsert
	r = strings.NewReader(requestBody)

	req, err = http.NewRequest("PUT", fmt.Sprintf("%s/posts/unexisting_id", s.srv.URL), r)
	s.NoError(err)

	res, err = http.DefaultClient.Do(req)
	s.NoError(err)
	s.Equal(200, res.StatusCode)

	fromStorage, err = s.storage.ByID("unexisting_id")

	s.NoError(err)
	s.Equal("title1", fromStorage.Title)
	s.Equal("content1", fromStorage.Content)
}

func (s *Suite) TestDeletePost() {
	s.Run("success", func() {
		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/posts/1", s.srv.URL), nil)
		s.NoError(err)

		res, err := http.DefaultClient.Do(req)
		s.NoError(err)
		s.Equal(204, res.StatusCode)

		_, err = s.storage.ByID("1")
		s.ErrorIs(err, post.ErrNotFound)
	})

	s.Run("no_data", func() {
		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/posts/unexisting_id", s.srv.URL), nil)
		s.NoError(err)

		res, err := http.DefaultClient.Do(req)
		s.NoError(err)
		s.Equal(204, res.StatusCode)
	})
}

func (s *Suite) TestGetPosts() {
	s.Run("success", func() {
		res, err := http.Get(fmt.Sprintf("%s/posts", s.srv.URL))
		s.NoError(err)
		s.Equal(200, res.StatusCode)

		var posts []post.Post

		err = jsoniter.ConfigFastest.NewDecoder(res.Body).Decode(&posts)
		s.NoError(err)
		s.Len(posts, 1)
	})

	s.Run("no_data", func() {
		res, err := http.Get(fmt.Sprintf("%s/posts?from=%s", s.srv.URL, time.Now().
			UTC().
			Add(time.Hour).
			Round(time.Millisecond).
			Format(time.RFC3339)))
		s.NoError(err)
		s.Equal(200, res.StatusCode)

		var posts []post.Post

		err = jsoniter.ConfigFastest.NewDecoder(res.Body).Decode(&posts)
		s.NoError(err)
		s.Len(posts, 0)
	})
}
