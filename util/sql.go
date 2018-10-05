package util

import (
	"io/ioutil"
	"strconv"
)

const SQL_FILE = "file/sql.sql"

var sqlArr = []string{}
//返回
func MakeSql(area map[int]map[int]string) string {
	for key, item := range area {
		process(key, item)
	}
	str := ""
	for _, item := range sqlArr {
		str += item
	}
	return str
}

func process(parent_id int, item map[int]string) {
	str := ""
	for key, value := range item {

		str = "INSERT INTO `area` (`id`,`name`,`parent_id`)VALUE ";
		str += "('" + strconv.Itoa(key) + "', '" + value + "', '" + strconv.Itoa(parent_id) + "');\n";
		sqlArr = append(sqlArr, str)
	}
}

func GetCreateTableSql() string {
	b, err := ioutil.ReadFile(SQL_FILE)
	checkErr(err)
	//fmt.Println(b)
	return string(b)
}
