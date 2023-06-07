package files

import "database/sql"

func Get(db *sql.DB, id int64) (*File, error) {
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
	FROM "files" WHERE "id"=$1`

	row := db.QueryRow(query, id)

	var file File
	err := row.Scan(
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
		return nil, err
	}

	return &file, nil
}
