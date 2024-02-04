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
