package organisations

type Organisation struct {
	Id   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"unique"`
	Head string `json:"head"`
}
