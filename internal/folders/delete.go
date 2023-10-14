package folders

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/marcosvto1/go-driver/internal/files"
)

func (h *handler) Delete(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = deleteFolderContent(h.db, int64(id))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = Delete(h.db, int64(id))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
}

func deleteFolderContent(db *sql.DB, folderID int64) error {
	err := deleteFiles(db, folderID)
	if err != nil {
		return err
	}

	err = deleteSubFolders(db, folderID)
	if err != nil {
		return err
	}

	return nil
}

func deleteSubFolders(db *sql.DB, folderId int64) error {
	listOfSubfolder, err := getSubFolder(db, folderId)
	if err != nil {
		return err
	}

	removedSubFolders := make([]Folder, 0, len(listOfSubfolder))
	for _, sf := range listOfSubfolder {
		err := Delete(db, sf.ID)
		if err != nil {
			break
		}

		err = deleteFolderContent(db, sf.ID)
		if err != nil {
			Update(db, sf.ID, &sf)
			break
		}

		removedSubFolders = append(removedSubFolders, sf)
	}

	if len(removedSubFolders) != len(listOfSubfolder) {
		for _, rf := range removedSubFolders {
			err = Delete(db, rf.ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func deleteFiles(db *sql.DB, folderId int64) error {
	listOfFiles, err := files.List(db, folderId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	removedFiles := make([]files.File, 0, len(listOfFiles))
	for _, file := range listOfFiles {
		fmt.Println("ok")
		file.Deleted = true
		err := files.Update(db, file.ID, &file)
		if err != nil {
			break
		}

		removedFiles = append(removedFiles, file)
	}

	if len(removedFiles) != len(listOfFiles) {
		for _, file := range removedFiles {
			file.Deleted = true
			err := files.Update(db, file.ID, &file)
			return err
		}
	}

	return nil
}

func Delete(db *sql.DB, id int64) error {
	query := `UPDATE "folders" SET "modified_at"=$1, "deleted"=true WHERE id=$2`

	_, err := db.Exec(query, time.Now(), id)

	return err
}
