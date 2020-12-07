package ip2country

import (
	"github.com/jlaffaye/ftp"
	"io/ioutil"
	"time"
)

type FtpFileData struct {
	FileDatas []byte
}

func Download(addr string, filename string) (fileData FtpFileData, err error){
	c, err := ftp.Dial(addr, ftp.DialWithTimeout(10*time.Second))
	if err != nil {
		return fileData, err
	}
	defer c.Quit()
	if err := c.Login("anonymous", "anonymous"); err != nil {
		return fileData, err
	}
	r, err := c.Retr(filename)
	if err != nil {
		return fileData, err
	}
	defer r.Close()
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return fileData, err
	}
	fileData.FileDatas = buf
	return fileData, nil
}