//go:build mocks

package mocks

import (
	// ensures this package is in vendors folder and fixes a bug in go:generate that appears because of use of reflection in mocks generation
	_ "go.uber.org/mock/mockgen"
)

//go:generate mockgen -destination=./app/feed/repo_mock.go -package=mocks -mock_names=Repository=Repository github.com/cubny/lite-reader/internal/app/feed Repository
//go:generate mockgen -destination=./app/feed/parser_mock.go -package=mocks -mock_names=Parser=Parser       github.com/cubny/lite-reader/internal/app/feed Parser
//go:generate mockgen -destination=./app/feed/finder_mock.go -package=mocks -mock_names=Finder=Finder       github.com/cubny/lite-reader/internal/app/feed Finder
//go:generate mockgen -destination=./app/item/repo_mock.go -package=mocks -mock_names=Repository=Repository github.com/cubny/lite-reader/internal/app/item Repository
//go:generate mockgen -destination=./infra/http/api/feed_mock.go -package=mocks -mock_names=FeedService=FeedService   github.com/cubny/lite-reader/internal/infra/http/api FeedService
//go:generate mockgen -destination=./infra/http/api/item_mock.go -package=mocks -mock_names=ItemService=ItemService   github.com/cubny/lite-reader/internal/infra/http/api ItemService

//go:generate mockgen -destination=./infra/job/dependencies_mock.go -package=mocks -mock_names=FeedService=FeedService,ItemService=ItemService github.com/cubny/lite-reader/internal/infra/job FeedService,ItemService
