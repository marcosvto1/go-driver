package files

import (
	"database/sql"
)

func ListRootFiles(db *sql.DB) ([]File, error) {
	query := `SELECT
		id,
		name,
		folder_id,
		owner_id,
		type,
		path,
		created_at,
		modified_at,
		deleted
		FROM files
		WHERE folder_id IS NULL AND deleted = false`
	return selectAllFiles(db, query, -1)
}

func List(db *sql.DB, folderID int64) ([]File, error) {
	query := `SELECT
		id,
		name,
		folder_id,
		owner_id,
		type,
		path,
		created_at,
		modified_at,
		deleted
		FROM files
		WHERE folder_id = $1 AND deleted = false`

	return selectAllFiles(db, query, folderID)
}

func selectAllFiles(db *sql.DB, query string, folderID int64) ([]File, error) {
	var rows *sql.Rows
	var err error

	if folderID != -1 {
		rs, er := db.Query(query, folderID)
		rows = rs
		err = er
	} else {
		rs, er := db.Query(query)
		rows = rs
		err = er
	}

	if err != nil {
		return nil, err
	}

	files := make([]File, 0)

	for rows.Next() {
		var file File
		err := rows.Scan(
			&file.ID,
			&file.Name,
			&file.FolderId,
			&file.OwnerId,
			&file.Type,
			&file.Path,
			&file.CreatedAt,
			&file.ModifiedAt,
			&file.Deleted,
		)

		if err != nil {
			continue
		}

		files = append(files, file)
	}

	return files, nil
}
