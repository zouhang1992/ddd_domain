package sqlite

import (
	"database/sql"
	"github.com/zouhang1992/ddd_domain/internal/domain/model"
	"github.com/zouhang1992/ddd_domain/internal/domain/repository"
)

// LandlordRepository SQLite 房东仓储实现
type LandlordRepository struct {
	conn *Connection
}

// NewLandlordRepository 创建房东仓储
func NewLandlordRepository(conn *Connection) repository.LandlordRepository {
	return &LandlordRepository{conn: conn}
}

// Save 保存房东
func (r *LandlordRepository) Save(landlord *model.Landlord) error {
	_, err := r.conn.DB().Exec(`
		INSERT OR REPLACE INTO landlords (id, name, phone, note, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, landlord.ID, landlord.Name, landlord.Phone, landlord.Note, landlord.CreatedAt, landlord.UpdatedAt)
	return err
}

// FindByID 根据ID查找房东
func (r *LandlordRepository) FindByID(id string) (*model.Landlord, error) {
	row := r.conn.DB().QueryRow(`
		SELECT id, name, phone, note, created_at, updated_at
		FROM landlords WHERE id = ?
	`, id)

	landlord := &model.Landlord{}
	err := row.Scan(&landlord.ID, &landlord.Name, &landlord.Phone, &landlord.Note,
		&landlord.CreatedAt, &landlord.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return landlord, nil
}

// FindAll 查找所有房东
func (r *LandlordRepository) FindAll() ([]*model.Landlord, error) {
	rows, err := r.conn.DB().Query(`
		SELECT id, name, phone, note, created_at, updated_at
		FROM landlords ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var landlords []*model.Landlord
	for rows.Next() {
		landlord := &model.Landlord{}
		err := rows.Scan(&landlord.ID, &landlord.Name, &landlord.Phone, &landlord.Note,
			&landlord.CreatedAt, &landlord.UpdatedAt)
		if err != nil {
			return nil, err
		}
		landlords = append(landlords, landlord)
	}
	return landlords, nil
}

// FindByCriteria 按条件查找房东
func (r *LandlordRepository) FindByCriteria(criteria repository.LandlordCriteria, offset, limit int) ([]*model.Landlord, error) {
	query := `
		SELECT id, name, phone, note, created_at, updated_at
		FROM landlords
		WHERE 1 = 1
	`
	var args []interface{}

	if criteria.Name != "" {
		query += " AND name LIKE ?"
		args = append(args, "%"+criteria.Name+"%")
	}
	if criteria.Phone != "" {
		query += " AND phone LIKE ?"
		args = append(args, "%"+criteria.Phone+"%")
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

	var landlords []*model.Landlord
	for rows.Next() {
		landlord := &model.Landlord{}
		err := rows.Scan(&landlord.ID, &landlord.Name, &landlord.Phone, &landlord.Note,
			&landlord.CreatedAt, &landlord.UpdatedAt)
		if err != nil {
			return nil, err
		}
		landlords = append(landlords, landlord)
	}
	return landlords, nil
}

// CountByCriteria 按条件统计房东数量
func (r *LandlordRepository) CountByCriteria(criteria repository.LandlordCriteria) (int, error) {
	query := `
		SELECT COUNT(*) FROM landlords
		WHERE 1 = 1
	`
	var args []interface{}

	if criteria.Name != "" {
		query += " AND name LIKE ?"
		args = append(args, "%"+criteria.Name+"%")
	}
	if criteria.Phone != "" {
		query += " AND phone LIKE ?"
		args = append(args, "%"+criteria.Phone+"%")
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

// Delete 删除房东
func (r *LandlordRepository) Delete(id string) error {
	_, err := r.conn.DB().Exec("DELETE FROM landlords WHERE id = ?", id)
	return err
}

// HasLeases 检查是否有关联租约
func (r *LandlordRepository) HasLeases(landlordID string) (bool, error) {
	var count int
	row := r.conn.DB().QueryRow("SELECT COUNT(*) FROM leases WHERE landlord_id = ?", landlordID)
	err := row.Scan(&count)
	return count > 0, err
}
