package auth_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/cubny/lite-reader/internal/app/auth"
	mocks "github.com/cubny/lite-reader/internal/mocks/app/auth"
)

func TestService_Signup(t *testing.T) {
	ctrl := gomock.NewController(t)
	tests := []struct {
		name      string
		command   *auth.SignupCommand
		mockSetup func(*mocks.Repository)
		wantErr   bool
		errMsg    string
	}{
		{
			name: "success",
			command: &auth.SignupCommand{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(r *mocks.Repository) {
				r.EXPECT().GetUserByEmail("test@example.com").Return(nil, errors.New("not found"))
				r.EXPECT().CreateUser("test@example.com", gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error - email already exists",
			command: &auth.SignupCommand{
				Email:    "existing@example.com",
				Password: "password123",
			},
			mockSetup: func(r *mocks.Repository) {
				r.EXPECT().GetUserByEmail("existing@example.com").Return(&auth.User{}, nil)
			},
			wantErr: true,
			errMsg:  "email already registered",
		},
		{
			name: "error - create user fails",
			command: &auth.SignupCommand{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(r *mocks.Repository) {
				r.EXPECT().GetUserByEmail("test@example.com").Return(nil, errors.New("not found"))
				r.EXPECT().CreateUser("test@example.com", gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewRepository(ctrl)
			tt.mockSetup(repo)

			s := auth.NewService(repo)
			err := s.Signup(tt.command)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
			ctrl.Finish()
		})
	}
}

func TestService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	tests := []struct {
		name      string
		command   *auth.LoginCommand
		mockSetup func(*mocks.Repository)
		want      *auth.LoginResponse
		wantErr   bool
		errMsg    string
	}{
		{
			name: "success",
			command: &auth.LoginCommand{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(r *mocks.Repository) {
				hashedPassword := "$2a$10$abcdefghijklmnopqrstuvwxyz" // pre-hashed "password123"
				user := &auth.User{
					ID:       1,
					Email:    "test@example.com",
					Password: hashedPassword,
				}
				session := &auth.Session{
					AccessToken:  "access-token",
					RefreshToken: "refresh-token",
					ExpiresAt:    time.Now().Add(time.Hour),
				}
				r.EXPECT().GetUserByEmail("test@example.com").Return(user, nil)
				r.EXPECT().CreateSession(user.ID).Return(session, nil)
			},
			wantErr: false,
		},
		{
			name: "error - user not found",
			command: &auth.LoginCommand{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			mockSetup: func(r *mocks.Repository) {
				r.EXPECT().GetUserByEmail("nonexistent@example.com").Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "invalid email or password",
		},
		{
			name: "error - invalid password",
			command: &auth.LoginCommand{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockSetup: func(r *mocks.Repository) {
				hashedPassword := "$2a$10$abcdefghijklmnopqrstuvwxyz" // pre-hashed "password123"
				user := &auth.User{
					ID:       1,
					Email:    "test@example.com",
					Password: hashedPassword,
				}
				r.EXPECT().GetUserByEmail("test@example.com").Return(user, nil)
			},
			wantErr: true,
			errMsg:  "invalid email or password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewRepository(ctrl)
			tt.mockSetup(repo)

			s := auth.NewService(repo)
			got, err := s.Login(tt.command)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.NotEmpty(t, got.AccessToken)
				assert.NotEmpty(t, got.RefreshToken)
				assert.Greater(t, got.ExpiresIn, float64(0))
			}
			ctrl.Finish()
		})
	}
}

func TestService_GetSession(t *testing.T) {
	tests := []struct {
		name      string
		token     string
		mockSetup func(*mocks.Repository)
		want      *auth.Session
		wantErr   bool
	}{
		{
			name:  "success",
			token: "valid-token",
			mockSetup: func(r *mocks.Repository) {
				session := &auth.Session{
					AccessToken:  "valid-token",
					RefreshToken: "refresh-token",
					ExpiresAt:    time.Now().Add(time.Hour),
				}
				r.EXPECT().GetSessionByToken("valid-token").Return(session, nil)
			},
			wantErr: false,
		},
		{
			name:  "error - invalid token",
			token: "invalid-token",
			mockSetup: func(r *mocks.Repository) {
				r.EXPECT().GetSessionByToken("invalid-token").Return(nil, errors.New("session not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewRepository(ctrl)
			tt.mockSetup(repo)

			s := auth.NewService(repo)
			got, err := s.GetSession(tt.token)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
			ctrl.Finish()
		})
	}
}

func TestService_GetAllUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	tests := []struct {
		name      string
		mockSetup func(*mocks.Repository)
		want      []*auth.User
		wantErr   bool
	}{
		{
			name: "success",
			mockSetup: func(r *mocks.Repository) {
				users := []*auth.User{
					{ID: 1, Email: "user1@example.com"},
					{ID: 2, Email: "user2@example.com"},
				}
				r.EXPECT().GetAllUsers().Return(users, nil)
			},
			wantErr: false,
		},
		{
			name: "error",
			mockSetup: func(r *mocks.Repository) {
				r.EXPECT().GetAllUsers().Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewRepository(ctrl)
			tt.mockSetup(repo)

			s := auth.NewService(repo)
			got, err := s.GetAllUsers()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
			ctrl.Finish()
		})
	}
}
