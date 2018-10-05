package mca_gov

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

//中国 行政区划代码
//http://www.mca.gov.cn/article/sj/xzqh/
//http://www.mca.gov.cn/article/sj/xzqh/2018/
//额外数据
var AreaFormatExt map[int]map[int]string

//扩展数据原始文件
const extFile = "file/ext.csv"

func init() {

}

func ReadExt() {
	//b, err := ioutil.ReadFile("test.log")
	//if err != nil {
	//	fmt.Print(err)
	//}
	//fmt.Println(b)
	//str := string(b)

	//path := util.GetCurrentDirectory()
	//file := path + "/" + extFile
	//fmt.Println(file)
	//fmt.Println(util.PathExists(extFile))
	//
	AreaFormatExt = make(map[int]map[int]string)
	AreaFormatExt[86] = make(map[int]string)
	maps := AreaFormatExt[86]
	maps[810000] = "香港"
	maps[820000] = "澳门"
	maps[710000] = "台湾"
	AreaFormatExt[86] = maps
	//
	getCsv(extFile)
}
func getCsv(path string) {

	f, err := os.Open(path)
	CheckErr(err)
	defer f.Close()

	reader := csv.NewReader(f)
	record := []string{}
	for {
		record, err = reader.Read()
		if err == io.EOF {
			break
		} else if nil != err {
			fmt.Println(err)
			return
		}
		key, _ := strconv.Atoi(record[2])
		if key == 86 {
			continue
		}
		if key < 86 {
			continue
		}
		if ok := AreaFormatExt[key]; ok == nil {
			AreaFormatExt[key] = make(map[int]string)
		}
		maps := AreaFormatExt[key]
		if len(record) > 2 && record[1] != "" {
			str := record[1]
			key2, _ := strconv.Atoi(record[0])
			maps[key2] = str
			AreaFormatExt[key] = maps
		}

	}
}

func CheckErr(err error) {
	if nil != err {
		panic(err)
	}
}
