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
}

type ComicDetailGetter interface {
	GetGroup(groupId string) (ComicGroup, error)
	AddGroup(group ComicGroup)
	UpdateGroup(groupId string, group ComicGroup) error
	DeleteGroup(groupId string) error
}

type ComicGroupGetter interface {
	GetChapter(chapterId string) (ComicChapter, error)
	AddChapter(chapter ComicChapter)
	UpdateChapter(chapterId string, chapter ComicChapter) error
	DeleteChapter(chapterId string) error
}

func (comicDetail *ComicDetail) GetGroup(groupId string) (ComicGroup, error) {
	for i := 0; i < len(comicDetail.Groups); i++ {
		group := comicDetail.Groups[i]
		if group.GroupId == groupId {
			return group, nil
		}
	}
	return ComicGroup{}, errors.New(fmt.Sprintf("no group named %s", groupId))
}

func (comicDetail *ComicDetail) AddGroup(group ComicGroup) {
	comicDetail.Groups = append(comicDetail.Groups, group)
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

func (group *ComicGroup) GetChapter(chapterId string) (ComicChapter, error) {
	for i := 0; i < len(group.Chapters); i++ {
		chapter := group.Chapters[i]
		if chapter.ChapterId == chapterId {
			return chapter, nil
		}
	}
	return ComicChapter{}, errors.New(fmt.Sprintf("no chapter named %s", chapterId))
}

func (group *ComicGroup) AddChapter(chapter ComicChapter) {
	group.Chapters = append(group.Chapters, chapter)
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
