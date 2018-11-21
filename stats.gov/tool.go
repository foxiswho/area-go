package stats_gov

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
	"github.com/foxiswho/area-go/consts"
	"github.com/foxiswho/area-go/mca.gov"
	"github.com/foxiswho/area-go/util"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//被爬地址
const SITE = "http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2017/index.html"

//存储变量
var area map[int]string
var areaSort = []int{}
//格式化后存储
var areaFormat map[int]map[int]string
//顶级ID
const topLevelId = 86
const fileSuffix = ".html"

//临时文件名
const tmp_file_index_name = "stats.index.html"

//
var collect_url map[int]string
//
var site_path = ""
var is_test = false

//数据 获取
func GetAreaData() {
	path, _ := util.GetCurrentPath()
	fmt.Println("当前可执行文件路径：", path)
	fmt.Println("当前可执行文件路径2：", util.GetCurrentPath2())
	fmt.Println("当前可执行文件路径3：", consts.APP_PATH)
	//time.Sleep(time.Duration(5) * time.Second)

	getProvince()
	////
	fmt.Println("collect_url", collect_url)
	if len(collect_url) > 0 {
		fmt.Println("第2级：")
		fmt.Println("第2级：")
		fmt.Println("第2级：")
		fmt.Println("第2级：")
		fmt.Println("第2级：")
		for _, url := range collect_url {
			fmt.Println(url)
		}
		collect_city := make(map[int]string)
		i := 1
		for _, url := range collect_url {
			data := getCity(url, 2)
			if is_test {
				i = i + 1
				if i > 2 {
					break
				}
			}
			if data != nil {
				for key, val := range data {
					collect_city[key] = val

				}
			}
		}
		if len(collect_city) > 0 {
			fmt.Println("第3级：")
			fmt.Println("第3级：")
			fmt.Println("第3级：")
			fmt.Println("第3级：")
			fmt.Println("第3级：")
			fmt.Println("第3级：")
			for _, url := range collect_city {
				fmt.Println(url)
			}
			collect_district := make(map[int]string)
			i := 1
			for _, url := range collect_city {
				data := getCity(url, 3)
				if is_test {
					i = i + 1
					if i > 2 {
						break
					}
				}
				if data != nil {

					for key, val := range data {
						collect_district[key] = val
					}
				}
			}
			if len(collect_district) > 0 {
				fmt.Println("第4级：")
				fmt.Println("第4级：")
				fmt.Println("第4级：")
				fmt.Println("第4级：")
				fmt.Println("第4级：")
				fmt.Println("第4级：")
				fmt.Println("第4级：")
				for _, url := range collect_district {
					fmt.Println(url)
				}
				for key, value := range area {
					fmt.Println(key, "=:::", value)
				}
				for _, url := range collect_district {
					data := getCity(url, 4)
					fmt.Println("第4级：", data)
				}
			}
		}
	}
	fmt.Println(len(area))
	fmt.Println("area:", area)
}

func getProvince() {
	input_byte := []byte{}
	//检查目录是否已创建
	is_dir, _ := util.PathExists(consts.APP_PATH + consts.TMP_DIR)
	if false == is_dir {
		// 创建文件夹
		err := os.Mkdir(consts.APP_PATH+consts.TMP_DIR, os.ModePerm)
		if err == nil {
			fmt.Printf("mkdir success!\n")
		}
	}
	tmp_index_file := consts.APP_PATH + consts.TMP_DIR + "/" + tmp_file_index_name
	fmt.Println("首页文件路径：", tmp_index_file)
	//
	is_file, _ := util.PathExists(tmp_index_file)
	if false == is_file {
		aConv, err := iconv.NewConverter("GBK", "utf-8")
		if err != nil {
			fmt.Println("iconv.Open failed!")
			return
		}
		defer aConv.Close()

		res, err := http.Get(SITE)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		}
		input, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println("err::::", err)
		}
		if len(input) < 100 {
			fmt.Println("字符：", input)
			fmt.Println("错误：ioutil.ReadAll")
			return
		}
		str, err := aConv.ConvertString(string(input))
		fmt.Println("字符：", str)
		util.SaveFile([]byte(str), tmp_index_file)

		input_byte = []byte(str)
	} else {
		b, err := ioutil.ReadFile(tmp_index_file)
		if err != nil {
			fmt.Print(err)
		}
		input_byte = b
	}

	//fmt.Println("body::::", str)
	//i, err := res.Body.Read([]byte(str))
	fmt.Println("字符：", len(input_byte))
	if len(input_byte) < 100 {
		fmt.Println("字符：", input_byte)
		fmt.Println("错误：aConv.ConvertString")
		return
	}
	reader := bytes.NewReader(input_byte)
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
		fmt.Println("错误：goquery.NewDocumentFromReader")
		return
	}
	//fmt.Println(doc)
	collect_url = make(map[int]string)
	area = make(map[int]string)
	files := SITE
	paths, fileName := filepath.Split(files)
	fmt.Println("路径：", paths, " 文件名：", fileName)
	site_path = paths

	doc.Find(".provincetable").Find("a").Each(func(i int, item *goquery.Selection) {
		link, _ := item.Eq(0).Attr("href")
		fmt.Println(link, " = ", item.Eq(0).Text())
		// 去除文件名后缀
		file := strings.TrimSuffix(link, fileSuffix)
		fmt.Println("文件去除后缀：", file)
		if len(file) > 0 {
			key, _ := strconv.Atoi(file)
			area[key*10] = item.Eq(0).Text()
			areaSort = append(areaSort, key)
			//
			collect_url[key] = paths + link
		}

	})
	//fmt.Println(area)
	fmt.Println(len(area))
}

