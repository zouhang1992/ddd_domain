package mysql

import (
	"database/sql"
	"time"

	landlordmodel "github.com/zouhang1992/ddd_domain/internal/domain/landlord/model"
	landlordrepo "github.com/zouhang1992/ddd_domain/internal/domain/landlord/repository"
)

// LandlordRepository MySQL 房东仓储实现
type LandlordRepository struct {
	conn *Connection
}

// NewLandlordRepository 创建房东仓储
func NewLandlordRepository(conn *Connection) landlordrepo.LandlordRepository {
	return &LandlordRepository{conn: conn}
}

// tempLandlord is a temporary struct for scanning
type tempLandlord struct {
	ID        string
	Name      string
	Phone     string
	Note      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Save 保存房东
func (r *LandlordRepository) Save(landlord *landlordmodel.Landlord) error {
	_, err := r.conn.DB().Exec(`
		INSERT INTO landlords (id, name, phone, note, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			name = VALUES(name),
			phone = VALUES(phone),
			note = VALUES(note),
			updated_at = VALUES(updated_at)
	`, landlord.IDField, landlord.Name, landlord.Phone, landlord.Note, landlord.CreatedAt, landlord.UpdatedAt)
	return err
}

// FindByID 根据ID查找房东
func (r *LandlordRepository) FindByID(id string) (*landlordmodel.Landlord, error) {
	row := r.conn.DB().QueryRow(`
		SELECT id, name, phone, note, created_at, updated_at
		FROM landlords WHERE id = ?
	`, id)

	var temp tempLandlord
	err := row.Scan(&temp.ID, &temp.Name, &temp.Phone, &temp.Note, &temp.CreatedAt, &temp.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	landlord := landlordmodel.NewLandlord(temp.ID, temp.Name, temp.Phone, temp.Note)
	landlord.CreatedAt = temp.CreatedAt
	landlord.UpdatedAt = temp.UpdatedAt
	landlord.ClearEvents()

	return landlord, nil
}

// FindAll 查找所有房东
func (r *LandlordRepository) FindAll() ([]*landlordmodel.Landlord, error) {
	rows, err := r.conn.DB().Query(`
		SELECT id, name, phone, note, created_at, updated_at
		FROM landlords ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var landlords []*landlordmodel.Landlord
	for rows.Next() {
		var temp tempLandlord
		err := rows.Scan(&temp.ID, &temp.Name, &temp.Phone, &temp.Note, &temp.CreatedAt, &temp.UpdatedAt)
		if err != nil {
			return nil, err
		}

		landlord := landlordmodel.NewLandlord(temp.ID, temp.Name, temp.Phone, temp.Note)
		landlord.CreatedAt = temp.CreatedAt
		landlord.UpdatedAt = temp.UpdatedAt

		landlords = append(landlords, landlord)
	}
	return landlords, nil
}

// Delete 删除房东
func (r *LandlordRepository) Delete(id string) error {
	_, err := r.conn.DB().Exec("DELETE FROM landlords WHERE id = ?", id)
	return err
}
