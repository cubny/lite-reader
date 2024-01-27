package item

import "time"

type ServiceImpl struct {
}

func NewService() (*ServiceImpl, error) {
	return &ServiceImpl{}, nil
}

func (s *ServiceImpl) GetUnreadItems(command *GetUnreadItemsCommand) ([]*Item, error) {
	return []*Item{
		{
			Id:        "1",
			Title:     "title 1",
			Desc:      "desc 1",
			Link:      "www.google.com",
			IsNew:     true,
			Starred:   true,
			Timestamp: time.Now(),
		},
		{
			Id:        "2",
			Title:     "title 2",
			Desc:      "desc 2",
			Link:      "www.google.com",
			IsNew:     true,
			Starred:   false,
			Timestamp: time.Now(),
		},
	}, nil
}

func (s *ServiceImpl) GetStarredItems(command *GetStarredItemsCommand) ([]*Item, error) {
	return []*Item{
		{
			Id:        "1",
			Title:     "title 1",
			Desc:      "desc 1",
			Link:      "www.google.com",
			IsNew:     true,
			Starred:   true,
			Timestamp: time.Now(),
		},
	}, nil
}

func (s *ServiceImpl) GetFeedItems(command *GetFeedItemsCommand) ([]*Item, error) {
	return []*Item{
		{
			Id:        "1",
			Title:     "title 1",
			Desc:      "desc 1",
			Link:      "www.google.com",
			IsNew:     true,
			Starred:   true,
			Timestamp: time.Now(),
		},
	}, nil
}
