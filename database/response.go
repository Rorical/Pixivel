package pixivel

type Tag interface {
	Name string
}
type Images interface {
	SquareMedium string
	Medium       string
	Large        string
	Original     string
}
type MetaSinglePage interface {
	OriginalImageURL string
}
type MetaPage interface {
	Images Images
}
type Illust interface {
	ID          uint64
	Title       string
	Type        string
	Images      Images
	Caption     string
	Restrict    int
	User        User
	Tags        []Tag
	Tools       []string
	CreateData  string
	PageCount   int
	Width       int
	Height      int
	SanityLevel int
	MetaSinglePage MetaSinglePage
	MetaPages      []MetaPage
	TotalView      int
	TotalBookmarks int
	IsBookmarked   bool
	Visible        bool
	IsMuted        bool
	TotalComments  int
}