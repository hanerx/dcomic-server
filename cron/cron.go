package cron

import (
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
		ip, ipErr := utils.GetExternalIP()
		if ipErr == nil {
			server := ip.String()
			if os.Getenv("hostname") != "" {
				server = os.Getenv("hostname")
			}
			for i := 0; i < len(nodes); i++ {
				response, httpErr := http.Get(fmt.Sprintf("http://%s/server/%s", nodes[i].Address, server))
				if httpErr == nil && response.StatusCode == 200 {
					body, readErr := ioutil.ReadAll(response.Body)
					if readErr == nil {
						type Response struct {
							Code int                 `json:"code"`
							Msg  string              `json:"msg"`
							Data []model.ComicDetail `json:"data"`
						}
						var jsonData Response
						err = json.Unmarshal(body, &jsonData)
						comics := jsonData.Data
						if err == nil {
							failed := 0
							success := 0
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
								}
							}
							log.Printf("完成与%s同步,成功：%d,失败: %d", nodes[i].Address, success, failed)
						} else {
							log.Println(err)
						}
					} else {
						log.Println(readErr)
					}
				}
			}
		} else {
			log.Println(ipErr)
		}
	} else {
		log.Println(err)
	}
}
