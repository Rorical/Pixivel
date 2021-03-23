package database

import "time"

type HashDB interface {
	Get(key []byte) ([]byte, error)
	Set(key []byte, value []byte) error
	Has(key []byte) (bool, error)
	Del(key []byte) error
	Close()
	IsErrNotFound(err error) bool
}

type DataTag struct {
	ID      uint64       `gorm:"AUTO_INCREMENT,PRIMARY_KEY"`
	Name    string       `gorm:"uniqueIndex"`
	Illusts []DataIllust `gorm:"many2many:illusts_tags;"`
}
type DataUser struct {
	ID                  uint64 `gorm:"PRIMARY_KEY"`
	Name                string
	Account             string
	ProfileImagesMedium string
	Illusts             []DataIllust `gorm:"foreignKey:User;references:ID"`
}
type DataMetaPage struct {
	ID           uint64 `gorm:"AUTO_INCREMENT,PRIMARY_KEY"`
	IllustID     uint64
	SquareMedium string
	Medium       string
	Large        string
	Original     string
}

type DataIllustCollection struct {
	ID      uint64       `gorm:"AUTO_INCREMENT,PRIMARY_KEY"`
	Illusts []DataIllust `gorm:"many2many:illust_collection;"`
}

type DataRankIllusts struct {
	ID      uint64 `gorm:"AUTO_INCREMENT,PRIMARY_KEY"`
	Type    uint
	Date    time.Time
	Illusts DataIllustCollection `gorm:"foreignKey:ID"`
}

type DataIllust struct {
	ID                             uint64 `gorm:"PRIMARY_KEY,many2many:illusts_illust"`
	Title                          string
	Type                           uint
	ImagesSquareMedium             string
	ImagesMedium                   string
	ImagesLarge                    string
	Caption                        string
	User                           uint64
	CreateDate                     time.Time
	Tags                           []DataTag `gorm:"many2many:illusts_tags;"`
	PageCount                      uint
	Width                          uint
	Height                         uint
	SanityLevel                    uint
	MetaSinglePageOriginalImageURL string
	MetaPages                      []DataMetaPage `gorm:"foreignKey:IllustID"`
	TotalView                      uint
	TotalBookmarks                 uint
}
type DataUgoiraMetadata struct {
	ID             uint64                    `gorm:"PRIMARY_KEY"`
	Frames         []DataUgoiraMetadataFrame `gorm:"foreignKey:UgoiraID"`
	ZipURLMedium   string
	ZipURLLarge    string
	ZipURLOriginal string
}

type DataUgoiraMetadataFrame struct {
	UgoiraID uint64
	Delay    int
	File     string
}
