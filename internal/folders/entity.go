package folders

import (
	"errors"
	"time"

	"gopkg.in/guregu/null.v4"
)

var (
	ErrNameRequired = errors.New("name is required")
)

type Folder struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	ParentID   null.Int  `json:"parent_id"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
	Deleted    bool      `json:"-"`
}

type FolderContent struct {
	Folder  Folder           `json:"folder"`
	Content []FolderResource `json:"content"`
}

type FolderResource struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

func New(name string, parentID null.Int) (*Folder, error) {
	f := Folder{
		Name: name,
	}

	if parentID.IsZero() {
		f.ParentID = parentID
	}

	err := f.Validate()
	if err != nil {
		return nil, err
	}

	return &f, nil
}

func (f *Folder) Validate() error {
	if f.Name == "" {
		return ErrNameRequired
	}

	return nil
}