func getCity(url string, level int) map[int]string {

	fmt.Println("url：", url)
	paths, fileName := filepath.Split(url)
	paths2 := strings.Replace(paths, site_path, "", -1)
	fmt.Println("路径：", paths2, " 文件名：", fileName)
	tmp_path := consts.APP_PATH + consts.TMP_DIR + "/" + paths2
	//检查目录是否已创建
	is_dir, _ := util.PathExists(tmp_path)
	if false == is_dir {
		// 创建文件夹
		err := os.Mkdir(tmp_path, os.ModePerm)
		if err == nil {
			fmt.Printf("mkdir success!\n")
		}
	}
	input_byte := []byte{}
	fmt.Println("文件是否存在，文件地址：", tmp_path+fileName)
	//检查文件是否已存在
	is_file, _ := util.PathExists(tmp_path + fileName)
	if false == is_file {
		fmt.Println("Sleep 2 Second begin")
		time.Sleep(time.Duration(1) * time.Second)
		fmt.Println("end")
		//
		fmt.Println("级别：", level)
		fmt.Println("级别：", level)
		fmt.Println("级别：", level)
		fmt.Println("级别：", level)
		fmt.Println("级别：", level)
		fmt.Println("级别：", level)
		fmt.Println("级别：", level)
		fmt.Println("级别：", level)
		aConv, err := iconv.NewConverter("GBK", "utf-8")
		if err != nil {
			fmt.Println("iconv.Open failed!")
			return nil
		}
		defer aConv.Close()
		fmt.Println("URL:", url)
		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
			return nil
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		}
		input, err := ioutil.ReadAll(res.Body)
		str, err := aConv.ConvertString(string(input))
		//fmt.Println("字符：", str)
		util.SaveFile([]byte(str), tmp_path+fileName)

		input_byte = []byte(str)
	} else {
		fmt.Println("文件已存在，直接读取该文件：", tmp_path+fileName)
		b, err := ioutil.ReadFile(tmp_path + fileName)
		if err != nil {
			fmt.Print(err)
		}
		input_byte = b
	}

	//fmt.Println("body::::", str)
	//i, err := res.Body.Read([]byte(str))
	//fmt.Println(i, err)
	reader := bytes.NewReader(input_byte)
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(doc)
	//area = make(map[int]string)
	fmt.Println("路径：", paths, " 文件名：", fileName)
	data_url := make(map[int]string)
	select_str := ".provincetable"
	if level == 2 {
		select_str = ".citytable"
	} else if level == 3 {
		select_str = ".towntable"
	} else if level == 4 {
		select_str = ".villagetable"
	}
	doc.Find(select_str).Find("tr").Each(func(i int, item *goquery.Selection) {
		str := item.Find("td").Eq(0).Find("a").Text()
		value := item.Find("td").Eq(1).Text()
		link := ""
		key := ""
		if str == "" {
			str = item.Find("td").Eq(0).Text()
			if str != "" {
				match, _ := regexp.MatchString(`\d`, str)
				if match {
					key_int, _ := strconv.Atoi(str)
					area[key_int] = value
					fmt.Println("市辖区：", str, area[key_int])
				}
			}
		} else {
			if level == 4 {
				key = item.Find("td").Eq(0).Text()
				value = item.Find("td").Eq(1).Text()
			} else {
				link, _ = item.Find("td").Eq(0).Find("a").Attr("href")
				key = item.Find("td").Eq(0).Find("a").Text()
				value = item.Find("td").Eq(1).Find("a").Text()
			}

			//fmt.Println(link, " = ", value)
			fmt.Println(key, " = ", value)
			// 去除文件名后缀
			file := strings.TrimSuffix(link, fileSuffix)
			fmt.Println("文件去除后缀：", file)
			if len(file) > 0 {
				key_int, _ := strconv.Atoi(key)
				area[key_int] = value
				areaSort = append(areaSort, key_int)
				//
				data_url[key_int] = paths + link
				//collect_url[key] = paths + "/" + link
			}
		}

	})
	if len(data_url) < 1 {
		return nil
	}
	//fmt.Println(area)
	fmt.Println(len(area))
	return data_url
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

func SaveFile() {
	str := util.ObjToJson(&areaFormat)
	//
	util.SaveFile(str, consts.JSON_FILE)
}

//扩展数据填充
func FormatExtData() {
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

func SaveSqlFile() {
	str := ""
	str = util.GetCreateTableSql() + "\n\n"
	str += util.MakeSql(areaFormat)
	//Sql
	util.SaveFile([]byte(str), consts.SQL_FILE)
}

func SaveCsvFile() {
	str := ""
	str += util.MakeCsv(areaFormat)
	//Sql
	util.SaveFile([]byte(str), consts.CSV_FILE)
}
