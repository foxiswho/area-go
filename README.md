# area-go
国家省市区获取

## 中国行政区划代码 (这里选择是)
http://www.mca.gov.cn/article/sj/xzqh/

http://www.mca.gov.cn/article/sj/xzqh/2018/

## 统计用区划和城乡划分代码
http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/

# 源码编译
```go
go get -u github.com/foxiswho/area-go
```
执行
```shell
cd $GOPATH/src/github.com/foxiswho/area-go/
go run area.go
```

生成后文件：`$GOPATH/src/github.com/foxiswho/area-go/area.js`

输出结果
```SEHLL
=======获取数据======

=======获取成功=====
=======格式化数据=====

=======格式化 扩展数据=====
=======处理成功=====
===================
=======  写入到文件 area.js    ============
=======  路径: $GOPATH/src/github.com/foxiswho/area-go/area.js    ============
=======  写入到文件 area.sql    ============
=======  路径: $GOPATH/src/github.com/foxiswho/area-go/area.sql    ============
===================
=======数据保存成功=======

```


# 编译成多环境执行文件
## linux
```SHELL
GOOS=linux GOARCH=amd64 go build -o area-linux area.go
```
## MAC
```SHELL
GOOS=darwin GOARCH=amd64 go build -o area-mac area.go
```

## Windows
```SHELL
GOOS=windows GOARCH=amd64 go build -o area-win area.go
```
