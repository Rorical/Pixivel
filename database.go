package pixivel

import (
	"github.com/Rorical/pixiv"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Database struct {
	db *gorm.DB
}

func GetDB() *Database {
	db, err := gorm.Open("sqlite3", "cache.db")
	if err != nil {
		panic("failed to connect database")
	}
	//defer db.Close()
	return &Database{
		db: db,
	}
}

func (self Database) Migrate() {
	self.db.AutoMigrate(&DataIllust{}, &DataMetaPage{}, &DataUser{})
}

func (self Database) CreateIllust(illust *pixiv.Illust) {
	var isExist bool
	var metaLen int
	metaLen = len(illust.MetaPages)
	var newMetaPages []DataMetaPage = make([]DataMetaPage, metaLen)
	for j := 0; j < metaLen; j++ {
		newMetaPages[j] = DataMetaPage{
			IllustID:     illust.ID,
			SquareMedium: illust.MetaPages[j].Images.SquareMedium,
			Medium:       illust.MetaPages[j].Images.Medium,
			Large:        illust.MetaPages[j].Images.Large,
			Original:     illust.MetaPages[j].Images.Original,
		}
	}
	newIllust := &DataIllust{
		ID:                             illust.ID,
		Title:                          illust.Title,
		Type:                           illust.Type,
		Caption:                        illust.Caption,
		Restrict:                       illust.Restrict,
		MetaPages:                      newMetaPages,
		CreateData:                     illust.CreateData,
		PageCount:                      illust.PageCount,
		Width:                          illust.Width,
		Height:                         illust.Height,
		SanityLevel:                    illust.SanityLevel,
		ImagesSquareMedium:             illust.Images.SquareMedium,
		ImagesMedium:                   illust.Images.Medium,
		ImagesLarge:                    illust.Images.Large,
		ImagesOriginal:                 illust.Images.Original,
		MetaSinglePageOriginalImageURL: illust.MetaSinglePage.OriginalImageURL,
		TotalView:                      illust.TotalView,
		TotalBookmarks:                 illust.TotalBookmarks,
		Visible:                        illust.Visible,
		IsMuted:                        illust.IsMuted,
		TotalComments:                  illust.TotalComments,
	}
	isExist = self.db.NewRecord(newIllust)
	if isExist {
		self.db.Save(&newIllust)
		return
	}
	self.db.Create(&newIllust)

	newUser := &DataUser{
		ID:                  illust.User.ID,
		Name:                illust.User.Name,
		Account:             illust.User.Account,
		ProfileImagesMedium: illust.User.ProfileImages.Medium,
	}
	isExist = self.db.NewRecord(newUser)
	if isExist {
		self.db.Save(&newUser)
		return
	}
}

func (self Database) QueryIllust(id uint64) *Illust {
	var illust Illust
	self.db.Preload("MetaPages").First(&illust, id)
	return &illust
}

func (self Database) Close() {
	self.db.Close()
}
