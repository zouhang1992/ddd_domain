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
