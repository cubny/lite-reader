package folder

type AddFolderCommand struct {
	Name   string
	UserID int
}

type RenameFolderCommand struct {
	FolderID int
	UserID   int
	Name     string
}

type ReorderFolderCommand struct {
	FolderID int
	UserID   int
	Position int
}

type DeleteFolderCommand struct {
	FolderID int
	UserID   int
}

type GetFolderCommand struct {
	FolderID int
	UserID   int
}
