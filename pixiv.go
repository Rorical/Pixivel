package pixivel

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dghubble/sling"
)

const (
	apiBase = "https://app-api.pixiv.net/"
)

//AppPixiv AppPixivAPI
type AppPixiv struct {
	sling   *sling.Sling
	BaseAPI *BasePixiv
}

//AppPixivAPI AppPixiv
func AppPixivAPI() *AppPixiv {
	s := sling.New().Base(apiBase).Set("User-Agent", "PixivIOSApp/7.6.2 (iOS 12.2; iPhone9,1)").Set("App-Version", "7.6.2").Set("App-OS-VERSION", "12.2").Set("App-OS", "ios")
	baseAPI := BasePixivAPI()
	return &AppPixiv{
		sling:   s,
		BaseAPI: baseAPI,
	}
}

/*
type _Params struct {
	_ _ `url:"_,omitemtpy"`
}

func (api *AppPixiv) _(_) (*_, error) {
	path := ""
	data := &_{}
	params := &_{
		_,
	}
	if err := api.request(path, params, data, true); err != nil {
		return nil, err
	}
	return &_, nil
}

*/

func (api *AppPixiv) request(path string, params, data interface{}, auth bool) (err error) {
requestStart:
	erro := &PixivResponseError{}
	if auth {
		if _, err := api.BaseAPI.Login("", ""); err != nil {
			return fmt.Errorf("refresh token failed: %v", err)
		}
		_, err = api.sling.New().Get(path).Set("Authorization", "Bearer "+api.BaseAPI.AccessToken).QueryStruct(params).Receive(data, erro)
		if strings.Contains(erro.Error.Message, "invalid_grant") {
			api.BaseAPI.TokenDeadline = time.Now()
			if _, err := api.BaseAPI.Login("", ""); err != nil {
				return fmt.Errorf("refresh token failed: %v", err)
			}
			goto requestStart
		}
	} else {
		_, err = api.sling.New().Get(path).QueryStruct(params).Receive(data, erro)
	}

	switch {
	case erro.Error.UserMessage != "":
		return errors.New(erro.Error.UserMessage)
	case erro.Error.Message != "":
		return errors.New(erro.Error.UserMessage)
	}
	return err
}

type illustDetailParams struct {
	IllustID uint64 `url:"illust_id,omitemtpy"`
}

//IllustDetail gets the detail of an illust
func (api *AppPixiv) IllustDetail(id uint64) (*Illust, error) {
	path := "v1/illust/detail"
	data := &IllustResponse{}
	params := &illustDetailParams{
		IllustID: id,
	}
	if err := api.request(path, params, data, true); err != nil {
		return nil, err
	}
	return &data.Illust, nil
}

type userDetailParams struct {
	UserID uint64 `url:"user_id,omitempty"`
	Filter string `url:"filter,omitempty"`
}

//UserDetail gets users information
func (api *AppPixiv) UserDetail(uid uint64) (*User, error) {
	path := "v1/user/detail"
	params := &userDetailParams{
		UserID: uid,
		Filter: "for_ios",
	}
	detail := &UserDetail{
		User: &User{},
	}
	if err := api.request(path, params, detail, true); err != nil {
		return nil, err
	}

	return detail.User, nil
}

type userIllustsParams struct {
	UserID uint64 `url:"user_id,omitempty"`
	Filter string `url:"filter,omitempty"`
	Type   string `url:"type,omitempty"`
	Offset int    `url:"offset,omitempty"`
}

// UserIllusts type: [illust, manga]
func (api *AppPixiv) UserIllusts(uid uint64, _type string, offset int) ([]Illust, int, error) {
	path := "v1/user/illusts"
	params := &userIllustsParams{
		UserID: uid,
		Filter: "for_ios",
		Type:   _type,
		Offset: offset,
	}
	data := &IllustsResponse{}
	if err := api.request(path, params, data, true); err != nil {
		return nil, 0, err
	}
	next, err := parseNextPageOffset(data.NextURL)
	return data.Illusts, next, err
}

type userBookmarkIllustsParams struct {
	UserID        uint64 `url:"user_id,omitempty"`
	Restrict      string `url:"restrict,omitempty"`
	Filter        string `url:"filter,omitempty"`
	MaxBookmarkID int    `url:"max_bookmark_id,omitempty"`
	Tag           string `url:"tag,omitempty"`
}

