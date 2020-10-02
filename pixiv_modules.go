package pixivel

import "time"

type UserImages struct {
	Medium string `json:"medium"`
}
type User struct {
	ID            uint64     `json:"id"`
	Name          string     `json:"name"`
	Account       string     `json:"account"`
	IsFollowed    bool       `json:"is_followed"`
	ProfileImages UserImages `json:"profile_image_urls"`
}
type UserDetail struct {
	User *User `json:"user"`
	// TODO:
	// Profile
	// ProfilePublicity
	// Workspace
}
type Tag struct {
	Name string `json:"name"`
}
type Images struct {
	SquareMedium string `json:"square_medium"`
	Medium       string `json:"medium"`
	Large        string `json:"large"`
	Original     string `json:"original"`
}
type MetaSinglePage struct {
	OriginalImageURL string `json:"original_image_url"`
}
type MetaPage struct {
	Images Images `json:"image_urls"`
}
type Illust struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	Images      Images `json:"image_urls"`
	Caption     string `json:"caption"`
	Restrict    int    `json:"restrict"`
	User        User   `json:"user"`
	Tags        []Tag  `json:"tags"`
	CreateData  string `json:"create_data"`
	PageCount   int    `json:"page_count"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	SanityLevel int    `json:"sanity_level"`
	// TODO:
	// Series `json:"series"`
	MetaSinglePage MetaSinglePage `json:"meta_single_page"`
	MetaPages      []MetaPage     `json:"meta_pages"`
	TotalView      int            `json:"total_view"`
	TotalBookmarks int            `json:"total_bookmarks"`
	Visible        bool           `json:"visible"`
	IsMuted        bool           `json:"is_muted"`
	TotalComments  int            `json:"total_comments"`
}

type IllustsResponse struct {
	Illusts []Illust `json:"illusts"`
	NextURL string   `json:"next_url"`
}
type IllustResponse struct {
	Illust Illust `json:"illust"`
}

//PixivResponseError PixivResponseError
type PixivResponseError struct {
	Error PixivError `json:"error"`
}
type UserMessageDetail struct {
}
type PixivError struct {
	Message            string            `json:"message"`
	Reason             string            `json:"reason"`
	UserMessage        string            `json:"user_message"`
	UserMessageDetails UserMessageDetail `json:"user_message_details"`
}
type ParentComment struct {
}
type IllustComments struct {
	ID            int           `json:"id"`
	Comment       string        `json:"comment"`
	Date          time.Time     `json:"date"`
	User          User          `json:"user"`
	ParentComment ParentComment `json:"parent_comment"`
	//_       _         `json:"_"`
}
type IllustCommentsResponse struct {
	TotalComments int              `json:"total_comments"`
	Comments      []IllustComments `json:"comments"`
	NextURL       string           `json:"next_url"`
}

type TrendingTagsIllust struct {
	TranslatedName string `json:"translated_name"`
	Tag            string `json:"tag"`
	Illust         Illust `json:"illust"`
}

type TrendingTagsIllustResponse struct {
	TrendTags *[]TrendingTagsIllust `json:"trend_tags"`
}

type UserPreviews struct {
	User    *User    `json:"user"`
	Illusts []Illust `json:"illusts"`
	IsMuted bool     `json:"is_muted"`
	//novels
}
type UserResponse struct {
	UserPreviews []UserPreviews `json:"user_previews"`
	NextURL      string         `json:"next_url"`
}

type UgoiraMetadataFrame struct {
	Delay int    `json:"delay"`
	File  string `json:"file"`
}
type UgoiraMetadataZipUrls struct {
	Medium   string `json:"medium"`
	Large    string `json:"large"`
	Original string `json:"original"`
}

type UgoiraMetadata struct {
	Frames  []UgoiraMetadataFrame `json:"frames"`
	ZipUrls UgoiraMetadataZipUrls `json:"zip_urls"`
}
type UgoiraMetadataResponse struct {
	UgoiraMetadata UgoiraMetadata `json:"ugoira_metadata"`
}
