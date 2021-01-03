package database

import (
	"Pixivel/internal/config"
	"Pixivel/internal/levelgo"
	"Pixivel/internal/pixiv"
	"errors"
	"strconv"

	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	db     *gorm.DB
	HashDB HashDB
}

var RECORD_NOT_FOUND = errors.New("No Result")
var illustTypes2Num map[string]uint = map[string]uint{"illust": 0, "manga": 1, "ugoira": 2}
var illustNum2Types map[uint]string = map[uint]string{0: "illust", 1: "manga", 2: "ugoira"}

func GetDB(settings *config.Setting) *Database {
	db, err := gorm.Open(settings.SQL.Type, settings.SQL.URI)
	if err != nil {
		panic(err)
	}
	//redisPool := NewRedisPool()
	leveldb := levelgo.RpcClient(settings.HashDB.URI)
	leveldb.Connect()
	return &Database{
		db: db,
		//RedisPool: redisPool,
		//Redis:     redisPool.NewRedisClient(),
		HashDB: leveldb,
	}

}

func (self *Database) Migrate() {
	self.db.AutoMigrate(&DataIllust{}, &DataMetaPage{}, &DataUser{}, &DataTag{}, &DataUgoiraMetadata{}, &DataUgoiraMetadataFrame{})
}

func (self *Database) IsTheSame(face interface{}, hashKey string) bool {
	hash := config.HashStruct(face)
	bytehash := config.StringIn(hashKey)
	res, err := self.HashDB.Get(bytehash)
	if self.HashDB.IsErrNotFound(err) {
		self.HashDB.Set(bytehash, config.StringIn(hash))
		return false
	}
	strres := config.StringOut(res)
	if err != nil {
		panic(err)
	}
	if hash == strres {
		return true
	}
	return false
}

func (self *Database) CreateIllust(illust *pixiv.Illust) {
	var err error
	same := self.IsTheSame(illust, "i"+strconv.FormatUint(illust.ID, 10))
	if same {
		return
	}
	newIllust := DataIllust{
		ID:                             illust.ID,
		Title:                          illust.Title,
		Type:                           illustTypes2Num[illust.Type],
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
	illustModel := self.db.Model(&newIllust)

	metaLen := len(illust.MetaPages)
	var singleMetaPage *DataMetaPage

	self.db.Where(&DataMetaPage{IllustID: illust.ID}).Delete(DataMetaPage{})
	for j := 0; j < metaLen; j++ {
		singleMetaPage = &DataMetaPage{
			IllustID:     illust.ID,
			SquareMedium: illust.MetaPages[j].Images.SquareMedium,
			Medium:       illust.MetaPages[j].Images.Medium,
			Large:        illust.MetaPages[j].Images.Large,
			Original:     illust.MetaPages[j].Images.Original,
		}
		//self.db.Save(singleMetaPage)
		illustModel.Association("MetaPages").Append(singleMetaPage)
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
		}
		self.db.Model(&newIllust).Association("Tags").Append(newDataTag)
	}
	newUser := &DataUser{
		ID:                  illust.User.ID,
		Name:                illust.User.Name,
		Account:             illust.User.Account,
		ProfileImagesMedium: illust.User.ProfileImages.Medium,
	}
	same = self.IsTheSame(illust.User, "u"+strconv.FormatUint(illust.User.ID, 10))
	if same {
		self.db.Model(&newUser).Association("Illusts").Append(&newIllust)
		return
	}

	self.db.Where(&DataUser{ID: illust.User.ID}).First(&DataUser{})
	self.db.Save(&newUser)

	self.db.Model(&newUser).Association("Illusts").Append(&newIllust)

}

func (self *Database) QueryIllust(id uint64) (*pixiv.Illust, error) {

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
	newTags := make([]pixiv.Tag, lena)
	for j := 0; j < lena; j++ {
		newTags[j] = pixiv.Tag{
			Name: tags[j].Name,
		}
	}

	lena = len(metapages)
	newMetaPages := make([]pixiv.MetaPage, lena)
	for j := 0; j < lena; j++ {
		newMetaPages[j] = pixiv.MetaPage{
			Images: pixiv.Images{
				SquareMedium: metapages[j].SquareMedium,
				Medium:       metapages[j].Medium,
				Large:        metapages[j].Large,
				Original:     metapages[j].Original,
			},
		}
	}

	ResponseIllust := pixiv.Illust{
		ID:          illust.ID,
		Title:       illust.Title,
		Type:        illustNum2Types[illust.Type],
		Caption:     illust.Caption,
		Restrict:    illust.Restrict,
		PageCount:   illust.PageCount,
		Width:       illust.Width,
		Height:      illust.Height,
		SanityLevel: illust.SanityLevel,
		Tags:        newTags,
		Images: pixiv.Images{
			SquareMedium: illust.ImagesSquareMedium,
			Medium:       illust.ImagesMedium,
			Large:        illust.ImagesLarge,
		},
		MetaSinglePage: pixiv.MetaSinglePage{
			OriginalImageURL: illust.MetaSinglePageOriginalImageURL,
		},
		User: pixiv.User{
			ID:      user.ID,
			Name:    user.Name,
			Account: user.Account,
			ProfileImages: pixiv.UserImages{
				Medium: user.ProfileImagesMedium,
			},
		},
		MetaPages:      newMetaPages,
		TotalView:      illust.TotalView,
		TotalBookmarks: illust.TotalBookmarks,
	}

	return &ResponseIllust, nil
}

func (self *Database) DeleteIllust(id uint64) error {
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

func (self *Database) Close() {
	self.db.Close()
	self.HashDB.Close()
}
