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
	self.db.AutoMigrate(&Illust{}, &MetaPages{}, &User{}, &Tag{})
}

func (self Database) CreateIllust(illust *pixiv.Illust) {
	var isNotExist bool
	var metaLen int
	metaLen = len(illust.MetaPages)
	var newMetaPages []MetaPages = make([]MetaPages, metaLen)
	for j := 0; j < metaLen; j++ {
		newMetaPages[j] = MetaPages{
			IllustID:     illust.ID,
			SquareMedium: illust.MetaPages[j].Images.SquareMedium,
			Medium:       illust.MetaPages[j].Images.Medium,
			Large:        illust.MetaPages[j].Images.Large,
			Original:     illust.MetaPages[j].Images.Original,
		}
	}
	newIllust := &Illust{
		ID:       illust.ID,
		Title:    illust.Title,
		Type:     illust.Type,
		Caption:  illust.Caption,
		Restrict: illust.Restrict,
		User:     illust.User.ID,
		//Tools          []string,
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
		IsBookmarked:                   illust.IsBookmarked,
		Visible:                        illust.Visible,
		IsMuted:                        illust.IsMuted,
		TotalComments:                  illust.TotalComments,
		//MetaPages:      newMetaPages,
		//MetaSinglePage: MetaSinglePage{
		//	OriginalImageURL: illust.MetaSinglePage.OriginalImageURL,
		//},
	}
	isNotExist = self.db.NewRecord(newIllust)
	if isNotExist {
		self.db.Create(newIllust)
	} else {
		self.db.Save(newIllust)
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
