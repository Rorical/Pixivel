package pixivel

type Tag struct {
	ID   uint64 `gorm:"AUTO_INCREMENT,PRIMARY_KEY"`
	Name string
}
type User struct {
	ID                  uint64 `gorm:"PRIMARY_KEY"`
	Name                string
	Account             string
	Comment             string
	IsFollowed          bool
	ProfileImagesMedium string
}
type MetaPage struct {
	ID           uint64 `gorm:"AUTO_INCREMENT,PRIMARY_KEY"`
	IllustID     uint64
	SquareMedium string
	Medium       string
	Large        string
	Original     string
}
type Illust struct {
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
	MetaPages                      []MetaPage `gorm:"foreignKey:IllustID"`
	TotalView                      int
	TotalBookmarks                 int
	IsBookmarked                   bool
	Visible                        bool
	IsMuted                        bool
	TotalComments                  int
}