// UserBookmarksIllust restrict: [public, private]
func (api *AppPixiv) UserBookmarksIllust(uid uint64, restrict string, maxBookmarkID int, tag string) ([]Illust, int, error) {
	path := "v1/user/bookmarks/illust"
	params := &userBookmarkIllustsParams{
		UserID:        uid,
		Restrict:      "public",
		Filter:        "for_ios",
		MaxBookmarkID: maxBookmarkID,
		Tag:           tag,
	}
	data := &IllustsResponse{}
	if err := api.request(path, params, data, true); err != nil {
		return nil, 0, err
	}
	next, err := parseNextPageOffset(data.NextURL)
	return data.Illusts, next, err
}

type illustFollowParams struct {
	Restrict string `url:"restrict,omitempty"`
	Offset   int    `url:"offset,omitempty"`
}

// IllustFollow restrict: [public, private]
func (api *AppPixiv) IllustFollow(restrict string, offset int) ([]Illust, int, error) {
	path := "v2/illust/follow"
	params := &illustFollowParams{
		Restrict: restrict,
		Offset:   offset,
	}
	data := &IllustsResponse{}
	if err := api.request(path, params, data, true); err != nil {
		return nil, 0, err
	}
	next, err := parseNextPageOffset(data.NextURL)
	return data.Illusts, next, err
}

//IllustCommentsParams is used by function IllustComments
type IllustCommentsParams struct {
	IllustID             uint64 `url:"illust_id,omitemtpy"`
	Offset               int    `url:"offset,omitempty"`
	IncludeTotalComments bool   `url:"include_total_comments,omitempty"`
}

//IllustComments get the comments of an Illust
func (api *AppPixiv) IllustComments(IllustID uint64, offset int, IncludeTotalComments bool) (*IllustCommentsResponse, int, error) {
	path := "v1/illust/comments"
	data := &IllustCommentsResponse{}
	params := &IllustCommentsParams{
		IllustID:             IllustID,
		Offset:               offset,
		IncludeTotalComments: IncludeTotalComments,
	}
	if err := api.request(path, params, data, true); err != nil {
		return nil, 0, err
	}
	next, err := parseNextPageOffset(data.NextURL)
	return data, next, err
}

type IllustRelatedParams struct {
	IllustID uint64 `url:"illust_id,omitemtpy"`
	Offset   int    `url:"offset,omitempty"`
	Filter   string `url:"filter,omitempty"`
}

//IllustRelated get the related illusts of an Illust
func (api *AppPixiv) IllustRelated(IllustID uint64, offset int) (*IllustsResponse, int, error) {
	path := "v2/illust/related"
	data := &IllustsResponse{}
	params := &IllustRelatedParams{
		IllustID: IllustID,
		Offset:   offset,
		Filter:   "for_ios",
	}
	if err := api.request(path, params, data, true); err != nil {
		return nil, 0, err
	}
	next, err := parseNextPageOffset(data.NextURL)
	return data, next, err
}

// IllustRecommendedParams IllustRecommended
type IllustRecommendedParams struct {
	ContentType                  string `url:"content_type,omitemtpy"`
	includeRankingLabel          bool   `url:"include_ranking_label,omitempty"`
	Filter                       string `url:"filter,omitempty"`
	MaxBookmarkIDForRecommend    uint64 `url:"max_bookmark_id_for_recommend,omitempty"`
	MinBookmarkIDForRecentIllust uint64 `url:"min_bookmark_id_for_recent_illust,omitempty"`
	Offset                       int    `url:"offset,omitempty"`
	IncludeRankingIllusts        bool   `url:"include_ranking_illusts,omitempty"`
}

//IllustRecommended ("illust", true, 0, 0, 0, true, true)
func (api *AppPixiv) IllustRecommended(contentType string, includeRankingLabel bool, maxBookmarkIDForRecommend uint64, minBookmarkIDForRecentIllust uint64, offset int, includeRankingIllusts bool, reqAuth bool) (*IllustsResponse, int, error) {
	var path string
	if reqAuth {
		path = "v1/illust/recommended"
	} else {
		path = "v1/illust/recommended-nologin"
	}

	data := &IllustsResponse{}
	params := &IllustRecommendedParams{
		ContentType:                  contentType,
		includeRankingLabel:          includeRankingLabel,
		Filter:                       "for_ios",
		MaxBookmarkIDForRecommend:    maxBookmarkIDForRecommend,
		MinBookmarkIDForRecentIllust: minBookmarkIDForRecentIllust,
		Offset:                       offset,
		IncludeRankingIllusts:        includeRankingIllusts,
	}
	if err := api.request(path, params, data, reqAuth); err != nil {
		return nil, 0, err
	}
	next, err := parseNextPageOffset(data.NextURL)
	return data, next, err
}

