package models

type Good struct {
	ID    int `gorm:"primary_key";"AUTO_INCREMENT"`
	Name  string
	Price int
}

type File struct {
	Name        string
	ContentType string
	Content     []byte
}
