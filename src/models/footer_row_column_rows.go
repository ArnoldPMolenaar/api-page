package models

type FooterRowColumnRow struct {
	ColumnID uint `gorm:"column:column_id"`
	RowID    uint `gorm:"column:row_id"`
}

func (FooterRowColumnRow) TableName() string {
	return "footer_row_column_rows"
}