// IllustRankingParams IllustRanking
type IllustRankingParams struct {
	Mode   string `url:"mode,omitemtpy"`
	Date   string `url:"date,omitempty"`
	Filter string `url:"filter,omitempty"`
	Offset int    `url:"offset,omitempty"`
}

//IllustRanking ("day", "", 0)
func (api *AppPixiv) IllustRanking(Mode string, Date string, offset int) (*IllustsResponse, int, error) {
	path := "v1/illust/ranking"
	data := &IllustsResponse{}
	params := &IllustRankingParams{
		Filter: "for_ios",
		Mode:   Mode,
		Date:   Date,
		Offset: offset,
	}
	if err := api.request(path, params, data, true); err != nil {
		return nil, 0, err
	}
	next, err := parseNextPageOffset(data.NextURL)
	return data, next, err
}

//TrendingTagsIllustParams TrendingTagsIllust
type TrendingTagsIllustParams struct {
	Filter string `url:"filter,omitempty"`
}

//TrendingTagsIllust ("day", "", 0)
func (api *AppPixiv) TrendingTagsIllust() (*TrendingTagsIllust, error) {
	path := "v1/trending-tags/illust"
	data := &TrendingTagsIllustResponse{
		TrendTags: *[]TrendingTagsIllust{},
	}
	params := &TrendingTagsIllustParams{
		Filter: "for_ios",
	}
	if err := api.request(path, params, data, true); err != nil {
		return nil, err
	}
	return data.TrendTags, nil
}

type SearchIllustParams struct {
	Word         string `url:"word,omitempty"`
	SearchTarget string `url:"search_target,omitempty"`
	Sort         string `url:"sort,omitempty"`
	StartDate    string `url:"start_date,omitempty"`
	EndDate      string `url:"end_date,omitempty"`
	Duration     string `url:"duration,omitempty"`
	Offset       int    `url:"offset,omitempty"`
	Filter       string `url:"filter,omitempty"`
}

//SearchIllust ("Sagiri", "partial_match_for_tags", "date_desc", "", "", "", 0)
func (api *AppPixiv) SearchIllust(Word string, searchTarget string, sort string, startDate string, endDate string, duration string, offset int) (*IllustsResponse, int, error) {
	path := "v1/search/illust"
	data := &IllustsResponse{}
	params := &SearchIllustParams{
		Filter:       "for_ios",
		Word:         Word,
		SearchTarget: searchTarget,
		Sort:         sort,
		StartDate:    startDate,
		EndDate:      endDate,
		Duration:     duration,
		Offset:       offset,
	}
	if err := api.request(path, params, data, true); err != nil {
		return nil, 0, err
	}
	next, err := parseNextPageOffset(data.NextURL)
	return data, next, err
}

type SearchUserParams struct {
	Word     string `url:"word,omitempty"`
	Sort     string `url:"sort,omitempty"`
	Duration string `url:"duration,omitempty"`
	Offset   int    `url:"offset,omitempty"`
	Filter   string `url:"filter,omitempty"`
}

//SearchUser ("Quan_", "date_desc", "", 0)
func (api *AppPixiv) SearchUser(Word string, sort string, duration string, offset int) (*UserResponse, int, error) {
	path := "v1/search/user"
	data := &UserResponse{}
	params := &SearchUserParams{
		Filter:   "for_ios",
		Word:     Word,
		Sort:     sort,
		Duration: duration,
		Offset:   offset,
	}
	if err := api.request(path, params, data, true); err != nil {
		return nil, 0, err
	}
	next, err := parseNextPageOffset(data.NextURL)
	return data, next, err
}

type UgoiraMetadataParams struct {
	IllustID uint64 `url:"illust_id,omitempty"`
}

//UgoiraMetadata ("day", "", 0)
func (api *AppPixiv) UgoiraMetadata(illustID uint64) (*UgoiraMetadataResponse, error) {
	path := "v1/ugoira/metadata"
	data := &UgoiraMetadataResponse{}
	params := &UgoiraMetadataParams{
		IllustID: illustID,
	}
	if err := api.request(path, params, data, true); err != nil {
		return nil, err
	}
	return data, nil
}
