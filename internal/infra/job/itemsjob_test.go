package job_test

import (
	"github.com/cubny/lite-reader/internal/app/feed"
	"github.com/cubny/lite-reader/internal/app/item"
	"github.com/cubny/lite-reader/internal/infra/job"
	mocks "github.com/cubny/lite-reader/internal/mocks/infra/job"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestItemsJob_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	feedService := mocks.NewFeedService(ctrl)
	itemService := mocks.NewItemService(ctrl)
	j := job.NewItemsJob(feedService, itemService)

	t.Run("Success", func(t *testing.T) {
		feedService.EXPECT().ListFeeds().Return([]*feed.Feed{
			{Id: 1},
			{Id: 2},
		}, nil)
		feedService.EXPECT().FetchItems(1).Return([]*item.Item{
			{Id: 1},
			{Id: 2},
		}, nil)
		feedService.EXPECT().FetchItems(2).Return([]*item.Item{
			{Id: 3},
			{Id: 4},
		}, nil)
		itemService.EXPECT().UpsertItems(&item.UpsertItemsCommand{
			FeedId: 1,
			Items: []*item.Item{
				{Id: 1},
				{Id: 2},
			},
		}).Return(nil)
		itemService.EXPECT().UpsertItems(&item.UpsertItemsCommand{
			FeedId: 2,
			Items: []*item.Item{
				{Id: 3},
				{Id: 4},
			},
		}).Return(nil)
		j.Execute()
	})

	t.Run("FailListFeeds", func(t *testing.T) {
		feedService.EXPECT().ListFeeds().Return(nil, assert.AnError)
		itemService.EXPECT().UpsertItems(gomock.Any()).Times(0)
		j.Execute()
	})

	t.Run("FailFetchItems", func(t *testing.T) {
		feedService.EXPECT().ListFeeds().Return([]*feed.Feed{
			{Id: 1},
		}, nil)
		feedService.EXPECT().FetchItems(1).Return(nil, assert.AnError)
		itemService.EXPECT().UpsertItems(gomock.Any()).Times(0)
		j.Execute()
	})

	t.Run("FailUpsertItems", func(t *testing.T) {
		feedService.EXPECT().ListFeeds().Return([]*feed.Feed{
			{Id: 1},
			{Id: 2},
		}, nil)
		feedService.EXPECT().FetchItems(1).Return([]*item.Item{
			{Id: 1},
		}, nil)
		feedService.EXPECT().FetchItems(2).Return([]*item.Item{
			{Id: 2},
		}, nil)
		itemService.EXPECT().UpsertItems(&item.UpsertItemsCommand{
			FeedId: 1,
			Items: []*item.Item{
				{Id: 1},
			},
		}).Return(assert.AnError)
		itemService.EXPECT().UpsertItems(&item.UpsertItemsCommand{
			FeedId: 2,
			Items: []*item.Item{
				{Id: 2},
			},
		}).Return(nil)
		j.Execute()
	})
}
