package pixivel

import (
	"fmt"

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
	self.db.AutoMigrate(&Illust{}, &MetaPages{}, &User{}, &Illust2Tag{}, &Tag{})
}

func (self Database) CreateIllust(illust *pixiv.Illust) {
	var isNotExist bool
	metaLen := len(illust.Tags)
	var newTag *Tag
	var newIllust2Tag *Illust2Tag
	newTag = new(Tag)
	newIllust2Tag = new(Illust2Tag)
	for j := 0; j < metaLen; j++ {

		newTag = &Tag{
			Name: illust.Tags[j].Name,
		}
		isNotExist = self.db.NewRecord(newTag)

		fmt.Println(isNotExist)

		if isNotExist {
			self.db.Create(newTag)
		} else {
			self.db.Save(newTag)
		}

		newIllust2Tag = &Illust2Tag{
			TagID:    newTag.ID,
			IllustID: illust.ID,
		}
		isNotExist = self.db.NewRecord(newIllust2Tag)
		fmt.Println(isNotExist)
		if isNotExist {
			self.db.Create(newIllust2Tag)
		} else {
			self.db.Save(newIllust2Tag)
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
	metaLen = len(illust.MetaPages)
	for j := 0; j < metaLen; j++ {
		newMetaPage := &MetaPages{
			IllustID:     illust.ID,
			SquareMedium: illust.MetaPages[j].Images.SquareMedium,
			Medium:       illust.MetaPages[j].Images.Medium,
			Large:        illust.MetaPages[j].Images.Large,
			Original:     illust.MetaPages[j].Images.Original,
		}
		isNotExist = self.db.NewRecord(newMetaPage)
		if isNotExist {
			self.db.Create(newMetaPage)
		} else {
			self.db.Save(newMetaPage)
		}
	}
	newUser := &User{
		ID:                  illust.User.ID,
		Name:                illust.User.Name,
		Account:             illust.User.Account,
		Comment:             illust.User.Comment,
		IsFollowed:          illust.User.IsFollowed,
		ProfileImagesMedium: illust.User.ProfileImages.Medium,
	}
	isNotExist = self.db.NewRecord(newUser)
	if isNotExist {
		self.db.Create(newUser)
	} else {
		self.db.Save(newUser)
	}

}

func (self Database) QueryIllust(id uint64) *Illust {
	var illust Illust
	var metaPages []MetaPages
	var illust2Tag []Illust2Tag
	self.db.Where(&Illust{
		ID: id,
	}).First(&illust)
	self.db.Where(&MetaPages{
		IllustID: id,
	}).Find(&metaPages)
	self.db.Where(&Illust2Tag{
		IllustID: id,
	}).Find(&illust2Tag)

	Len := len(illust2Tag)
	var tag *Tag
	tag = new(Tag)
	for j := 0; j < Len; j++ {
		self.db.Where(&Tag{
			ID: illust2Tag[j].TagID,
		}).First(tag)
		fmt.Println(tag.Name)
	}

	return &illust
}

func (self Database) Close() {
	self.db.Close()
}
