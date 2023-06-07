package files

import "database/sql"

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

	rows, err := db.Query(query)
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

	rows, err := db.Query(query, folderID)
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
