package files

import (
	"errors"
	"time"
)

var (
	ErrNameRequired  = errors.New("name is required")
	ErrOwnerRequired = errors.New("owner is required")
	ErrTypoRequired  = errors.New("type is required")
	ErrPathRequired  = errors.New("path is required")
)

type File struct {
	ID         int64     `json:"id"`
	FolderId   int64     `json:"-"`
	OwnerId    int64     `json:"owner_id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Path       string    `json:"-"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
	Deleted    bool      `json:"deleted"`
}

func New(ownerID int64, name, fileType, path string) (*File, error) {
	f := File{
		OwnerId:    ownerID,
		Name:       name,
		Type:       fileType,
		Path:       path,
		ModifiedAt: time.Now(),
	}

	err := f.Validate()
	if err != nil {
		return nil, err
	}

	return &f, nil
}

func (f *File) Validate() error {
	if f.Name == "" {
		return ErrNameRequired
	}

	if f.OwnerId == 0 {
		return ErrOwnerRequired
	}

	if f.Path == "" {
		return ErrPathRequired
	}

	if f.Type == "" {
		return ErrTypoRequired
	}

	return nil
}
