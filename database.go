package pixivel

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Database struct {
	db *gorm.DB
}

var RECORD_NOT_FOUND = errors.New("No Result")

func GetDB(uri string) *Database {
	db, err := gorm.Open("sqlite3", uri)
	if err != nil {
		panic("failed to connect database")
	}
	//defer db.Close()
	return &Database{
		db: db,
	}
}

func (self Database) Migrate() {
	self.db.AutoMigrate(&DataIllust{}, &DataMetaPage{}, &DataUser{}, &DataTag{}, &DataUgoiraMetadata{}, &DataUgoiraMetadataFrame{})
}

func (self Database) CreateIllust(illust *Illust) {
	var err error
	newIllust := DataIllust{
		ID:                             illust.ID,
		Title:                          illust.Title,
		Type:                           illust.Type,
		Caption:                        illust.Caption,
		Restrict:                       illust.Restrict,
		PageCount:                      illust.PageCount,
		Width:                          illust.Width,
		Height:                         illust.Height,
		SanityLevel:                    illust.SanityLevel,
		ImagesSquareMedium:             illust.Images.SquareMedium,
		ImagesMedium:                   illust.Images.Medium,
		ImagesLarge:                    illust.Images.Large,
		MetaSinglePageOriginalImageURL: illust.MetaSinglePage.OriginalImageURL,
		TotalView:                      illust.TotalView,
		TotalBookmarks:                 illust.TotalBookmarks,
	}
	self.db.Save(&newIllust)

	var metaLen int

	metaLen = len(illust.MetaPages)
	var newDataMetaPage *DataMetaPage
	self.db.Model(&newIllust).Association("MetaPages").Clear()
	for j := 0; j < metaLen; j++ {
		newDataMetaPage = &DataMetaPage{}
		err = self.db.Where(&DataMetaPage{Original: illust.MetaPages[j].Images.Original}).First(newDataMetaPage).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newDataMetaPage = &DataMetaPage{
				IllustID:     illust.ID,
				SquareMedium: illust.MetaPages[j].Images.SquareMedium,
				Medium:       illust.MetaPages[j].Images.Medium,
				Large:        illust.MetaPages[j].Images.Large,
				Original:     illust.MetaPages[j].Images.Original,
			}
			self.db.Create(newDataMetaPage)
		} else {
			newDataMetaPage.SquareMedium = illust.MetaPages[j].Images.SquareMedium
			newDataMetaPage.Medium = illust.MetaPages[j].Images.Medium
			newDataMetaPage.Large = illust.MetaPages[j].Images.Large
			newDataMetaPage.Original = illust.MetaPages[j].Images.Original
			self.db.Save(newDataMetaPage)
		}
		self.db.Model(&newIllust).Association("MetaPages").Append(newDataMetaPage)

	}

	metaLen = len(illust.Tags)
	var tagName string
	var newDataTag *DataTag
	self.db.Model(&newIllust).Association("Tags").Clear()
	for j := 0; j < metaLen; j++ {
		tagName = illust.Tags[j].Name
		newDataTag = &DataTag{}
		err = self.db.Where(&DataTag{Name: tagName}).First(newDataTag).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newDataTag = &DataTag{
				Name: tagName,
			}
			self.db.Create(newDataTag)
		} else {
			newDataTag.Name = tagName
			self.db.Save(newDataTag)
		}
		self.db.Model(&newIllust).Association("Tags").Append(newDataTag)
	}

	newUser := &DataUser{
		ID:                  illust.User.ID,
		Name:                illust.User.Name,
		Account:             illust.User.Account,
		ProfileImagesMedium: illust.User.ProfileImages.Medium,
	}
	err = self.db.Where(&DataUser{ID: illust.User.ID}).First(&DataUser{}).Error
	self.db.Save(&newUser)

	self.db.Model(&newUser).Association("Illusts").Append(&newIllust)

}

func (self Database) QueryIllust(id uint64) (*Illust, error) {

	var illust DataIllust
	var user DataUser
	var tags []DataTag
	var metapages []DataMetaPage
	illudb := self.db.First(&illust, id)
	if errors.Is(illudb.Error, gorm.ErrRecordNotFound) {
		return nil, RECORD_NOT_FOUND
	}
	illudb.Association("Tags").Find(&tags)
	illudb.Association("MetaPages").Find(&metapages)
	self.db.First(&user, illust.User)

	lena := len(tags)
	newTags := make([]Tag, lena)
	for j := 0; j < lena; j++ {
		newTags[j] = Tag{
			Name: tags[j].Name,
		}
	}

	lena = len(metapages)
	newMetaPages := make([]MetaPage, lena)
	for j := 0; j < lena; j++ {
		newMetaPages[j] = MetaPage{
			Images: Images{
				SquareMedium: metapages[j].SquareMedium,
				Medium:       metapages[j].Medium,
				Large:        metapages[j].Large,
				Original:     metapages[j].Original,
			},
		}
	}

	ResponseIllust := Illust{
		ID:          illust.ID,
		Title:       illust.Title,
		Type:        illust.Type,
		Caption:     illust.Caption,
		Restrict:    illust.Restrict,
		PageCount:   illust.PageCount,
		Width:       illust.Width,
		Height:      illust.Height,
		SanityLevel: illust.SanityLevel,
		Tags:        newTags,
		Images: Images{
			SquareMedium: illust.ImagesSquareMedium,
			Medium:       illust.ImagesMedium,
			Large:        illust.ImagesLarge,
		},
		MetaSinglePage: MetaSinglePage{
			OriginalImageURL: illust.MetaSinglePageOriginalImageURL,
		},
		User: User{
			ID:      user.ID,
			Name:    user.Name,
			Account: user.Account,
			ProfileImages: UserImages{
				Medium: user.ProfileImagesMedium,
			},
		},
		MetaPages:      newMetaPages,
		TotalView:      illust.TotalView,
		TotalBookmarks: illust.TotalBookmarks,
	}

	return &ResponseIllust, nil
}

func (self Database) DeleteIllust(id uint64) error {
	var illust DataIllust
	var user DataUser
	var tags []DataTag
	var metapages []DataMetaPage
	illudb := self.db.First(&illust, id)
	if errors.Is(illudb.Error, gorm.ErrRecordNotFound) {
		return RECORD_NOT_FOUND
	}

	illudb.Association("Tags").Find(&tags)
	illudb.Association("MetaPages").Find(&metapages)

	self.db.First(&user, illust.User)

	self.db.Model(&user).Association("Illusts").Delete(&illust)

	self.db.Model(&illust).Association("Tags").Delete(&tags)

	self.db.Model(&illust).Association("MetaPages").Delete(&metapages)
	self.db.Delete(&metapages)
	self.db.Delete(&illust)
	return nil
}

func (self Database) Close() {
	self.db.Close()
}
