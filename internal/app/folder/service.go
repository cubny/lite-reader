package folder

import (
	"errors"
	"fmt"
	"strings"
)

const maxNameLen = 100

var (
	ErrEmptyName   = errors.New("folder name cannot be empty")
	ErrNameTooLong = errors.New("folder name is too long")
	ErrNotFound    = errors.New("folder not found")
	ErrNotOwner    = errors.New("folder does not belong to user")
)

type ServiceImpl struct {
	repository Repository
}

func NewService(repository Repository) *ServiceImpl {
	return &ServiceImpl{repository: repository}
}

func (s *ServiceImpl) AddFolder(command *AddFolderCommand) (*Folder, error) {
	name, err := normalizeName(command.Name)
	if err != nil {
		return nil, err
	}
	f := &Folder{
		Name:   name,
		UserID: command.UserID,
	}
	id, err := s.repository.AddFolder(f)
	if err != nil {
		return nil, fmt.Errorf("cannot add folder: %w", err)
	}
	f.ID = id
	return f, nil
}

func (s *ServiceImpl) ListFolders(userID int) ([]*Folder, error) {
	return s.repository.ListFolders(userID)
}

func (s *ServiceImpl) RenameFolder(command *RenameFolderCommand) error {
	name, err := normalizeName(command.Name)
	if err != nil {
		return err
	}
	if err := s.assertOwner(command.FolderID, command.UserID); err != nil {
		return err
	}
	return s.repository.RenameFolder(command.FolderID, name)
}

func (s *ServiceImpl) ReorderFolder(command *ReorderFolderCommand) error {
	if err := s.assertOwner(command.FolderID, command.UserID); err != nil {
		return err
	}
	return s.repository.ReorderFolder(command.FolderID, command.Position)
}

func (s *ServiceImpl) DeleteFolder(command *DeleteFolderCommand) error {
	if err := s.assertOwner(command.FolderID, command.UserID); err != nil {
		return err
	}
	return s.repository.DeleteFolder(command.FolderID)
}

func (s *ServiceImpl) GetFolder(command *GetFolderCommand) (*Folder, error) {
	f, err := s.repository.GetFolder(command.FolderID)
	if err != nil {
		return nil, err
	}
	if f == nil {
		return nil, ErrNotFound
	}
	if f.UserID != command.UserID {
		return nil, ErrNotOwner
	}
	return f, nil
}

func (s *ServiceImpl) assertOwner(folderID, userID int) error {
	f, err := s.repository.GetFolder(folderID)
	if err != nil {
		return err
	}
	if f == nil {
		return ErrNotFound
	}
	if f.UserID != userID {
		return ErrNotOwner
	}
	return nil
}

func normalizeName(name string) (string, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "", ErrEmptyName
	}
	if len(trimmed) > maxNameLen {
		return "", ErrNameTooLong
	}
	return trimmed, nil
}
