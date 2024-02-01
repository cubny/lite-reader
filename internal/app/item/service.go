package item

type ServiceImpl struct {
	repository Repository
}

func NewService(repository Repository) (*ServiceImpl, error) {
	return &ServiceImpl{repository: repository}, nil
}

func (s *ServiceImpl) GetUnreadItems(command *GetUnreadItemsCommand) ([]*Item, error) {
	return s.repository.GetUnreadItems()
}

func (s *ServiceImpl) GetStarredItems(command *GetStarredItemsCommand) ([]*Item, error) {
	return s.repository.GetStarredItems()
}

func (s *ServiceImpl) GetFeedItems(command *GetFeedItemsCommand) ([]*Item, error) {
	return s.repository.GetFeedItems(command.FeedId)
}

func (s *ServiceImpl) UpsertItems(command *UpsertItemsCommand) error {
	return s.repository.UpsertItems(command.FeedId, command.Items)
}

func (s *ServiceImpl) UpdateItem(command *UpdateItemCommand) error {
	return s.repository.UpdateItem(command.Id, command.Starred, command.IsNew)
}

func (s *ServiceImpl) ReadFeedItems(command *ReadFeedItemsCommand) error {
	return s.repository.ReadFeedItems(command.FeedId)
}

func (s *ServiceImpl) UnreadFeedItems(command *UnreadFeedItemsCommand) error {
	return s.repository.UnreadFeedItems(command.FeedId)
}

func (s *ServiceImpl) GetStarredItemsCount() (int, error) {
	return s.repository.GetStarredItemsCount()
}

func (s *ServiceImpl) GetUnreadItemsCount() (int, error) {
	return s.repository.GetUnreadItemsCount()
}

func (s *ServiceImpl) DeleteFeedItems(command *DeleteFeedItemsCommand) error {
	return s.repository.DeleteFeedItems(command.FeedId)
}
