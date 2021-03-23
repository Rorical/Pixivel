package pixiv

import "time"

type UserImages struct {
	Medium string `json:"medium"`
}
type User struct {
	ID            uint64     `json:"id"`
	Name          string     `json:"name"`
	Account       string     `json:"account"`
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
	Name string `json:"name,omitempty"`
}
type Images struct {
	SquareMedium string `json:"square_medium,omitempty"`
	Medium       string `json:"medium,omitempty"`
	Large        string `json:"large,omitempty"`
	Original     string `json:"original,omitempty"`
}
type MetaSinglePage struct {
	OriginalImageURL string `json:"original_image_url,omitempty"`
}
type MetaPage struct {
	Images Images `json:"image_urls,omitempty"`
}
type Illust struct {
	ID          uint64    `json:"id"`
	Title       string    `json:"title,omitempty"`
	Type        string    `json:"type,omitempty"`
	Images      Images    `json:"image_urls,omitempty"`
	Caption     string    `json:"caption,omitempty"`
	User        User      `json:"user,omitempty"`
	Tags        []Tag     `json:"tags,omitempty"`
	PageCount   uint      `json:"page_count,omitempty"`
	Width       uint      `json:"width,omitempty"`
	Height      uint      `json:"height,omitempty"`
	SanityLevel uint      `json:"sanity_level,omitempty"`
	CreateDate  time.Time `json:"create_date,omitempty"`
	// TODO:
	// Series `json:"series"`
	MetaSinglePage MetaSinglePage `json:"meta_single_page,omitempty"`
	MetaPages      []MetaPage     `json:"meta_pages,omitempty"`
	TotalView      uint           `json:"total_view,omitempty"`
	TotalBookmarks uint           `json:"total_bookmarks,omitempty"`
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
	ID            uint          `json:"id"`
	Comment       string        `json:"comment"`
	Date          time.Time     `json:"date"`
	User          User          `json:"user"`
	ParentComment ParentComment `json:"parent_comment"`
	//_       _         `json:"_"`
}
type IllustCommentsResponse struct {
	TotalComments uint             `json:"total_comments"`
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
	Delay uint   `json:"delay"`
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
