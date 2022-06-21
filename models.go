package main

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	db, err := gorm.Open(sqlite.Open("db.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(new(File))
	DB = db
}

type File struct {
	ID        uint `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	Name      string `json:"name" gorm:"column:name;type:varchar;size:255"`
	Size      int64  `json:"size" gorm:"column:size;type:bigint"`
	Path      string `json:"path" gorm:"column:path;type:varchar;size:255;index"`
	Srouce    string `json:"source" gorm:"column:srouce;type:varchar;size:255"`
}

func (t *File) GetFileAll() (all []File, err error) {
	err = DB.Find(&all).Error
	return
}

func (t *File) Delete(id uint) (err error) {
	return DB.Delete(t, id).Error
}

func (t *File) Create() (err error) {
	return DB.Create(t).Error
}

func (t *File) Get(id uint) (one File, err error) {
	err = DB.First(&one, "id = ?", id).Error
	return
}
