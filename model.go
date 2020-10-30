package mysql

// Model base model
type Model struct {
	ID        int64 `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
	CreatedAt int64 `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt int64 `gorm:"column:updated_at" json:"updatedAt"`
}
