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
