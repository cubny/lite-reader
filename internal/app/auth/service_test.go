package auth_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"

	"github.com/cubny/lite-reader/internal/app/auth"
	mocks "github.com/cubny/lite-reader/internal/mocks/app/auth"
)

const (
	tcSuccess    = "success"
	testEmail    = "test@example.com"
	testPassword = "password123"
)

func TestService_Signup(t *testing.T) {
	tests := []struct {
		name      string
		command   *auth.SignupCommand
		mockSetup func(*mocks.Repository)
		wantErr   bool
		errMsg    string
	}{
		{
			name: tcSuccess,
			command: &auth.SignupCommand{
				Email:    testEmail,
				Password: testPassword,
			},
			mockSetup: func(r *mocks.Repository) {
				r.EXPECT().GetSetting("allow_signup").Return("true", nil)
				r.EXPECT().GetUserByEmail(testEmail).Return(nil, errors.New("not found"))
				r.EXPECT().CreateUser(testEmail, gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error - email already exists",
			command: &auth.SignupCommand{
				Email:    "existing@example.com",
				Password: testPassword,
			},
			mockSetup: func(r *mocks.Repository) {
				r.EXPECT().GetSetting("allow_signup").Return("true", nil)
				r.EXPECT().GetUserByEmail("existing@example.com").Return(&auth.User{}, nil)
			},
			wantErr: true,
			errMsg:  "email already registered",
		},
		{
			name: "error - create user fails",
			command: &auth.SignupCommand{
				Email:    testEmail,
				Password: testPassword,
			},
			mockSetup: func(r *mocks.Repository) {
				r.EXPECT().GetSetting("allow_signup").Return("true", nil)
				r.EXPECT().GetUserByEmail(testEmail).Return(nil, errors.New("not found"))
				r.EXPECT().CreateUser(testEmail, gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
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
	tests := []struct {
		name      string
		command   *auth.LoginCommand
		mockSetup func(*mocks.Repository)
		want      *auth.LoginResponse
		wantErr   bool
		errMsg    string
	}{
		{
			name: tcSuccess,
			command: &auth.LoginCommand{
				Email:    testEmail,
				Password: testPassword,
			},
			mockSetup: func(r *mocks.Repository) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)
				user := &auth.User{
					ID:       1,
					Email:    testEmail,
					Password: string(hashedPassword),
				}
				session := &auth.Session{
					AccessToken:  "access-token",
					RefreshToken: "refresh-token",
					ExpiresAt:    time.Now().Add(time.Hour),
				}
				r.EXPECT().GetUserByEmail(testEmail).Return(user, nil)
				r.EXPECT().CreateSession(user.ID).Return(session, nil)
			},
			wantErr: false,
		},
		{
			name: "error - user not found",
			command: &auth.LoginCommand{
				Email:    "nonexistent@example.com",
				Password: testPassword,
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
				Email:    testEmail,
				Password: "wrongpassword",
			},
			mockSetup: func(r *mocks.Repository) {
				hashedPassword := "$2a$10$abcdefghijklmnopqrstuvwxyz" // pre-hashed testPassword
				user := &auth.User{
					ID:       1,
					Email:    testEmail,
					Password: hashedPassword,
				}
				r.EXPECT().GetUserByEmail(testEmail).Return(user, nil)
			},
			wantErr: true,
			errMsg:  "invalid email or password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
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
			name:  tcSuccess,
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
			name: tcSuccess,
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

func TestService_Setup(t *testing.T) {
	tests := []struct {
		name      string
		command   *auth.SetupCommand
		mockSetup func(*mocks.Repository)
		wantErr   bool
		errMsg    string
	}{
		{
			name: "success - creates admin and sets allow_signup to false",
			command: &auth.SetupCommand{
				Email:           testEmail,
				Password:        testPassword,
				ConfirmPassword: testPassword,
				AllowSignup:     false,
			},
			mockSetup: func(r *mocks.Repository) {
				r.EXPECT().CountUsers().Return(0, nil)
				r.EXPECT().CreateAdmin(testEmail, gomock.Any()).Return(nil)
				r.EXPECT().SetSetting("allow_signup", "false").Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success - creates admin and sets allow_signup to true",
			command: &auth.SetupCommand{
				Email:           testEmail,
				Password:        testPassword,
				ConfirmPassword: testPassword,
				AllowSignup:     true,
			},
			mockSetup: func(r *mocks.Repository) {
				r.EXPECT().CountUsers().Return(0, nil)
				r.EXPECT().CreateAdmin(testEmail, gomock.Any()).Return(nil)
				r.EXPECT().SetSetting("allow_signup", "true").Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error - users already exist",
			command: &auth.SetupCommand{
				Email:           testEmail,
				Password:        testPassword,
				ConfirmPassword: testPassword,
			},
			mockSetup: func(r *mocks.Repository) {
				r.EXPECT().CountUsers().Return(1, nil)
			},
			wantErr: true,
			errMsg:  "setup has already been completed",
		},
		{
			name: "error - count users fails",
			command: &auth.SetupCommand{
				Email:           testEmail,
				Password:        testPassword,
				ConfirmPassword: testPassword,
			},
			mockSetup: func(r *mocks.Repository) {
				r.EXPECT().CountUsers().Return(0, errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "db error",
		},
		{
			name: "error - create admin fails",
			command: &auth.SetupCommand{
				Email:           testEmail,
				Password:        testPassword,
				ConfirmPassword: testPassword,
			},
			mockSetup: func(r *mocks.Repository) {
				r.EXPECT().CountUsers().Return(0, nil)
				r.EXPECT().CreateAdmin(testEmail, gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewRepository(ctrl)
			tt.mockSetup(repo)

			s := auth.NewService(repo)
			err := s.Setup(tt.command)

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

func TestService_NeedsSetup(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*mocks.Repository)
		want      bool
		wantErr   bool
	}{
		{
			name: "needs setup - zero users",
			mockSetup: func(r *mocks.Repository) {
				r.EXPECT().CountUsers().Return(0, nil)
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "does not need setup - users exist",
			mockSetup: func(r *mocks.Repository) {
				r.EXPECT().CountUsers().Return(1, nil)
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "error counting users",
			mockSetup: func(r *mocks.Repository) {
				r.EXPECT().CountUsers().Return(0, errors.New("db error"))
			},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewRepository(ctrl)
			tt.mockSetup(repo)

			s := auth.NewService(repo)
			got, err := s.NeedsSetup()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			ctrl.Finish()
		})
	}
}

func TestService_IsSignupAllowed(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*mocks.Repository)
		want      bool
		wantErr   bool
	}{
		{
			name: "signup allowed",
			mockSetup: func(r *mocks.Repository) {
				r.EXPECT().GetSetting("allow_signup").Return("true", nil)
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "signup not allowed",
			mockSetup: func(r *mocks.Repository) {
				r.EXPECT().GetSetting("allow_signup").Return("false", nil)
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "error getting setting - defaults to false",
			mockSetup: func(r *mocks.Repository) {
				r.EXPECT().GetSetting("allow_signup").Return("", errors.New("not found"))
			},
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewRepository(ctrl)
			tt.mockSetup(repo)

			s := auth.NewService(repo)
			got := s.IsSignupAllowed()

			assert.Equal(t, tt.want, got)
			ctrl.Finish()
		})
	}
}

func TestService_SetAllowSignup(t *testing.T) {
	tests := []struct {
		name      string
		allow     bool
		mockSetup func(*mocks.Repository)
		wantErr   bool
	}{
		{
			name:  "set allow signup to true",
			allow: true,
			mockSetup: func(r *mocks.Repository) {
				r.EXPECT().SetSetting("allow_signup", "true").Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "set allow signup to false",
			allow: false,
			mockSetup: func(r *mocks.Repository) {
				r.EXPECT().SetSetting("allow_signup", "false").Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "error setting",
			allow: true,
			mockSetup: func(r *mocks.Repository) {
				r.EXPECT().SetSetting("allow_signup", "true").Return(errors.New("db error"))
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
			err := s.SetAllowSignup(tt.allow)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			ctrl.Finish()
		})
	}
}

func TestService_Signup_GuardedByAllowSignup(t *testing.T) {
	tests := []struct {
		name      string
		command   *auth.SignupCommand
		mockSetup func(*mocks.Repository)
		wantErr   bool
		errMsg    string
	}{
		{
			name: "signup allowed",
			command: &auth.SignupCommand{
				Email:    testEmail,
				Password: testPassword,
			},
			mockSetup: func(r *mocks.Repository) {
				r.EXPECT().GetSetting("allow_signup").Return("true", nil)
				r.EXPECT().GetUserByEmail(testEmail).Return(nil, errors.New("not found"))
				r.EXPECT().CreateUser(testEmail, gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "signup not allowed",
			command: &auth.SignupCommand{
				Email:    testEmail,
				Password: testPassword,
			},
			mockSetup: func(r *mocks.Repository) {
				r.EXPECT().GetSetting("allow_signup").Return("false", nil)
			},
			wantErr: true,
			errMsg:  "registration is currently disabled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
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
