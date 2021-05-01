package cron

import (
	"bytes"
	"dcomicServer/database"
	"dcomicServer/model"
	"dcomicServer/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func AutoSync() {
	var nodes []model.Node
	err := database.Databases.C("server").Find(nil).All(&nodes)
	if err == nil {
		switch os.Getenv("server_mode") {
		case "0":
			return
		case "1":
			PutSyncComic(nodes)
			break
		case "2":
			GetSyncComic(nodes)
			break
		case "3":
			GetSyncComic(nodes)
			PutSyncComic(nodes)
			break
		default:
			return
		}
	} else {
		log.Println(err)
	}
}

func GetSyncComic(nodes []model.Node) {
	ip, ipErr := utils.GetExternalIP()
	if ipErr == nil {
		server := ip.String()
		if os.Getenv("hostname") != "" {
			server = os.Getenv("hostname")
		}
		for i := 0; i < len(nodes); i++ {
			response, httpErr := http.Get(fmt.Sprintf("http://%s/server/node/%s/comic", nodes[i].Address, server))
			if httpErr == nil && response.StatusCode == 200 {
				body, readErr := ioutil.ReadAll(response.Body)
				if readErr == nil {
					type Response struct {
						Code int                 `json:"code"`
						Msg  string              `json:"msg"`
						Data []model.ComicDetail `json:"data"`
					}
					var jsonData Response
					err := json.Unmarshal(body, &jsonData)
					comics := jsonData.Data
					if err == nil {
						failed := 0
						success := 0
						ignore := 0
						for j := 0; j < len(comics); j++ {
							var comic model.ComicDetail
							err = database.Databases.C("comic").Find(map[string]string{"comic_id": comics[j].ComicId}).One(&comic)
							if err == nil && comic.Timestamp < comics[i].Timestamp {
								err = database.Databases.C("comic").Update(map[string]string{"comic_id": comics[j].ComicId}, comics[j])
								if err == nil {
									success++
								} else {
									failed++
								}
							} else if err != nil {
								err = database.Databases.C("comic").Insert(comics[j])
								if err == nil {
									success++
								} else {
									failed++
								}
							} else {
								ignore++
							}
						}
						log.Printf("完成从 %s 的同步操作,总计: %d,成功：%d,失败: %d,忽略: %d", nodes[i].Address, len(comics), success, failed, ignore)
					} else {
						log.Println(err)
					}
				} else {
					log.Println(readErr)
				}
			} else {
				log.Println(httpErr)
			}
		}
	} else {
		log.Println(ipErr)
	}
}

func PutSyncComic(nodes []model.Node) {
	var comics []model.ComicDetail
	err := database.Databases.C("comic").Find(nil).All(&comics)
	if err == nil {
		ip, ipErr := utils.GetExternalIP()
		if ipErr == nil {
			data := make([]model.ComicDetail, len(comics))
			for i := 0; i < len(comics); i++ {
				if comics[i].Redirect {
					data[i] = comics[i]
				} else {
					if os.Getenv("hostname") == "" {
						data[i] = model.ComicDetail{
							Title:       comics[i].Title,
							Timestamp:   comics[i].Timestamp,
							Description: comics[i].Description,
							Cover:       comics[i].Cover,
							Tags:        comics[i].Tags,
							Authors:     comics[i].Authors,
							ComicId:     comics[i].ComicId,
							Redirect:    true,
							RedirectUrl: ip.String(),
						}
					} else {
						data[i] = model.ComicDetail{
							Title:       comics[i].Title,
							Timestamp:   comics[i].Timestamp,
							Description: comics[i].Description,
							Cover:       comics[i].Cover,
							Tags:        comics[i].Tags,
							Authors:     comics[i].Authors,
							ComicId:     comics[i].ComicId,
							Redirect:    true,
							RedirectUrl: os.Getenv("hostname"),
						}
					}
				}
			}
			server := ip.String()
			if os.Getenv("hostname") != "" {
				server = os.Getenv("hostname")
			}
			for i := 0; i < len(nodes); i++ {

				jsonData, jsonErr := json.Marshal(data)
				if jsonErr == nil {
					client := &http.Client{}
					request, httpErr := http.NewRequest("POST", fmt.Sprintf("http://%s/server/node/%s/comic", nodes[i].Address, server), bytes.NewBuffer(jsonData))
					if httpErr == nil {
						request.Header.Add("Content-Type", "application/json")
						request.Header.Add("token", nodes[i].Token)
						response, responseErr := client.Do(request)
						if responseErr == nil {
							body, readErr := ioutil.ReadAll(response.Body)
							if readErr == nil {
								type Data struct {
									Success int `json:"success"`
									Failed  int `json:"failed"`
									Ignore  int `json:"ignore"`
								}
								type Response struct {
									Code int    `json:"code"`
									Msg  string `json:"msg"`
									Data Data   `json:"data"`
								}
								var responseData Response
								err = json.Unmarshal(body, &responseData)
								if err == nil {
									log.Printf("完成向 %s 的同步操作,总计: %d,成功：%d,失败: %d,忽略: %d", nodes[i].Address, len(comics), responseData.Data.Success, responseData.Data.Failed, responseData.Data.Ignore)
								} else {
									log.Println(err)
								}
							} else {
								log.Println(readErr)
							}
						} else {
							log.Println(responseErr)
						}
					} else {
						log.Println(httpErr)
					}
				} else {
					log.Println(jsonErr)
				}
			}
		} else {
			log.Println(ipErr)
		}
	} else {
		log.Println(err)
	}
}

func GetSyncUser(nodes []model.Node) {

}

func PutSyncUser(nodes []model.Node) {

}

func GetSyncServer(nodes []model.Node) {

}

func PutSyncServer(nodes []model.Node) {

}
