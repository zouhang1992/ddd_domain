package sqlite

import (
	"database/sql"
	"github.com/zouhang1992/ddd_domain/internal/domain/model"
	"github.com/zouhang1992/ddd_domain/internal/domain/repository"
)

// RoomRepository SQLite 房间仓储实现
type RoomRepository struct {
	conn *Connection
}

// NewRoomRepository 创建房间仓储
func NewRoomRepository(conn *Connection) repository.RoomRepository {
	return &RoomRepository{conn: conn}
}

// Save 保存房间
func (r *RoomRepository) Save(room *model.Room) error {
	_, err := r.conn.DB().Exec(`
		INSERT OR REPLACE INTO rooms (id, location_id, room_number, tags, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, room.ID, room.LocationID, room.RoomNumber, room.TagsString(), room.CreatedAt, room.UpdatedAt)
	return err
}

// FindByID 根据ID查找房间
func (r *RoomRepository) FindByID(id string) (*model.Room, error) {
	row := r.conn.DB().QueryRow(`
		SELECT id, location_id, room_number, tags, created_at, updated_at
		FROM rooms WHERE id = ?
	`, id)

	room := &model.Room{}
	var tagsStr string
	err := row.Scan(&room.ID, &room.LocationID, &room.RoomNumber, &tagsStr, &room.CreatedAt, &room.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	room.Tags = model.ParseTags(tagsStr)
	return room, nil
}

// FindAll 查找所有房间
func (r *RoomRepository) FindAll() ([]*model.Room, error) {
	rows, err := r.conn.DB().Query(`
		SELECT id, location_id, room_number, tags, created_at, updated_at
		FROM rooms ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []*model.Room
	for rows.Next() {
		room := &model.Room{}
		var tagsStr string
		err := rows.Scan(&room.ID, &room.LocationID, &room.RoomNumber, &tagsStr, &room.CreatedAt, &room.UpdatedAt)
		if err != nil {
			return nil, err
		}
		room.Tags = model.ParseTags(tagsStr)
		rooms = append(rooms, room)
	}
	return rooms, nil
}

// FindByLocationIDAndRoomNumber 根据位置ID和房间号查找房间
func (r *RoomRepository) FindByLocationIDAndRoomNumber(locationID, roomNumber string) (*model.Room, error) {
	row := r.conn.DB().QueryRow(`
		SELECT id, location_id, room_number, tags, created_at, updated_at
		FROM rooms WHERE location_id = ? AND room_number = ?
	`, locationID, roomNumber)

	room := &model.Room{}
	var tagsStr string
	err := row.Scan(&room.ID, &room.LocationID, &room.RoomNumber, &tagsStr, &room.CreatedAt, &room.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	room.Tags = model.ParseTags(tagsStr)
	return room, nil
}

// Delete 删除房间
func (r *RoomRepository) Delete(id string) error {
	_, err := r.conn.DB().Exec("DELETE FROM rooms WHERE id = ?", id)
	return err
}
