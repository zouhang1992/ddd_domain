package sqlite

import (
	"database/sql"
	"strings"
	"time"

	roommodel "github.com/zouhang1992/ddd_domain/internal/domain/room/model"
	roomrepo "github.com/zouhang1992/ddd_domain/internal/domain/room/repository"
)

// RoomRepository SQLite 房间仓储实现
type RoomRepository struct {
	conn *Connection
}

// NewRoomRepository 创建房间仓储
func NewRoomRepository(conn *Connection) roomrepo.RoomRepository {
	return &RoomRepository{conn: conn}
}

// tempRoom is a temporary struct for scanning
type tempRoom struct {
	ID         string
	LocationID string
	RoomNumber string
	Status     sql.NullString
	Tags       sql.NullString
	Note       sql.NullString
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Save 保存房间
func (r *RoomRepository) Save(room *roommodel.Room) error {
	tagsStr := strings.Join(room.Tags, ",")
	_, err := r.conn.DB().Exec(`
		INSERT OR REPLACE INTO rooms (id, location_id, room_number, status, tags, note, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, room.IDField, room.LocationID, room.RoomNumber, string(room.Status), tagsStr, room.Note, room.CreatedAt, room.UpdatedAt)
	return err
}

// FindByID 根据ID查找房间
func (r *RoomRepository) FindByID(id string) (*roommodel.Room, error) {
	row := r.conn.DB().QueryRow(`
		SELECT id, location_id, room_number, status, tags, note, created_at, updated_at
		FROM rooms WHERE id = ?
	`, id)

	var temp tempRoom
	err := row.Scan(&temp.ID, &temp.LocationID, &temp.RoomNumber, &temp.Status, &temp.Tags, &temp.Note, &temp.CreatedAt, &temp.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var tags []string
	if temp.Tags.Valid && temp.Tags.String != "" {
		tags = strings.Split(temp.Tags.String, ",")
	} else {
		tags = []string{}
	}

	note := ""
	if temp.Note.Valid {
		note = temp.Note.String
	}

	status := roommodel.RoomStatusAvailable
	if temp.Status.Valid && temp.Status.String != "" {
		status = roommodel.RoomStatus(temp.Status.String)
	}

	room := roommodel.NewRoom(temp.ID, temp.LocationID, temp.RoomNumber, tags, note)
	room.Status = status
	room.CreatedAt = temp.CreatedAt
	room.UpdatedAt = temp.UpdatedAt

	return room, nil
}

// FindAll 查找所有房间
func (r *RoomRepository) FindAll() ([]*roommodel.Room, error) {
	rows, err := r.conn.DB().Query(`
		SELECT id, location_id, room_number, status, tags, note, created_at, updated_at
		FROM rooms ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []*roommodel.Room
	for rows.Next() {
		var temp tempRoom
		err := rows.Scan(&temp.ID, &temp.LocationID, &temp.RoomNumber, &temp.Status, &temp.Tags, &temp.Note, &temp.CreatedAt, &temp.UpdatedAt)
		if err != nil {
			return nil, err
		}

		var tags []string
		if temp.Tags.Valid && temp.Tags.String != "" {
			tags = strings.Split(temp.Tags.String, ",")
		} else {
			tags = []string{}
		}

		note := ""
		if temp.Note.Valid {
			note = temp.Note.String
		}

		status := roommodel.RoomStatusAvailable
		if temp.Status.Valid && temp.Status.String != "" {
			status = roommodel.RoomStatus(temp.Status.String)
		}

		room := roommodel.NewRoom(temp.ID, temp.LocationID, temp.RoomNumber, tags, note)
		room.Status = status
		room.CreatedAt = temp.CreatedAt
		room.UpdatedAt = temp.UpdatedAt

		rooms = append(rooms, room)
	}
	return rooms, nil
}

// Delete 删除房间
func (r *RoomRepository) Delete(id string) error {
	_, err := r.conn.DB().Exec("DELETE FROM rooms WHERE id = ?", id)
	return err
}
