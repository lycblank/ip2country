package ip2country

import "fmt"

var ipPool *IpPool

func Init(dirs ...string) {
	ipPool = &IpPool{}
	ipPool.Init(dirs...)
}

func SearchCountry(ip string) string {
	if ipPool != nil {
		if ipRange := ipPool.Search(ip); ipRange != nil {
			return ipRange.Country
		}
	}
	return ""
}



type FtpAddrInfo struct {
	Addr string
	FileName string
}

var addrInfos = [5]FtpAddrInfo{
	{Addr: "ftp.apnic.net:21", FileName: "/pub/stats/apnic/assigned-apnic-latest"},
	{Addr: "ftp.afrinic.net:21", FileName: "/pub/stats/afrinic/delegated-afrinic-latest"},
	{Addr: "ftp.lacnic.net:21", FileName: "/pub/stats/lacnic/delegated-lacnic-latest"},
	{Addr: "ftp.arin.net:21", FileName: "/pub/stats/arin/delegated-arin-extended-latest"},
	{Addr: "ftp.ripe.net:21", FileName: "/pub/stats/ripencc/delegated-ripencc-latest"},
}

func GetIpPoolData() ([]FtpFileData, error) {
	fileDatas := make([]FtpFileData, 0, len(addrInfos))
	for i,cnt:=0,len(addrInfos);i<cnt;i++{
		fileData, err := Download(addrInfos[i].Addr, addrInfos[i].FileName)
		if err != nil {
			fmt.Printf("addr:%s filename:%s err:%s",
				addrInfos[i].Addr, addrInfos[i].FileName, err.Error())
			return fileDatas, err
		}
		fileDatas = append(fileDatas, fileData)
	}
	return fileDatas, nil
}
