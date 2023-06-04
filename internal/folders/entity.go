package folders

import (
	"errors"
	"time"
)

var (
	ErrNameRequired = errors.New("name is required")
)

type Folder struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	ParentID   int64     `json:"parent_id"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
	Deleted    bool      `json:"-"`
}

func New(name string, parentId int64) (*Folder, error) {
	f := Folder{
		Name:     name,
		ParentID: parentId,
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
