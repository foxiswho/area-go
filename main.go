package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/foxiswho/area-go/mca.gov"
	"github.com/foxiswho/area-go/util"
	"log"
	"math"
	"net/http"
	"strconv"
)

//被爬地址
const SITE = "http://www.mca.gov.cn/article/sj/xzqh/2018/201804-12/20180708230813.html"

//存储变量
var area map[int]string
var areaSort = []int{}
//格式化后存储
var areaFormat map[int]map[int]string
//顶级ID
const topLevelId = 86

const JSON_FILE = "area.js"

//数据 获取
func GetAreaData() {
	res, err := http.Get(SITE)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	area = make(map[int]string)

	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		if node := s.Find(".xl708733").Nodes; node != nil {
			if len(node) >= 2 {
				//fmt.Println(node[0]," = ",node[1])
				//fmt.Println(s.Find("td.xl708733").Eq(0).Text())
				id := s.Find("td.xl708733").Eq(0).Text()
				if id != "" {
					key, _ := strconv.Atoi(id)
					value := s.Find("td.xl708733").Eq(1).Text()
					area[key] = value
					areaSort = append(areaSort, key)
				}
			}
		}
	})
	//fmt.Println(area)
	fmt.Println(len(area))
}

//格式化数据
func FormatArea() {
	//初始化
	areaFormat = make(map[int]map[int]string)
	//
	for index := range areaSort {
		key := areaSort[index]
		value := area[key]
		modTmp := key % 1000
		//省份
		if modTmp == 0 {
			//是否初始化
			if areaFormat[topLevelId] == nil {
				areaFormat[topLevelId] = make(map[int]string)
			}
			maps := areaFormat[topLevelId]
			maps[key] = value
			areaFormat[topLevelId] = maps
			//直辖市处理
			formatMunicipality(key)
		} else { //市县
			modCity := modTmp % 100
			//市
			if modCity == 0 {
				tmp := math.Floor(float64(key) / 1000)
				province := int(tmp) * 1000
				if areaFormat[province] == nil {
					areaFormat[province] = make(map[int]string)
				}
				maps := areaFormat[province]
				maps[key] = value
				areaFormat[province] = maps
			} else {
				//县区
				tmp := math.Floor(float64(key) / 100)
				city := int(tmp) * 100
				if areaFormat[city] == nil {
					areaFormat[city] = make(map[int]string)
				}
				maps := areaFormat[city]
				maps[key] = value
				areaFormat[city] = maps
			}
		}
		//break
	}

	//fmt.Println(areaFormat)
	fmt.Println(len(areaFormat))
	//str, err := json.Marshal(&areaFormat)
	//if err != nil {
	//	fmt.Println("序列化失败:", err)
	//}
	//fmt.Println("JSON:", str)
}

//直辖市处理
func formatMunicipality(id int) {
	province := int(id/1000) * 1000
	city := int((id/100)+1) * 100
	//直辖市
	municipality := make(map[int]string)
	municipality[110100] = "北京市"
	municipality[120100] = "天津市"
	municipality[310100] = "上海市"
	municipality[500100] = "重庆市"
	if ok := municipality[city]; ok != "" {
		if areaFormat[province] == nil {
			areaFormat[province] = make(map[int]string)
		}
		maps := areaFormat[province]
		maps[city] = ok
		areaFormat[province] = maps
	}
}

func saveFile() {
	str := util.ObjToJson(&areaFormat)
	//
	util.SaveFile(str, JSON_FILE)
}

//扩展数据填充
func formatExtData() {
	//
	mca_gov.ReadExt()
	//
	for key, _ := range mca_gov.AreaFormatExt {
		if key > 86 {
			areaFormat[key] = mca_gov.AreaFormatExt[key]
		}
	}

	//str := util.ObjToJson(&areaFormat)
	////
	//util.SaveFile(str, "ext.js")
}

func main() {
	fmt.Println("=======获取数据======")
	GetAreaData()
	fmt.Println("=======获取成功=====")
	fmt.Println("=======格式化数据=====")
	//
	FormatArea()
	fmt.Println("=======格式化 扩展数据=====")
	//
	formatExtData()
	fmt.Println("=======处理成功=====")
	//
	fmt.Println("===================")
	//
	fmt.Println("=======  写入到文件 " + JSON_FILE + "    ============")
	fmt.Println("=======  路径: $GOPATH/src/github.com/foxiswho/area-go/" + JSON_FILE + "    ============")
	saveFile()

	fmt.Println("===================")
	fmt.Println("=======数据保存成功=======")
}
