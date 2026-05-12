package folder

type Repository interface {
	AddFolder(folder *Folder) (int, error)
	GetFolder(id int) (*Folder, error)
	ListFolders(userID int) ([]*Folder, error)
	RenameFolder(id int, name string) error
	ReorderFolder(id, position int) error
	DeleteFolder(id int) error
}
