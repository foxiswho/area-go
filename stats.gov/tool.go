package stats_gov

import (
	"bytes"
	"encoding/json"
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
const SITE = "http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2018/index.html"

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

//暂停爬取秒数
const pause_second_web_crawler = 1

//
var collect_url map[int]string
//
var site_path = ""
var is_test = false
//是否使用 指定测试URL
var is_test_url = false
//没有获取到数据的连接
var url_data_null = []string{}
var city_is_over = false

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
		//是否指定读取URL
		if is_test_url {
			url := "http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2018/11.html"
			data := getCity(url, 2)
			if data != nil {
				for key, val := range data {
					collect_city[key] = val
				}
			}
		} else {
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
				if data == nil {
					url_data_null = append(url_data_null, url)
				} else if len(data) == 0 {
					url_data_null = append(url_data_null, url)
				}
			}
		}

		//fmt.Println("collect_city:", collect_city)
		//fmt.Println("area:", area)
		//return

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

			//是否指定读取URL
			if is_test_url {
				url := "http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2018/11/1101.html"
				data := getCity(url, 3)
				if data != nil {
					for key, val := range data {
						collect_district[key] = val
					}
				}
			} else {
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
					if data == nil {
						url_data_null = append(url_data_null, url)
					} else if len(data) == 0 {
						url_data_null = append(url_data_null, url)
					}
				}
			}

			//fmt.Println("collect_district:", collect_district)
			//fmt.Println("area:", area)
			//return
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
				//是否指定读取URL
				if is_test_url {
					url := "http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2017/44/19/441900104.html"
					data := getCity(url, 4)
					if data != nil {
						for key, val := range data {
							collect_district[key] = val
						}
					}
				} else {
					for _, url := range collect_district {
						city_is_over = false
						data := getCity(url, 4)
						fmt.Println("第4级：", data)
						if city_is_over == false {
							if data == nil {
								url_data_null = append(url_data_null, url)
							} else if len(data) == 0 {
								url_data_null = append(url_data_null, url)
							}
						}
					}
				}

			}
		}
	}
	fmt.Println(len(area))
	fmt.Println("area:", area)
	fmt.Println("url data null")
	fmt.Println("url data null")
	fmt.Println("url data null")
	fmt.Println("url data null")
	fmt.Println("url data null")
	fmt.Println("url_data_null:", url_data_null)
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
			area[key] = item.Eq(0).Text()
			areaSort = append(areaSort, key)
			//
			collect_url[key] = paths + link
		}

	})
	fmt.Println("area:", area)
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
	fmt.Println("准备爬取的文件地址：", tmp_path+fileName)
	//检查文件是否已存在
	is_file, _ := util.PathExists(tmp_path + fileName)
	if false == is_file {

		fmt.Println(fmt.Sprintf("暂停 %d 秒后后开始爬取", pause_second_web_crawler))
		fmt.Println(fmt.Sprintf("Sleep %d Second begin", pause_second_web_crawler))
		time.Sleep(time.Duration(pause_second_web_crawler) * time.Second)
		fmt.Println("end")
		fmt.Println("开始爬取该URl")
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
		fmt.Println("爬取完成")
		fmt.Println("检测文件数据")
		if len(str) > 100 {
			fmt.Println("保存文件:", tmp_path+fileName)
			util.SaveFile([]byte(str), tmp_path+fileName)
			input_byte = []byte(str)
		} else {
			fmt.Println("文件数据错误")
			return nil
		}

	} else {
		fmt.Println("文件已存在，直接读取该文件：", tmp_path+fileName)
		b, err := ioutil.ReadFile(tmp_path + fileName)
		if err != nil {
			fmt.Print(err)
		}
		if len(b) < 100 {
			fmt.Println("文件数据错误:", tmp_path+fileName)
			os.Remove(tmp_path + fileName)
			return nil
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
		// 第三级   市下面的 区县
		select_str = ".countytable"
	} else if level == 4 {
		select_str = ".towntable"
	} else if level == 5 {
		select_str = ".villagetable"
	}
	//海南
	if level == 3 && doc.Find(select_str).Find("tr").Nodes == nil {
		level = 4
		select_str = ".towntable"
	} else if level == 4 && doc.Find(select_str).Find("tr").Nodes == nil {
		//中山
		level = 5
		select_str = ".villagetable"
		city_is_over = true
	}
	//是否最后一级，居委会，如果是，不再执行后续
	is_villagetable := false
	doc.Find(select_str).Find("tr").Each(func(i int, item *goquery.Selection) {
		//if is_villagetable {
		//	return
		//}
		str := item.Find("td").Eq(0).Find("a").Text()
		value := item.Find("td").Eq(1).Text()
		link := ""
		key := item.Find("td").Eq(0).Text()
		if key == "统计用区划代码" {
			//跳过第一行标题
			return
		}
		//key = item.Find("td").Eq(1).Text()
		//if key == "城乡分类代码" {
		//	//跳过
		//	is_villagetable = true
		//	return
		//}

		//第4级
		if level == 4 {
			key = item.Find("td").Eq(0).Text()
			value = item.Find("td").Eq(1).Text()
			link, _ = item.Find("td").Eq(0).Find("a").Attr("href")
			// 去除文件名后缀
			fmt.Println("link：", link)
			file := strings.TrimSuffix(link, fileSuffix)
			fmt.Println("文件去除后缀：", file)
			fmt.Println(key, " = ", value)

			if len(file) > 0 {
				key_int, _ := strconv.Atoi(key)
				key_int = formatKeyint(level, key_int)
				area[key_int] = value
				areaSort = append(areaSort, key_int)
				//
				data_url[key_int] = paths + link
				//collect_url[key] = paths + "/" + link
			}
		} else if level == 5 {
			key = item.Find("td").Eq(0).Text()
			value = item.Find("td").Eq(2).Text()
			fmt.Println(key, " = ", value)
			key_int, _ := strconv.Atoi(key)
			key_int = formatKeyint(level, key_int)
			area[key_int] = value
			areaSort = append(areaSort, key_int)
		} else if str == "" {
			str = item.Find("td").Eq(0).Text()
			if str != "" {
				match, _ := regexp.MatchString(`\d`, str)
				if match {
					key_int, _ := strconv.Atoi(str)
					key_int = formatKeyint(level, key_int)
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
			//if level == 2 {
			//	key = Substr(key, 0, 5)
			//	fmt.Println(" 字符串截取后的 ", key)
			//}

			// 去除文件名后缀
			file := strings.TrimSuffix(link, fileSuffix)
			fmt.Println("文件去除后缀：", file)
			if len(file) > 0 {
				key_int, _ := strconv.Atoi(key)
				key_int = formatKeyint(level, key_int)
				area[key_int] = value
				areaSort = append(areaSort, key_int)
				//
				data_url[key_int] = paths + link
				//collect_url[key] = paths + "/" + link
			}
		}

	})
	if is_villagetable {
		return nil
	}
	if len(data_url) < 1 {
		return nil
	}
	//fmt.Println(area)
	//fmt.Println("============")
	//fmt.Println("============")
	//fmt.Println("============")
	fmt.Println(len(area))
	return data_url
}

func formatKeyint(level, key_int int) int {
	if level == 2 {
		key_int = key_int / 1000000
	} else if level == 3 {
		key_int = key_int / 1000
	}
	return key_int
}

//格式化数据
func FormatArea() {
	//初始化
	areaFormat = make(map[int]map[int]string)
	areaFormat[topLevelId] = make(map[int]string)
	//
	for index := range areaSort {
		key := areaSort[index]
		value := area[key]

		fmt.Println("key:", key)
		fmt.Println("value:", value)
		//省份
		if key < 1000 {
			//是否初始化
			if areaFormat[topLevelId] == nil {
				areaFormat[topLevelId] = make(map[int]string)
			}
			maps := areaFormat[topLevelId]
			maps[key*10*1000] = value
			areaFormat[topLevelId] = maps
			//直辖市处理
			formatMunicipality(key * 10 * 1000)
		} else if key < 1000000 {
			// 省下的市
			fmt.Println("市县 key:", key)
			tmp := math.Floor(float64(key) / 10000)
			//获得 上级省
			city := int(tmp) * 10000
			if areaFormat[city] == nil {
				areaFormat[city] = make(map[int]string)
			}
			maps := areaFormat[city]
			maps[key] = value
			areaFormat[city] = maps

		} else {
			// 市下的 县
			fmt.Println("乡镇 key:", key)
			street := float64(key) / 1000
			street_tmp := math.Floor(street)
			value = strings.Replace(value, "居民委员会", "", -1)
			value = strings.Replace(value, "办事处", "", -1)
			//如果相等，表示该 是乡镇  不是村
			if street == street_tmp {
				tmp :=street_tmp
				//乡镇
				tmp = math.Floor(street_tmp / 1000)
				//if len(strconv.Itoa(int(street_tmp)))>6 {
				//	//乡镇
				//	tmp = math.Floor(street_tmp / 1000)
				//}
				//获得 上级市
				parent := int(tmp)
				if len(strconv.Itoa(parent))==3 {
					tmp = math.Floor(street_tmp / 10)*10
					parent = int(tmp)
				}
				if areaFormat[parent] == nil {
					areaFormat[parent] = make(map[int]string)
				}
				maps := areaFormat[parent]
				maps[int(street_tmp)] = value
				areaFormat[parent] = maps
			} else {
				//村
				fmt.Println("村 key:", key)
				tmp := math.Floor(float64(key) / 1000)
				//获得 上级县
				parent := int(tmp)
				if areaFormat[parent] == nil {
					areaFormat[parent] = make(map[int]string)
				}
				maps := areaFormat[parent]
				maps[key] = value
				areaFormat[parent] = maps
			}

		}
		//break
	}
	fmt.Println("area:", area)
	fmt.Println("areaSort:", areaSort)
	fmt.Println("areaFormat:", areaFormat)
	fmt.Println(len(areaFormat))
	str, err := json.Marshal(&areaFormat)
	if err != nil {
		fmt.Println("序列化失败:", err)
	}
	fmt.Println("JSON:", string(str))
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
	fmt.Println("formatMunicipality city:", city)
	fmt.Println("formatMunicipality province:", province)
}

func SaveFile() {
	str := util.ObjToJson(&areaFormat)
	//
	util.SaveFile(str, consts.JSON_FILE)
	str2 := string(str)
	byte2 := []byte("module.exports=" + str2)
	util.SaveFile(byte2, consts.JSON_FILE_VUE)
}

//扩展数据填充
func FormatExtData() {
	//
	mca_gov.ReadExt()
	fmt.Printf("AreaFormatExt:", mca_gov.AreaFormatExt)
	//
	tmpMap := make(map[int]string)
	for key, _ := range mca_gov.AreaFormatExt {
		if key > 86 {
			areaFormat[key] = mca_gov.AreaFormatExt[key]
		}
		if key == 86 {
			tmpMap = mca_gov.AreaFormatExt[key]
			maps := areaFormat[topLevelId]
			for k, v := range tmpMap {
				maps[k] = v
			}
			areaFormat[topLevelId] = maps
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

//截取字符串 start 起点下标 end 终点下标(不包括)
func Substr(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < 0 || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}
