package pixivel

type DataTag struct {
	ID      uint64 `gorm:"AUTO_INCREMENT,PRIMARY_KEY"`
	Name    string
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
type DataIllust struct {
	ID                             uint64 `gorm:"PRIMARY_KEY"`
	Title                          string
	Type                           string
	ImagesSquareMedium             string
	ImagesMedium                   string
	ImagesLarge                    string
	Caption                        string
	Restrict                       int
	User                           uint64
	Tags                           []DataTag `gorm:"many2many:illusts_tags;"`
	PageCount                      int
	Width                          int
	Height                         int
	SanityLevel                    int
	MetaSinglePageOriginalImageURL string
	MetaPages                      []DataMetaPage `gorm:"foreignKey:IllustID"`
	TotalView                      int
	TotalBookmarks                 int
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
