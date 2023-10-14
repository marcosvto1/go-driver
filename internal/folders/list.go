package folders

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/marcosvto1/go-driver/internal/files"
)

func (h *handler) List(rw http.ResponseWriter, r *http.Request) {
	c, err := GetRootFolderContent(h.db)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	fc := FolderContent{
		Folder: Folder{
			Name: "root",
		},
		Content: c,
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(fc)
}

func getRootSubFolder(db *sql.DB) ([]Folder, error) {
	query := `SELECT
	id,
	name,
	parent_id,
	created_at,
	modified_at,
	deleted
	FROM "folders" WHERE "parent_id" IS NULL "deleted"=false`

	rows, err := db.Query(query)
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

func GetRootFolderContent(db *sql.DB) ([]FolderResource, error) {
	listSubfolder, err := getRootSubFolder(db)
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

	folderFiles, err := files.ListRootFiles(db)
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
