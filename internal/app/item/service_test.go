package item_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/cubny/lite-reader/internal/app/item"
	mocks "github.com/cubny/lite-reader/internal/mocks/app/item"
)

const (
	tcHappy      = "happy"
	tcRepoErrors = "repo errors"
	itemTitle    = "title"
	itemDesc     = "desc"
	itemDir      = "dir"
	itemLink     = "link"
)

func TestServiceImpl_GetUnreadItems(t *testing.T) {
	type Repo struct {
		result []*item.Item
		err    error
	}
	tests := []struct {
		name    string
		repo    Repo
		want    []*item.Item
		wantErr bool
	}{
		{
			name: tcHappy,
			repo: Repo{
				result: []*item.Item{
					{
						ID:        1,
						Title:     itemTitle,
						Desc:      itemDesc,
						Dir:       itemDir,
						Link:      itemLink,
						IsNew:     true,
						Starred:   false,
						Timestamp: time.Now(),
					},
				},
				err: nil,
			},
			want: []*item.Item{
				{
					ID:        1,
					Title:     itemTitle,
					Desc:      itemDesc,
					Dir:       itemDir,
					Link:      itemLink,
					IsNew:     true,
					Starred:   false,
					Timestamp: time.Now(),
				},
			}, wantErr: false,
		},
		{
			name: tcRepoErrors,
			repo: Repo{
				result: nil,
				err:    assert.AnError,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewRepository(ctrl)
			repo.EXPECT().GetUnreadItems().Return(tt.repo.result, tt.repo.err)

			s := item.NewService(repo, nil)
			got, err := s.GetUnreadItems()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUnreadItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("GetUnreadItems() got = %v, want %v", got, tt.want)
				return
			}
			for i := range got {
				assert.ObjectsAreEqual(got[i], tt.want[i])
			}
		})
	}
}

func TestServiceImpl_GetStarredItems(t *testing.T) {
	type Repo struct {
		result []*item.Item
		err    error
	}
	tests := []struct {
		name    string
		repo    Repo
		want    []*item.Item
		wantErr bool
	}{
		{
			name: tcHappy,
			repo: Repo{
				result: []*item.Item{
					{
						ID:        1,
						Title:     itemTitle,
						Desc:      itemDesc,
						Dir:       itemDir,
						Link:      itemLink,
						IsNew:     true,
						Starred:   false,
						Timestamp: time.Now(),
					},
				},
				err: nil,
			},
			want: []*item.Item{
				{
					ID:        1,
					Title:     itemTitle,
					Desc:      itemDesc,
					Dir:       itemDir,
					Link:      itemLink,
					IsNew:     true,
					Starred:   false,
					Timestamp: time.Now(),
				},
			}, wantErr: false,
		},
		{
			name: tcRepoErrors,
			repo: Repo{
				result: nil,
				err:    assert.AnError,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewRepository(ctrl)
			repo.EXPECT().GetStarredItems().Return(tt.repo.result, tt.repo.err)

			s := item.NewService(repo, nil)
			got, err := s.GetStarredItems()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStarredItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("GetStarredItems() got = %v, want %v", got, tt.want)
				return
			}
			for i := range got {
				assert.ObjectsAreEqual(got[i], tt.want[i])
			}
		})
	}
}

func TestServiceImpl_GetFeedItems(t *testing.T) {
	type Repo struct {
		result []*item.Item
		err    error
	}
	type Command struct {
		FeedID int
	}
	tests := []struct {
		name    string
		repo    Repo
		command Command
		want    []*item.Item
		wantErr bool
	}{
		{
			name: tcHappy,
			repo: Repo{
				result: []*item.Item{
					{
						ID:        1,
						Title:     itemTitle,
						Desc:      itemDesc,
						Dir:       itemDir,
						Link:      itemLink,
						IsNew:     true,
						Starred:   false,
						Timestamp: time.Now(),
					},
				},
				err: nil,
			},
			command: Command{
				FeedID: 1,
			},
			want: []*item.Item{
				{
					ID:        1,
					Title:     itemTitle,
					Desc:      itemDesc,
					Dir:       itemDir,
					Link:      itemLink,
					IsNew:     true,
					Starred:   false,
					Timestamp: time.Now(),
				},
			}, wantErr: false,
		},
		{
			name: tcRepoErrors,
			repo: Repo{
				result: nil,
				err:    assert.AnError,
			},
			command: Command{
				FeedID: 1,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewRepository(ctrl)
			repo.EXPECT().GetFeedItems(tt.command.FeedID).Return(tt.repo.result, tt.repo.err)

			s := item.NewService(repo, nil)
			got, err := s.GetFeedItems(&item.GetFeedItemsCommand{FeedID: tt.command.FeedID})
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFeedItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("GetFeedItems() got = %v, want %v", got, tt.want)
				return
			}
			for i := range got {
				assert.ObjectsAreEqual(got[i], tt.want[i])
			}
		})
	}
}

