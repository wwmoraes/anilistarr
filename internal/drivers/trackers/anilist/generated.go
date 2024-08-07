// Code generated by github.com/Khan/genqlient, DO NOT EDIT.

package anilist

import (
	"context"

	"github.com/Khan/genqlient/graphql"
)

// GetUserByNameResponse is returned by GetUserByName on success.
type GetUserByNameResponse struct {
	// User query
	User GetUserByNameUser `json:"User"`
}

// GetUser returns GetUserByNameResponse.User, and is useful for accessing the field via an interface.
func (v *GetUserByNameResponse) GetUser() GetUserByNameUser { return v.User }

// GetUserByNameUser includes the requested fields of the GraphQL type User.
// The GraphQL type's documentation follows.
//
// A user
type GetUserByNameUser struct {
	// The id of the user
	Id int `json:"id"`
}

// GetId returns GetUserByNameUser.Id, and is useful for accessing the field via an interface.
func (v *GetUserByNameUser) GetId() int { return v.Id }

// GetWatchingPage includes the requested fields of the GraphQL type Page.
// The GraphQL type's documentation follows.
//
// Page of data
type GetWatchingPage struct {
	MediaList []GetWatchingPageMediaList `json:"mediaList"`
}

// GetMediaList returns GetWatchingPage.MediaList, and is useful for accessing the field via an interface.
func (v *GetWatchingPage) GetMediaList() []GetWatchingPageMediaList { return v.MediaList }

// GetWatchingPageMediaList includes the requested fields of the GraphQL type MediaList.
// The GraphQL type's documentation follows.
//
// List of anime or manga
type GetWatchingPageMediaList struct {
	Media GetWatchingPageMediaListMedia `json:"media"`
}

// GetMedia returns GetWatchingPageMediaList.Media, and is useful for accessing the field via an interface.
func (v *GetWatchingPageMediaList) GetMedia() GetWatchingPageMediaListMedia { return v.Media }

// GetWatchingPageMediaListMedia includes the requested fields of the GraphQL type Media.
// The GraphQL type's documentation follows.
//
// Anime or Manga
type GetWatchingPageMediaListMedia struct {
	// The id of the media
	Id int `json:"id"`
	// The mal id of the media
	IdMal int `json:"idMal"`
	// The official titles of the media in various languages
	Title GetWatchingPageMediaListMediaTitle `json:"title"`
}

// GetId returns GetWatchingPageMediaListMedia.Id, and is useful for accessing the field via an interface.
func (v *GetWatchingPageMediaListMedia) GetId() int { return v.Id }

// GetIdMal returns GetWatchingPageMediaListMedia.IdMal, and is useful for accessing the field via an interface.
func (v *GetWatchingPageMediaListMedia) GetIdMal() int { return v.IdMal }

// GetTitle returns GetWatchingPageMediaListMedia.Title, and is useful for accessing the field via an interface.
func (v *GetWatchingPageMediaListMedia) GetTitle() GetWatchingPageMediaListMediaTitle { return v.Title }

// GetWatchingPageMediaListMediaTitle includes the requested fields of the GraphQL type MediaTitle.
// The GraphQL type's documentation follows.
//
// The official titles of the media in various languages
type GetWatchingPageMediaListMediaTitle struct {
	// The romanization of the native language title
	Romaji string `json:"romaji"`
}

// GetRomaji returns GetWatchingPageMediaListMediaTitle.Romaji, and is useful for accessing the field via an interface.
func (v *GetWatchingPageMediaListMediaTitle) GetRomaji() string { return v.Romaji }

// GetWatchingResponse is returned by GetWatching on success.
type GetWatchingResponse struct {
	Page GetWatchingPage `json:"Page"`
}

// GetPage returns GetWatchingResponse.Page, and is useful for accessing the field via an interface.
func (v *GetWatchingResponse) GetPage() GetWatchingPage { return v.Page }

// __GetUserByNameInput is used internally by genqlient
type __GetUserByNameInput struct {
	Name string `json:"name"`
}

// GetName returns __GetUserByNameInput.Name, and is useful for accessing the field via an interface.
func (v *__GetUserByNameInput) GetName() string { return v.Name }

// __GetWatchingInput is used internally by genqlient
type __GetWatchingInput struct {
	UserId  int `json:"userId"`
	Page    int `json:"page"`
	PerPage int `json:"perPage"`
}

// GetUserId returns __GetWatchingInput.UserId, and is useful for accessing the field via an interface.
func (v *__GetWatchingInput) GetUserId() int { return v.UserId }

// GetPage returns __GetWatchingInput.Page, and is useful for accessing the field via an interface.
func (v *__GetWatchingInput) GetPage() int { return v.Page }

// GetPerPage returns __GetWatchingInput.PerPage, and is useful for accessing the field via an interface.
func (v *__GetWatchingInput) GetPerPage() int { return v.PerPage }

// The query or mutation executed by GetUserByName.
const GetUserByName_Operation = `
query GetUserByName ($name: String!) {
	User(name: $name) {
		id
	}
}
`

func GetUserByName(
	ctx_ context.Context,
	client_ graphql.Client,
	name string,
) (*GetUserByNameResponse, error) {
	req_ := &graphql.Request{
		OpName: "GetUserByName",
		Query:  GetUserByName_Operation,
		Variables: &__GetUserByNameInput{
			Name: name,
		},
	}
	var err_ error

	var data_ GetUserByNameResponse
	resp_ := &graphql.Response{Data: &data_}

	err_ = client_.MakeRequest(
		ctx_,
		req_,
		resp_,
	)

	return &data_, err_
}

// The query or mutation executed by GetWatching.
const GetWatching_Operation = `
query GetWatching ($userId: Int!, $page: Int!, $perPage: Int!) {
	Page(page: $page, perPage: $perPage) {
		mediaList(userId: $userId, type: ANIME, status_in: [CURRENT,PLANNING]) {
			media {
				id
				idMal
				title {
					romaji
				}
			}
		}
	}
}
`

func GetWatching(
	ctx_ context.Context,
	client_ graphql.Client,
	userId int,
	page int,
	perPage int,
) (*GetWatchingResponse, error) {
	req_ := &graphql.Request{
		OpName: "GetWatching",
		Query:  GetWatching_Operation,
		Variables: &__GetWatchingInput{
			UserId:  userId,
			Page:    page,
			PerPage: perPage,
		},
	}
	var err_ error

	var data_ GetWatchingResponse
	resp_ := &graphql.Response{Data: &data_}

	err_ = client_.MakeRequest(
		ctx_,
		req_,
		resp_,
	)

	return &data_, err_
}
