package files

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/marcosvto1/go-driver/internal/queue"
	"gopkg.in/guregu/null.v4"
)

func (h *handler) Create(rw http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20) // 32mb

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	path := fmt.Sprintf("/%s", fileHeader.Filename)

	err = h.bucket.Upload(file, path)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	userId := r.Context().Value("user_id").(int64)
	entity, err := New(userId, fileHeader.Filename, fileHeader.Header.Get("Content-Type"), path)
	if err != nil {
		h.bucket.Delete(path)

		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	folderID := r.Form.Get("folder_id")
	if folderID != "" {
		fid, err := strconv.Atoi(folderID)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		entity.FolderId = null.NewInt(int64(fid), true)
	}

	id, err := InsertOne(h.db, entity)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	entity.ID = id

	dto := queue.QueueDTO{
		Filename: entity.Name,
		Path:     entity.Path,
		ID:       int(entity.ID),
	}

	msg, err := dto.Marshal()
	if err != nil {
		// TODO: ROLLBACK
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.queue.Publish(msg)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	json.NewEncoder(rw).Encode(entity)
}

func InsertOne(db *sql.DB, file *File) (id int64, err error) {
	file.ModifiedAt = time.Now()

	query := `INSERT INTO files (folder_id, owner_id, name, type, path, modified_at)
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	var lastInsertId int64

	err = db.QueryRow(query,
		file.FolderId,
		file.OwnerId,
		file.Name,
		file.Type,
		file.Path,
		file.ModifiedAt,
	).Scan(&lastInsertId)

	if err != nil {
		return -1, err
	}

	return lastInsertId, nil
}
