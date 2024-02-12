package feed_test

import (
	"testing"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/cubny/lite-reader/internal/app/feed"
	"github.com/cubny/lite-reader/internal/app/item"
	mocks "github.com/cubny/lite-reader/internal/mocks/app/feed"
)

func TestServiceImpl_AddFeed(t *testing.T) {
	ctrl := gomock.NewController(t)
	tests := []struct {
		name           string
		parseURLError  error
		parseURLResult *gofeed.Feed
		finderCalled   bool
		finderLinks    []string
		finderError    error
		repoCalled     bool
		repoResult     int
		repoError      error
		want           *feed.Feed
		wantErr        bool
	}{
		{
			name:          "success",
			parseURLError: nil,
			parseURLResult: &gofeed.Feed{
				Title:       "Example Feed",
				Description: "Example Description",
				Link:        "https://example.com",
				FeedLink:    "https://example.com/feed",
				Language:    "en",
				Items:       []*gofeed.Item{},
			},
			finderCalled: false,
			finderLinks:  []string{},
			finderError:  nil,
			repoCalled:   true,
			repoResult:   1,
			repoError:    nil,
			want: &feed.Feed{
				Id:          1,
				Title:       "Example Feed",
				Link:        "https//example.com",
				URL:         "https://example.com/feed",
				UnreadCount: 0,
			},
			wantErr: false,
		},
		{
			name:           "error - feed not detected",
			finderCalled:   true,
			finderLinks:    []string{},
			finderError:    nil,
			parseURLError:  gofeed.ErrFeedTypeNotDetected,
			parseURLResult: nil,
			repoCalled:     false,
			repoResult:     0,
			repoError:      nil,
			want:           nil,
			wantErr:        true,
		},
		{
			name:           "error - feed not found",
			parseURLError:  gofeed.ErrFeedTypeNotDetected,
			parseURLResult: nil,
			finderCalled:   true,
			finderLinks:    []string{},
			finderError:    assert.AnError,
			repoCalled:     false,
			repoResult:     0,
			repoError:      nil,
			want:           nil,
			wantErr:        true,
		},
		{
			name:           "error - parser unknown error",
			parseURLError:  assert.AnError,
			parseURLResult: nil,
			finderCalled:   false,
			finderLinks:    []string{},
			finderError:    nil,
			repoCalled:     false,
			repoResult:     0,
			repoError:      nil,
			want:           nil,
			wantErr:        true,
		},
		{
			name:           "error - repo error",
			parseURLError:  nil,
			parseURLResult: &gofeed.Feed{},
			finderCalled:   false,
			finderLinks:    []string{},
			finderError:    nil,
			repoCalled:     true,
			repoResult:     0,
			repoError:      assert.AnError,
			want:           nil,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewRepository(ctrl)
			if tt.repoCalled {
				repo.EXPECT().AddFeed(gomock.Any()).Return(tt.repoResult, tt.repoError)
			}

			parser := mocks.NewParser(ctrl)
			parser.EXPECT().ParseURL(gomock.Any()).Return(tt.parseURLResult, tt.parseURLError)

			finder := mocks.NewFinder(ctrl)
			if tt.finderCalled {
				finder.EXPECT().FindFeeds(gomock.Any()).Return(tt.finderLinks, tt.finderError)
			}

			s := feed.NewService(repo, parser, finder)
			cmd := &feed.AddFeedCommand{URL: "https://example.com/feed"}
			got, err := s.AddFeed(cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceImpl.AddFeed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.ObjectsAreEqualValues(got, tt.want)
		})
	}

	// finder.EXPECT().FindFeeds(gomock.Any()).Return([]string{"https://example.com/feed", "https://example.com/feed2", "https://example.com/feed3"}, nil)
	t.Run("parser is called as many times as there are links", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repoMock := mocks.NewRepository(ctrl)
		repoMock.EXPECT().AddFeed(gomock.Any()).Return(1, nil)

		parserMock := mocks.NewParser(ctrl)
		parserMock.EXPECT().ParseURL("https://example.com").Return(nil, gofeed.ErrFeedTypeNotDetected)
		parserMock.EXPECT().ParseURL("https://example.com/feed").Return(nil, assert.AnError)
		parserMock.EXPECT().ParseURL("https://example.com/feed2").Return(&gofeed.Feed{
			Title:       "Example Feed",
			Description: "Example Description",
			Link:        "https://example.com",
			FeedLink:    "https://example.com/feed",
			Language:    "en",
			Items:       []*gofeed.Item{},
		}, nil)

		finderMock := mocks.NewFinder(ctrl)
		finderMock.EXPECT().FindFeeds(gomock.Any()).Return([]string{"https://example.com/feed", "https://example.com/feed2"}, nil)

		s := feed.NewService(repoMock, parserMock, finderMock)
		cmd := &feed.AddFeedCommand{URL: "https://example.com"}
		_, err := s.AddFeed(cmd)
		assert.NoError(t, err)
	})
}

