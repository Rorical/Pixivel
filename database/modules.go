package pixivel

type DataTag struct {
	ID   uint64 `gorm:"AUTO_INCREMENT,PRIMARY_KEY"`
	Name string
}
type DataUser struct {
	ID                  uint64 `gorm:"PRIMARY_KEY"`
	Name                string
	Account             string
	ProfileImagesMedium string
	Illusts             []DataIllust `gorm:"foreignKey:User"`
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
	ID                 uint64 `gorm:"PRIMARY_KEY"`
	Title              string
	Type               string
	ImagesSquareMedium string
	ImagesMedium       string
	ImagesLarge        string
	ImagesOriginal     string
	Caption            string
	Restrict           int
	User               uint64
	//Tags     []Tag
	//Tools          []Tools
	CreateData                     string
	PageCount                      int
	Width                          int
	Height                         int
	SanityLevel                    int
	MetaSinglePageOriginalImageURL string
	MetaPages                      []DataMetaPage `gorm:"foreignKey:IllustID"`
	TotalView                      int
	TotalBookmarks                 int
	Visible                        bool
	IsMuted                        bool
	TotalComments                  int
}
