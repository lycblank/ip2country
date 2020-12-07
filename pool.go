package ip2country

import (
	"bufio"
	"bytes"
	"io"
	"math/big"
	"net"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
)

type IpRangData struct {
	Start *big.Int
	End *big.Int
	Country string
}

type IpPool struct {
	IpDatas []*IpRangData
}

func (p *IpPool) Init(dirs ...string) {
	rets := make([]*IpRangData, 0, 4096)
	if len(dirs) > 0 {
		for _, dir := range dirs {
			filepath.Walk(dir, func(filePath string, f os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if f.IsDir() {
					return nil
				}
				if path.Ext(f.Name()) == ".txt" {
					datas, _ := p.ParseFromFile(filePath)
					rets = append(rets, datas...)
				}
				return nil
			})
		}
	}
	if len(rets) <= 0 {
		ftpDatas, _ := GetIpPoolData()
		for i,cnt:=0,len(ftpDatas);i<cnt;i++{
			datas, _ := p.ParseFromBytes(ftpDatas[i].FileDatas)
			rets = append(rets, datas...)
		}
	}

	p.IpDatas = rets
	p.Sort()
}

func (p *IpPool) Sort() {
	sort.Slice(p.IpDatas, func(i,j int) bool {
		return p.IpDatas[i].Start.Cmp(p.IpDatas[j].Start) == -1
	})
}

func (p *IpPool) ParseFromBytes(datas []byte) ([]*IpRangData, error) {
	buf := bytes.NewBuffer(datas)
	return p.ParseFromReader(buf)
}

func (p *IpPool) ParseFromReader(reader io.Reader) ([]*IpRangData, error) {
	seq := []byte("|")
	bufReader := bufio.NewReader(reader)
	rets := make([]*IpRangData, 0, 4096)
	for line, _, err := bufReader.ReadLine(); err == nil; line, _, err = bufReader.ReadLine() {
		cells := bytes.Split(line, seq)
		if len(cells) < 7 {
			continue
		}
		ipType := string(cells[2])
		if ipType != "ipv4" && ipType != "ipv6" {
			continue
		}
		status := string(cells[6])
		if status != "allocated" && status != "assigned" {
			continue
		}
		data := &IpRangData{
			Country: string(cells[1]),
			Start:big.NewInt(0),
			End: big.NewInt(0),
		}
		startIp := net.ParseIP(string(cells[3]))
		data.Start.SetBytes(startIp.To16())

		if ipType == "ipv4" {
			num, _ := strconv.ParseInt(string(cells[4]), 10, 64)
			data.End.Add(data.Start, big.NewInt(num-1))
		} else { // ipv6
			num, _ := strconv.Atoi(string(cells[4]))
			data.End.Set(data.Start)
			for i := 0; i < num; i++ {
				data.End.SetBit(data.End, i, 1)
			}
		}
		rets = append(rets, data)
	}
	return nil, nil
}

func (p *IpPool) ParseFromFile(fileName string) ([]*IpRangData, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return p.ParseFromReader(f)
}

func (p *IpPool) Search(ip string) *IpRangData {
	ipValue := big.NewInt(0).SetBytes(net.ParseIP(ip).To16())
	idx := sort.Search(len(p.IpDatas), func(i int) bool {
		return p.IpDatas[i].Start.Cmp(ipValue) >= 0
	})

	if idx < len(p.IpDatas) {
		if p.IpDatas[idx].Start.Cmp(ipValue) <= 0 && p.IpDatas[idx].End.Cmp(ipValue) >= 0 {
			tmp := &IpRangData{
				Start: big.NewInt(0).Set(p.IpDatas[idx].Start),
				End: big.NewInt(0).Set(p.IpDatas[idx].End),
				Country:p.IpDatas[idx].Country,
			}
			return tmp
		}
	}
	return nil
}

