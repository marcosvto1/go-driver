package folders

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/marcosvto1/go-driver/internal/files"
)

func (h *handler) Get(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	folder, err := GetFolder(h.db, int64(id))
	if err != nil {
		fmt.Println(err)

		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	c, err := GetFolderContent(h.db, int64(id))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	fc := FolderContent{Folder: *folder, Content: c}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(fc)
}

func GetFolder(db *sql.DB, folderId int64) (*Folder, error) {
	fmt.Println(folderId)

	query := `SELECT
	id,
	name,
	parent_id,
	created_at,
	modified_at,
	deleted
	FROM "folders" where id=$1`

	fmt.Println(query)

	row := db.QueryRow(query, folderId)

	folder := new(Folder)
	err := row.Scan(
		&folder.ID,
		&folder.Name,
		&folder.ParentID,
		&folder.CreatedAt,
		&folder.ModifiedAt,
		&folder.Deleted,
	)
	if err != nil {
		return nil, err
	}

	return folder, nil
}

func getSubFolder(db *sql.DB, folderId int64) ([]Folder, error) {
	query := `SELECT
	id,
	name,
	parent_id,
	created_at,
	modified_at,
	deleted
	FROM "folders" where parent_id=$1`

	rows, err := db.Query(query, folderId)
	if err != nil {
		return nil, err
	}

	folders := make([]Folder, 0)
	for rows.Next() {
		var folder Folder
		err := rows.Scan(
			&folder.ID,
			&folder.Name,
			&folder.ParentID,
			&folder.CreatedAt,
			&folder.ModifiedAt,
			&folder.Deleted,
		)

		if err != nil {
			continue
		}

		folders = append(folders, folder)
	}

	return folders, nil
}

func GetFolderContent(db *sql.DB, folderId int64) ([]FolderResource, error) {
	listSubfolder, err := getSubFolder(db, folderId)
	if err != nil {
		return nil, err
	}

	fr := make([]FolderResource, 0, len(listSubfolder))
	for _, sf := range listSubfolder {
		r := FolderResource{
			ID:         sf.ID,
			Name:       sf.Name,
			Type:       "directory",
			CreatedAt:  sf.CreatedAt,
			ModifiedAt: sf.ModifiedAt,
		}

		fr = append(fr, r)
	}

	folderFiles, err := files.List(db, folderId)
	if err != nil {
		return nil, err
	}

	for _, ff := range folderFiles {
		r := FolderResource{
			ID:         ff.ID,
			Name:       ff.Name,
			Type:       ff.Type,
			CreatedAt:  ff.CreatedAt,
			ModifiedAt: ff.ModifiedAt,
		}

		fr = append(fr, r)
	}

	return fr, nil
}
