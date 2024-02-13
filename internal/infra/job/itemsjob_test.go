package job_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/cubny/lite-reader/internal/app/feed"
	"github.com/cubny/lite-reader/internal/app/item"
	"github.com/cubny/lite-reader/internal/infra/job"
	mocks "github.com/cubny/lite-reader/internal/mocks/infra/job"
)

func TestItemsJob_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	feedService := mocks.NewFeedService(ctrl)
	itemService := mocks.NewItemService(ctrl)
	j := job.NewItemsJob(feedService, itemService)

	t.Run("Success", func(_ *testing.T) {
		feedService.EXPECT().ListFeeds().Return([]*feed.Feed{
			{ID: 1},
			{ID: 2},
		}, nil)
		feedService.EXPECT().FetchItems(1).Return([]*item.Item{
			{ID: 1},
			{ID: 2},
		}, nil)
		feedService.EXPECT().FetchItems(2).Return([]*item.Item{
			{ID: 3},
			{ID: 4},
		}, nil)
		itemService.EXPECT().UpsertItems(&item.UpsertItemsCommand{
			FeedID: 1,
			Items: []*item.Item{
				{ID: 1},
				{ID: 2},
			},
		}).Return(nil)
		itemService.EXPECT().UpsertItems(&item.UpsertItemsCommand{
			FeedID: 2,
			Items: []*item.Item{
				{ID: 3},
				{ID: 4},
			},
		}).Return(nil)
		j.Execute()
	})

	t.Run("FailListFeeds", func(_ *testing.T) {
		feedService.EXPECT().ListFeeds().Return(nil, assert.AnError)
		itemService.EXPECT().UpsertItems(gomock.Any()).Times(0)
		j.Execute()
	})

	t.Run("FailFetchItems", func(_ *testing.T) {
		feedService.EXPECT().ListFeeds().Return([]*feed.Feed{
			{ID: 1},
		}, nil)
		feedService.EXPECT().FetchItems(1).Return(nil, assert.AnError)
		itemService.EXPECT().UpsertItems(gomock.Any()).Times(0)
		j.Execute()
	})

	t.Run("FailUpsertItems", func(_ *testing.T) {
		feedService.EXPECT().ListFeeds().Return([]*feed.Feed{
			{ID: 1},
			{ID: 2},
		}, nil)
		feedService.EXPECT().FetchItems(1).Return([]*item.Item{
			{ID: 1},
		}, nil)
		feedService.EXPECT().FetchItems(2).Return([]*item.Item{
			{ID: 2},
		}, nil)
		itemService.EXPECT().UpsertItems(&item.UpsertItemsCommand{
			FeedID: 1,
			Items: []*item.Item{
				{ID: 1},
			},
		}).Return(assert.AnError)
		itemService.EXPECT().UpsertItems(&item.UpsertItemsCommand{
			FeedID: 2,
			Items: []*item.Item{
				{ID: 2},
			},
		}).Return(nil)
		j.Execute()
	})
}
