package main

import (
	"flag"
	"fmt"
	"github.com/lycblank/ip2country"
	"io/ioutil"
	"path"
)

func main() {
	var dir string
	flag.StringVar(&dir, "dir", "data", "download dir")
	flag.Parse()
	datas, err := ip2country.GetIpPoolData()
	if err != nil {
		fmt.Printf("download failed err:%s\n", err.Error())
		return
	}
	for i,cnt:=0,len(datas);i<cnt;i++{
		fname := fmt.Sprintf("delegated-ip-%d.txt", i+1)
		err := ioutil.WriteFile(path.Join(dir, fname), datas[i].FileDatas, 0666)
		if err != nil {
			fmt.Printf("write file failed. fname:%s err:%s\n", fname, err.Error())
			return
		}
	}
	fmt.Println("download file success")
}

