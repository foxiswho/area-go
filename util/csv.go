package util

import "strconv"

var csvArr = []string{}

func MakeCsv(area map[int]map[int]string) string {
	csvArr = append(csvArr, "id\t名称\t上级id\n")
	for key, item := range area {
		processCsv(key, item)
	}
	str := ""
	for _, item := range csvArr {
		str += item
	}
	return str
}

func processCsv(parent_id int, item map[int]string) {
	str := ""
	for key, value := range item {

		str = strconv.Itoa(key) + "\t" + value + "\t" + strconv.Itoa(parent_id) + "\n";
		csvArr = append(csvArr, str)
	}
}
