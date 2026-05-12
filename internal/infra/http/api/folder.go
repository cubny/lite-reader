package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"

	"github.com/cubny/lite-reader/internal/app/folder"
	"github.com/cubny/lite-reader/internal/app/item"
	"github.com/cubny/lite-reader/internal/infra/http/api/cxutil"
)

func (h *Router) listFolders(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID := r.Context().Value(cxutil.UserIDKey).(int)
	folders, err := h.folderService.ListFolders(userID)
	if err != nil {
		log.WithError(err).Error("listFolders: service")
		_ = InternalError(w, "cannot list folders")
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(toListFolderResponse(folders)); err != nil {
		log.WithError(err).Error("listFolders: encoder")
		_ = InternalError(w, "cannot encode response")
	}
}

func (h *Router) addFolder(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	command, err := toAddFolderCommand(w, r, p)
	if err != nil {
		return
	}

	f, err := h.folderService.AddFolder(command)
	if err != nil {
		if errors.Is(err, folder.ErrEmptyName) || errors.Is(err, folder.ErrNameTooLong) {
			_ = InvalidParams(w, err.Error())
			return
		}
		log.WithError(err).Error("addFolder: service")
		_ = InternalError(w, "cannot add folder")
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(toFolderResponse(f)); err != nil {
		log.WithError(err).Error("addFolder: encoder")
		_ = InternalError(w, "cannot encode response")
	}
}

func (h *Router) updateFolder(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	folderID, userID, ok := folderIDAndUser(w, r, p)
	if !ok {
		return
	}

	request := &UpdateFolderRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		_ = BadRequest(w, "cannot decode request body")
		return
	}

	if request.Name != nil {
		err := h.folderService.RenameFolder(&folder.RenameFolderCommand{
			FolderID: folderID, UserID: userID, Name: *request.Name,
		})
		if err != nil {
			if handleFolderServiceError(w, err) {
				return
			}
			log.WithError(err).Error("updateFolder: rename")
			_ = InternalError(w, "cannot rename folder")
			return
		}
	}

	if request.Position != nil {
		err := h.folderService.ReorderFolder(&folder.ReorderFolderCommand{
			FolderID: folderID, UserID: userID, Position: *request.Position,
		})
		if err != nil {
			if handleFolderServiceError(w, err) {
				return
			}
			log.WithError(err).Error("updateFolder: reorder")
			_ = InternalError(w, "cannot reorder folder")
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Router) deleteFolder(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	folderID, userID, ok := folderIDAndUser(w, r, p)
	if !ok {
		return
	}
	err := h.folderService.DeleteFolder(&folder.DeleteFolderCommand{FolderID: folderID, UserID: userID})
	if err != nil {
		if handleFolderServiceError(w, err) {
			return
		}
		log.WithError(err).Error("deleteFolder: service")
		_ = InternalError(w, "cannot delete folder")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Router) getFolderItems(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	folderID, userID, ok := folderIDAndUser(w, r, p)
	if !ok {
		return
	}
	items, err := h.itemService.GetFolderItems(&item.GetFolderItemsCommand{FolderID: folderID, UserID: userID})
	if err != nil {
		log.WithError(err).Error("getFolderItems: service")
		_ = InternalError(w, "cannot get folder items")
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(toGetItemsResponse(items)); err != nil {
		log.WithError(err).Error("getFolderItems: encoder")
		_ = InternalError(w, "cannot encode response")
	}
}

func (h *Router) readFolderItems(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	folderID, userID, ok := folderIDAndUser(w, r, p)
	if !ok {
		return
	}
	if err := h.itemService.ReadFolderItems(&item.ReadFolderItemsCommand{FolderID: folderID, UserID: userID}); err != nil {
		log.WithError(err).Error("readFolderItems: service")
		_ = InternalError(w, "cannot mark folder read")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func handleFolderServiceError(w http.ResponseWriter, err error) bool {
	switch {
	case errors.Is(err, folder.ErrEmptyName), errors.Is(err, folder.ErrNameTooLong):
		_ = InvalidParams(w, err.Error())
		return true
	case errors.Is(err, folder.ErrNotFound):
		_ = NotFound(w, err.Error())
		return true
	case errors.Is(err, folder.ErrNotOwner):
		_ = NotFound(w, "folder not found")
		return true
	}
	return false
}
