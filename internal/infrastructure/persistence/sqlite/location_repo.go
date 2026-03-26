package sqlite

import (
	"database/sql"
	"github.com/zouhang1992/ddd_domain/internal/domain/model"
	"github.com/zouhang1992/ddd_domain/internal/domain/repository"
)

// LocationRepository SQLite 位置仓储实现
type LocationRepository struct {
	conn *Connection
}

// NewLocationRepository 创建位置仓储
func NewLocationRepository(conn *Connection) repository.LocationRepository {
	return &LocationRepository{conn: conn}
}

// Save 保存位置
func (r *LocationRepository) Save(location *model.Location) error {
	_, err := r.conn.DB().Exec(`
		INSERT OR REPLACE INTO locations (id, short_name, detail, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`, location.ID, location.ShortName, location.Detail, location.CreatedAt, location.UpdatedAt)
	return err
}

// FindByID 根据ID查找位置
func (r *LocationRepository) FindByID(id string) (*model.Location, error) {
	row := r.conn.DB().QueryRow(`
		SELECT id, short_name, detail, created_at, updated_at
		FROM locations WHERE id = ?
	`, id)

	loc := &model.Location{}
	err := row.Scan(&loc.ID, &loc.ShortName, &loc.Detail, &loc.CreatedAt, &loc.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return loc, nil
}

// FindAll 查找所有位置
func (r *LocationRepository) FindAll() ([]*model.Location, error) {
	rows, err := r.conn.DB().Query(`
		SELECT id, short_name, detail, created_at, updated_at
		FROM locations ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []*model.Location
	for rows.Next() {
		loc := &model.Location{}
		err := rows.Scan(&loc.ID, &loc.ShortName, &loc.Detail, &loc.CreatedAt, &loc.UpdatedAt)
		if err != nil {
			return nil, err
		}
		locations = append(locations, loc)
	}
	return locations, nil
}

// Delete 删除位置
func (r *LocationRepository) Delete(id string) error {
	_, err := r.conn.DB().Exec("DELETE FROM locations WHERE id = ?", id)
	return err
}

// HasRooms 检查位置是否有关联房间
func (r *LocationRepository) HasRooms(locationID string) (bool, error) {
	var count int
	err := r.conn.DB().QueryRow("SELECT COUNT(*) FROM rooms WHERE location_id = ?", locationID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindByCriteria 按条件查找位置
func (r *LocationRepository) FindByCriteria(criteria repository.LocationCriteria, offset, limit int) ([]*model.Location, error) {
	query := `
		SELECT id, short_name, detail, created_at, updated_at
		FROM locations
		WHERE 1 = 1
	`
	var args []interface{}

	if criteria.ShortName != "" {
		query += " AND short_name LIKE ?"
		args = append(args, "%"+criteria.ShortName+"%")
	}
	if criteria.Detail != "" {
		query += " AND detail LIKE ?"
		args = append(args, "%"+criteria.Detail+"%")
	}
	if criteria.StartTime != nil {
		query += " AND created_at >= ?"
		args = append(args, criteria.StartTime)
	}
	if criteria.EndTime != nil {
		query += " AND created_at <= ?"
		args = append(args, criteria.EndTime)
	}

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.conn.DB().Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []*model.Location
	for rows.Next() {
		loc := &model.Location{}
		err := rows.Scan(&loc.ID, &loc.ShortName, &loc.Detail, &loc.CreatedAt, &loc.UpdatedAt)
		if err != nil {
			return nil, err
		}
		locations = append(locations, loc)
	}
	return locations, nil
}

// CountByCriteria 按条件统计位置数量
func (r *LocationRepository) CountByCriteria(criteria repository.LocationCriteria) (int, error) {
	query := `
		SELECT COUNT(*) FROM locations
		WHERE 1 = 1
	`
	var args []interface{}

	if criteria.ShortName != "" {
		query += " AND short_name LIKE ?"
		args = append(args, "%"+criteria.ShortName+"%")
	}
	if criteria.Detail != "" {
		query += " AND detail LIKE ?"
		args = append(args, "%"+criteria.Detail+"%")
	}
	if criteria.StartTime != nil {
		query += " AND created_at >= ?"
		args = append(args, criteria.StartTime)
	}
	if criteria.EndTime != nil {
		query += " AND created_at <= ?"
		args = append(args, criteria.EndTime)
	}

	var count int
	row := r.conn.DB().QueryRow(query, args...)
	err := row.Scan(&count)
	return count, err
}