func TestServiceImpl_ListFeeds(t *testing.T) {
	ctrl := gomock.NewController(t)
	repoMock := mocks.NewRepository(ctrl)
	repoMock.EXPECT().ListFeeds().Return([]*feed.Feed{}, nil)

	parserMock := mocks.NewParser(ctrl)
	finderMock := mocks.NewFinder(ctrl)

	s := feed.NewService(repoMock, parserMock, finderMock)
	_, err := s.ListFeeds()
	assert.NoError(t, err)
}

func TestServiceImpl_FetchItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	samplePublishedParsed := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

	type Parser struct {
		result *gofeed.Feed
		error  error
		called bool
	}

	type Repo struct {
		result *feed.Feed
		error  error
	}
	tests := []struct {
		name    string
		parser  Parser
		repo    Repo
		want    []*item.Item
		wantErr bool
	}{
		{
			name: "success",
			parser: Parser{
				result: &gofeed.Feed{
					Title:       "Example Feed",
					Description: "Example Description",
					Link:        "https://example.com",
					FeedLink:    "https://example.com/feed",
					Language:    "en",
					Items: []*gofeed.Item{
						{Title: "Example Item", Description: "Example Description", Link: "https://example.com/item", Published: "2021-01-01T00:00:00Z", PublishedParsed: &samplePublishedParsed},
						{Title: "Example Item 2", Description: "Example Description 2", Link: "https://example.com/item2", Published: "2021-01-01T00:00:00Z", PublishedParsed: &samplePublishedParsed},
					},
				},
				error:  nil,
				called: true,
			},
			repo: Repo{
				result: &feed.Feed{
					Id:          1,
					Title:       "Example Feed",
					Link:        "https//example.com",
					URL:         "https://example.com/feed",
					UnreadCount: 0,
				},
				error: nil,
			},
			want: []*item.Item{
				{Title: "Example Item", Desc: "Example Description", Link: "https://example.com/item", Timestamp: samplePublishedParsed, Dir: "Example Description", IsNew: true, Starred: false},
				{Title: "Example Item 2", Desc: "Example Description 2", Link: "https://example.com/item2", Timestamp: samplePublishedParsed, Dir: "Example Description 2", IsNew: true, Starred: false},
			},
			wantErr: false,
		},
		{
			name: "error - repo error",
			parser: Parser{
				result: nil,
				error:  nil,
				called: false,
			},
			repo: Repo{
				result: nil,
				error:  assert.AnError,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error - parser error",
			parser: Parser{
				result: nil,
				error:  assert.AnError,
				called: true,
			},
			repo: Repo{
				result: &feed.Feed{
					Id:          1,
					Title:       "Example Feed",
					Link:        "https//example.com",
					URL:         "https://example.com/feed",
					UnreadCount: 0,
				},
				error: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := mocks.NewRepository(ctrl)
			repoMock.EXPECT().GetFeed(1).Return(tt.repo.result, tt.repo.error)

			parserMock := mocks.NewParser(ctrl)
			if tt.parser.called {
				parserMock.EXPECT().ParseURL(gomock.Any()).Return(tt.parser.result, tt.parser.error)
			}

			finderMock := mocks.NewFinder(ctrl)

			s := feed.NewService(repoMock, parserMock, finderMock)
			got, err := s.FetchItems(1)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceImpl.FetchItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.ObjectsAreEqualValues(got, tt.want)
		})
	}
}

func TestServiceImpl_DeleteFeed(t *testing.T) {
	ctrl := gomock.NewController(t)
	repoMock := mocks.NewRepository(ctrl)
	repoMock.EXPECT().DeleteFeed(1).Return(nil)

	parserMock := mocks.NewParser(ctrl)
	finderMock := mocks.NewFinder(ctrl)

	s := feed.NewService(repoMock, parserMock, finderMock)
	cmd := &feed.DeleteFeedCommand{FeedId: 1}
	err := s.DeleteFeed(cmd)
	assert.NoError(t, err)
}