func TestServiceImpl_UpsertItems(t *testing.T) {
	type Repo struct {
		err error
	}
	type Command struct {
		FeedID int
		Items  []*item.Item
	}
	tests := []struct {
		name    string
		repo    Repo
		command Command
		wantErr bool
	}{
		{
			name: tcHappy,
			repo: Repo{
				err: nil,
			},
			command: Command{
				FeedID: 1,
				Items: []*item.Item{
					{
						ID:        1,
						Title:     itemTitle,
						Desc:      itemDesc,
						Dir:       itemDir,
						Link:      itemLink,
						IsNew:     true,
						Starred:   false,
						Timestamp: time.Now(),
					},
				},
			}, wantErr: false,
		},
		{
			name: tcRepoErrors,
			repo: Repo{
				err: assert.AnError,
			},
			command: Command{
				FeedID: 1,
				Items: []*item.Item{
					{
						ID:        1,
						Title:     itemTitle,
						Desc:      itemDesc,
						Dir:       itemDir,
						Link:      itemLink,
						IsNew:     true,
						Starred:   false,
						Timestamp: time.Now(),
					},
				},
			}, wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewRepository(ctrl)
			repo.EXPECT().UpsertItems(tt.command.FeedID, tt.command.Items).Return(tt.repo.err)

			s := item.NewService(repo, nil)
			err := s.UpsertItems(&item.UpsertItemsCommand{FeedID: tt.command.FeedID, Items: tt.command.Items})
			if (err != nil) != tt.wantErr {
				t.Errorf("UpsertItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestServiceImpl_UpdateItem(t *testing.T) {
	type Repo struct {
		err error
	}
	type Command struct {
		ID      int
		Starred bool
		IsNew   bool
	}
	tests := []struct {
		name    string
		repo    Repo
		command Command
		wantErr bool
	}{
		{
			name: tcHappy,
			repo: Repo{
				err: nil,
			},
			command: Command{
				ID:      1,
				Starred: true,
				IsNew:   false,
			}, wantErr: false,
		},
		{
			name: tcRepoErrors,
			repo: Repo{
				err: assert.AnError,
			},
			command: Command{
				ID:      1,
				Starred: true,
				IsNew:   false,
			}, wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewRepository(ctrl)
			repo.EXPECT().UpdateItem(tt.command.ID, tt.command.Starred, tt.command.IsNew).Return(tt.repo.err)

			s := item.NewService(repo, nil)
			err := s.UpdateItem(&item.UpdateItemCommand{ID: tt.command.ID, Starred: tt.command.Starred, IsNew: tt.command.IsNew})
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestServiceImpl_ReadFeedItems(t *testing.T) {
	type Repo struct {
		err error
	}
	type Command struct {
		FeedID int
	}
	tests := []struct {
		name    string
		repo    Repo
		command Command
		wantErr bool
	}{
		{
			name: tcHappy,
			repo: Repo{
				err: nil,
			},
			command: Command{
				FeedID: 1,
			}, wantErr: false,
		},
		{
			name: tcRepoErrors,
			repo: Repo{
				err: assert.AnError,
			},
			command: Command{
				FeedID: 1,
			}, wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewRepository(ctrl)
			repo.EXPECT().ReadFeedItems(tt.command.FeedID).Return(tt.repo.err)

			s := item.NewService(repo, nil)
			err := s.ReadFeedItems(&item.ReadFeedItemsCommand{FeedID: tt.command.FeedID})
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFeedItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestServiceImpl_UnreadFeedItems(t *testing.T) {
	type Repo struct {
		err error
	}
	type Command struct {
		FeedID int
	}
	tests := []struct {
		name    string
		repo    Repo
		command Command
		wantErr bool
	}{
		{
			name: tcHappy,
			repo: Repo{
				err: nil,
			},
			command: Command{
				FeedID: 1,
			}, wantErr: false,
		},
		{
			name: tcRepoErrors,
			repo: Repo{
				err: assert.AnError,
			},
			command: Command{
				FeedID: 1,
			}, wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewRepository(ctrl)
			repo.EXPECT().UnreadFeedItems(tt.command.FeedID).Return(tt.repo.err)

			s := item.NewService(repo, nil)
			err := s.UnreadFeedItems(&item.UnreadFeedItemsCommand{FeedID: tt.command.FeedID})
			if (err != nil) != tt.wantErr {
				t.Errorf("UnreadFeedItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestServiceImpl_GetStarredItemsCount(t *testing.T) {
	type Repo struct {
		result int
		err    error
	}
	tests := []struct {
		name    string
		repo    Repo
		want    int
		wantErr bool
	}{
		{
			name: tcHappy,
			repo: Repo{
				result: 1,
				err:    nil,
			},
			want:    1,
			wantErr: false,
		},
		{
			name: tcRepoErrors,
			repo: Repo{
				result: 0,
				err:    assert.AnError,
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewRepository(ctrl)
			repo.EXPECT().GetStarredItemsCount().Return(tt.repo.result, tt.repo.err)

			s := item.NewService(repo, nil)
			got, err := s.GetStarredItemsCount()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStarredItemsCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetStarredItemsCount() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceImpl_GetUnreadItemsCount(t *testing.T) {
	type Repo struct {
		result int
		err    error
	}
	tests := []struct {
		name    string
		repo    Repo
		want    int
		wantErr bool
	}{
		{
			name: tcHappy,
			repo: Repo{
				result: 1,
				err:    nil,
			},
			want:    1,
			wantErr: false,
		},
		{
			name: tcRepoErrors,
			repo: Repo{
				result: 0,
				err:    assert.AnError,
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewRepository(ctrl)
			repo.EXPECT().GetUnreadItemsCount().Return(tt.repo.result, tt.repo.err)

			s := item.NewService(repo, nil)
			got, err := s.GetUnreadItemsCount()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUnreadItemsCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUnreadItemsCount() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceImpl_DeleteFeedItems(t *testing.T) {
	type Repo struct {
		err error
	}
	type Command struct {
		FeedID int
	}
	tests := []struct {
		name    string
		repo    Repo
		command Command
		wantErr bool
	}{
		{
			name: tcHappy,
			repo: Repo{
				err: nil,
			},
			command: Command{
				FeedID: 1,
			}, wantErr: false,
		},
		{
			name: tcRepoErrors,
			repo: Repo{
				err: assert.AnError,
			},
			command: Command{
				FeedID: 1,
			}, wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewRepository(ctrl)
			repo.EXPECT().DeleteFeedItems(tt.command.FeedID).Return(tt.repo.err)

			s := item.NewService(repo, nil)
			err := s.DeleteFeedItems(&item.DeleteFeedItemsCommand{FeedID: tt.command.FeedID})
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteFeedItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestServiceImpl_ScrapeItem(t *testing.T) {
	const userID = 7
	const itemID = 42
	const link = "https://example.test/article"

	t.Run("ok path persists scraped content and returns updated item", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := mocks.NewRepository(ctrl)
		scraper := mocks.NewScraper(ctrl)

		repo.EXPECT().GetItemForUser(itemID, userID).Return(&item.Item{ID: itemID, Link: link}, nil)
		scraper.EXPECT().Scrape(gomock.Any(), link).Return("<p>full body</p>", nil)
		repo.EXPECT().UpdateItemContent(itemID, "<p>full body</p>", "ok", gomock.Any()).Return(nil)

		s := item.NewService(repo, scraper)
		got, err := s.ScrapeItem(t.Context(), &item.ScrapeItemCommand{ID: itemID, UserID: userID})
		assert.NoError(t, err)
		assert.Equal(t, "<p>full body</p>", got.FullContent)
		assert.Equal(t, "ok", got.ScrapeStatus)
		assert.NotNil(t, got.ScrapedAt)
	})

	t.Run("scraper error persists 'error' status and returns the scraper err", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := mocks.NewRepository(ctrl)
		scraper := mocks.NewScraper(ctrl)

		repo.EXPECT().GetItemForUser(itemID, userID).Return(&item.Item{ID: itemID, Link: link}, nil)
		scraper.EXPECT().Scrape(gomock.Any(), link).Return("", assert.AnError)
		repo.EXPECT().UpdateItemContent(itemID, "", "error", gomock.Any()).Return(nil)

		s := item.NewService(repo, scraper)
		_, err := s.ScrapeItem(t.Context(), &item.ScrapeItemCommand{ID: itemID, UserID: userID})
		assert.Error(t, err)
	})

	t.Run("not found returns ErrItemNotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := mocks.NewRepository(ctrl)
		scraper := mocks.NewScraper(ctrl)

		repo.EXPECT().GetItemForUser(itemID, userID).Return(nil, nil)

		s := item.NewService(repo, scraper)
		_, err := s.ScrapeItem(t.Context(), &item.ScrapeItemCommand{ID: itemID, UserID: userID})
		assert.ErrorIs(t, err, item.ErrItemNotFound)
	})
}
