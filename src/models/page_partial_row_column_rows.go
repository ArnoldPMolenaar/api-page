package models

type PagePartialRowColumnRow struct {
	ColumnID uint `gorm:"column:column_id"`
	RowID    uint `gorm:"column:row_id"`
}

func (PagePartialRowColumnRow) TableName() string {
	return "page_partial_row_column_rows"
}
