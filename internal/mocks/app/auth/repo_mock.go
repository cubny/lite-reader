// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/cubny/lite-reader/internal/app/auth (interfaces: Repository)
//
// Generated by this command:
//
//	mockgen -destination=./app/auth/repo_mock.go -package=mocks -mock_names=Repository=Repository github.com/cubny/lite-reader/internal/app/auth Repository
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"

	auth "github.com/cubny/lite-reader/internal/app/auth"
)

// Repository is a mock of Repository interface.
type Repository struct {
	ctrl     *gomock.Controller
	recorder *RepositoryMockRecorder
	isgomock struct{}
}

// RepositoryMockRecorder is the mock recorder for Repository.
type RepositoryMockRecorder struct {
	mock *Repository
}

// NewRepository creates a new mock instance.
func NewRepository(ctrl *gomock.Controller) *Repository {
	mock := &Repository{ctrl: ctrl}
	mock.recorder = &RepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Repository) EXPECT() *RepositoryMockRecorder {
	return m.recorder
}

// CreateSession mocks base method.
func (m *Repository) CreateSession(userID int) (*auth.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", userID)
	ret0, _ := ret[0].(*auth.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSession indicates an expected call of CreateSession.
func (mr *RepositoryMockRecorder) CreateSession(userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*Repository)(nil).CreateSession), userID)
}

// CreateUser mocks base method.
func (m *Repository) CreateUser(email, password string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", email, password)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *RepositoryMockRecorder) CreateUser(email, password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*Repository)(nil).CreateUser), email, password)
}

// GetAllUsers mocks base method.
func (m *Repository) GetAllUsers() ([]*auth.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllUsers")
	ret0, _ := ret[0].([]*auth.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllUsers indicates an expected call of GetAllUsers.
func (mr *RepositoryMockRecorder) GetAllUsers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllUsers", reflect.TypeOf((*Repository)(nil).GetAllUsers))
}

// GetSessionByToken mocks base method.
func (m *Repository) GetSessionByToken(token string) (*auth.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSessionByToken", token)
	ret0, _ := ret[0].(*auth.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSessionByToken indicates an expected call of GetSessionByToken.
func (mr *RepositoryMockRecorder) GetSessionByToken(token any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSessionByToken", reflect.TypeOf((*Repository)(nil).GetSessionByToken), token)
}

// GetUserByEmail mocks base method.
func (m *Repository) GetUserByEmail(email string) (*auth.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", email)
	ret0, _ := ret[0].(*auth.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *RepositoryMockRecorder) GetUserByEmail(email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*Repository)(nil).GetUserByEmail), email)
}
