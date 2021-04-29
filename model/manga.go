package model

import (
	"errors"
	"fmt"
)

type ComicChapter struct {
	Title     string   `json:"title" bson:"title"`
	ChapterId string   `json:"chapter_id" bson:"chapter_id"`
	ComicId   string   `json:"comic_id" bson:"comic_id"`
	Pages     []string `json:"data" bson:"pages"`
	Timestamp int      `json:"timestamp" bson:"timestamp"`
}

type ComicGroup struct {
	Title    string         `json:"title" bson:"title"`
	GroupId  string         `json:"name" bson:"group_id"`
	Chapters []ComicChapter `json:"data" bson:"chapters"`
}



type ComicDetail struct {
	Title       string       `json:"title" bson:"title"`
	Cover       string       `json:"cover" bson:"cover"`
	Description string       `json:"description" bson:"description"`
	ComicId     string       `json:"comic_id" bson:"comic_id"`
	Groups      []ComicGroup `json:"data" bson:"groups"`
	Timestamp   int64        `json:"timestamp" bson:"timestamp"`
	Redirect    bool         `json:"redirect" bson:"redirect"`
	RedirectUrl string       `json:"redirect_url" bson:"redirect_url"`
	Tags        []ComicTag   `json:"tags" bson:"tags"`
	Authors     []ComicTag   `json:"authors" bson:"authors"`
	Status      string       `json:"status" bson:"status"`
}

type ComicDetailGetter interface {
	GetGroup(groupId string) (*ComicGroup, error)
	AddGroup(group ComicGroup) error
	UpdateGroup(groupId string, group ComicGroup) error
	DeleteGroup(groupId string) error
	GetChapter(groupId string, chapterId string) (*ComicChapter, error)
	AddChapter(groupId string, chapter ComicChapter) error
	UpdateChapter(groupId string, chapterId string, chapter ComicChapter) error
	DeleteChapter(groupId string, chapterId string) error
}

type ComicGroupGetter interface {
	GetChapter(chapterId string) (*ComicChapter, error)
	AddChapter(chapter ComicChapter) error
	UpdateChapter(chapterId string, chapter ComicChapter) error
	DeleteChapter(chapterId string) error
}

func (comicDetail *ComicDetail) GetGroup(groupId string) (*ComicGroup, error) {
	for i := 0; i < len(comicDetail.Groups); i++ {
		group := comicDetail.Groups[i]
		if group.GroupId == groupId {
			return &comicDetail.Groups[i], nil
		}
	}
	return &ComicGroup{}, errors.New(fmt.Sprintf("no group named %s", groupId))
}

func (comicDetail *ComicDetail) AddGroup(group ComicGroup) error {
	_, err := comicDetail.GetGroup(group.GroupId)
	if err != nil {
		comicDetail.Groups = append(comicDetail.Groups, group)
		return nil
	}
	return errors.New("group already exist")
}

func (comicDetail *ComicDetail) UpdateGroup(groupId string, group ComicGroup) error {
	for i := 0; i < len(comicDetail.Groups); i++ {
		originGroup := comicDetail.Groups[i]
		if originGroup.GroupId == groupId {
			comicDetail.Groups[i] = group
			return nil
		}
	}
	return errors.New(fmt.Sprintf("no group named %s", groupId))
}

func (comicDetail *ComicDetail) DeleteGroup(groupId string) error {
	for i := 0; i < len(comicDetail.Groups); i++ {
		originGroup := comicDetail.Groups[i]
		if originGroup.GroupId == groupId {
			comicDetail.Groups = append(comicDetail.Groups[:i], comicDetail.Groups[i+1:]...)
			return nil
		}
	}
	return errors.New(fmt.Sprintf("no group named %s", groupId))
}

func (comicDetail *ComicDetail) GetChapter(groupId string, chapterId string) (*ComicChapter, error) {
	group, err := comicDetail.GetGroup(groupId)
	if err != nil {
		return nil, err
	}
	return group.GetChapter(chapterId)
}

func (comicDetail *ComicDetail) AddChapter(groupId string, chapter ComicChapter) error {
	group, err := comicDetail.GetGroup(groupId)
	if err != nil {
		return err
	}
	return group.AddChapter(chapter)
}

func (comicDetail *ComicDetail) UpdateChapter(groupId string, chapterId string, chapter ComicChapter) error {
	group, err := comicDetail.GetGroup(groupId)
	if err != nil {
		return err
	}
	return group.UpdateChapter(chapterId, chapter)
}

func (comicDetail *ComicDetail) DeleteChapter(groupId string, chapterId string) error {
	group, err := comicDetail.GetGroup(groupId)
	if err != nil {
		return err
	}
	return group.DeleteChapter(chapterId)
}

func (group *ComicGroup) GetChapter(chapterId string) (*ComicChapter, error) {
	for i := 0; i < len(group.Chapters); i++ {
		chapter := group.Chapters[i]
		if chapter.ChapterId == chapterId {
			return &group.Chapters[i], nil
		}
	}
	return nil, errors.New(fmt.Sprintf("no chapter named %s", chapterId))
}

func (group *ComicGroup) AddChapter(chapter ComicChapter) error {
	_, err := group.GetChapter(chapter.ChapterId)
	if err == nil {
		return errors.New("chapter already exist")
	}
	group.Chapters = append(group.Chapters, chapter)
	return nil
}

func (group *ComicGroup) UpdateChapter(chapterId string, chapter ComicChapter) error {
	for i := 0; i < len(group.Chapters); i++ {
		originChapter := group.Chapters[i]
		if originChapter.ChapterId == chapterId {
			group.Chapters[i] = chapter
			return nil
		}
	}
	return errors.New(fmt.Sprintf("no chapter named %s", chapterId))
}

func (group *ComicGroup) DeleteChapter(chapterId string) error {
	for i := 0; i < len(group.Chapters); i++ {
		originChapter := group.Chapters[i]
		if originChapter.ChapterId == chapterId {
			group.Chapters = append(group.Chapters[:i], group.Chapters[i+1:]...)
			return nil
		}
	}
	return errors.New(fmt.Sprintf("no chapter named %s", chapterId))
}
