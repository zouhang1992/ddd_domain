package mysql

import (
	"database/sql"
	"time"

	locationmodel "github.com/zouhang1992/ddd_domain/internal/domain/location/model"
	locationrepo "github.com/zouhang1992/ddd_domain/internal/domain/location/repository"
)

// LocationRepository MySQL 位置仓储实现
type LocationRepository struct {
	conn *Connection
}

// NewLocationRepository 创建位置仓储
func NewLocationRepository(conn *Connection) locationrepo.LocationRepository {
	return &LocationRepository{conn: conn}
}

// tempLocation is a temporary struct for scanning
type tempLocation struct {
	ID        string
	ShortName string
	Detail    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Save 保存位置
func (r *LocationRepository) Save(location *locationmodel.Location) error {
	_, err := r.conn.DB().Exec(`
		INSERT INTO locations (id, short_name, detail, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			short_name = VALUES(short_name),
			detail = VALUES(detail),
			updated_at = VALUES(updated_at)
	`, location.IDField, location.ShortName, location.Detail, location.CreatedAt, location.UpdatedAt)
	return err
}

// FindByID 根据ID查找位置
func (r *LocationRepository) FindByID(id string) (*locationmodel.Location, error) {
	row := r.conn.DB().QueryRow(`
		SELECT id, short_name, detail, created_at, updated_at
		FROM locations WHERE id = ?
	`, id)

	var temp tempLocation
	err := row.Scan(&temp.ID, &temp.ShortName, &temp.Detail, &temp.CreatedAt, &temp.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	location := locationmodel.NewLocation(temp.ID, temp.ShortName, temp.Detail)
	location.CreatedAt = temp.CreatedAt
	location.UpdatedAt = temp.UpdatedAt
	location.ClearEvents()

	return location, nil
}

// FindAll 查找所有位置
func (r *LocationRepository) FindAll() ([]*locationmodel.Location, error) {
	rows, err := r.conn.DB().Query(`
		SELECT id, short_name, detail, created_at, updated_at
		FROM locations ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []*locationmodel.Location
	for rows.Next() {
		var temp tempLocation
		err := rows.Scan(&temp.ID, &temp.ShortName, &temp.Detail, &temp.CreatedAt, &temp.UpdatedAt)
		if err != nil {
			return nil, err
		}

		location := locationmodel.NewLocation(temp.ID, temp.ShortName, temp.Detail)
		location.CreatedAt = temp.CreatedAt
		location.UpdatedAt = temp.UpdatedAt

		locations = append(locations, location)
	}
	return locations, nil
}

// Delete 删除位置
func (r *LocationRepository) Delete(id string) error {
	_, err := r.conn.DB().Exec("DELETE FROM locations WHERE id = ?", id)
	return err
}
