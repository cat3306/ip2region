package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lionsoul2014/ip2region/v1.0/binding/golang/ip2region"
)

var (
	ipDbPath  = "./ip2region.db"
	ipDbFile1 = "https://fastly.jsdelivr.net/gh/lionsoul2014/ip2region@master/v1.0/data/ip2region.db"
	ipDbFile2 = "https://fastly.jsdelivr.net/gh/bqf9979/ip2region@master/data/ip2region.db"
	ip2Region *ip2region.Ip2Region
)

func main() {
	err := initIpDb()
	if err != nil {
		panic(err)
	}
	engine := gin.Default()
	engine.GET("/ip2regin", Ip2Regin)
	engine.GET("/download/db", DownloadDb)
	engine.POST("/batch/ip2regin", BatchIp2Regin)
	err = engine.Run("0.0.0.0:7878")
	log.Println(err)
}

func BatchIp2Regin(c *gin.Context) {
	ips := make([]string, 0)
	err := c.BindJSON(&ips)
	if err != nil {
		RspError(c, err.Error())
		return
	}
	list := make([]*IpInfo, 0)
	for _, v := range ips {
		info, err := ip2Regin(&Ip2ReginReq{
			Ip: v,
		})
		if err != nil {
			RspError(c, err.Error())
			return
		}
		list = append(list, info)
	}
	RspOk(c, list)
}
func initIpDb() error {
	region, err := ip2region.New(ipDbPath)
	if err != nil {
		return err
	}
	ip2Region = region
	return nil
}
func DownloadDb(c *gin.Context) {
	n, err := downloadFile(ipDbPath, ipDbFile2)
	if err != nil {
		log.Printf("DownloadDb err:%s", err.Error())
		RspError(c, err.Error())
		return
	}
	RspOk(c, n)
}
func downloadFile(filepath string, url string) (int64, error) {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	out, err := os.Create(filepath)
	if err != nil {
		return 0, err
	}
	defer out.Close()

	// Write the body to file
	n, err := io.Copy(out, resp.Body)
	return n, err
}

type Ip2ReginReq struct {
	Ip string `json:"ip" form:"ip"`
}

func Ip2Regin(c *gin.Context) {
	req := &Ip2ReginReq{}
	err := c.BindQuery(req)
	if err != nil {
		log.Printf("err:%s", err.Error())
		RspError(c, err.Error())
		return
	}
	v, err := ip2Regin(req)
	if err != nil {
		log.Printf("err:%s", err.Error())
		RspError(c, err.Error())
		return
	}
	RspOk(c, v)
}

type IpInfo struct {
	CityId   int64  `json:"city_id"`
	Country  string `json:"country"`
	Region   string `json:"region"`
	Province string `json:"province"`
	City     string `json:"city"`
	ISP      string `json:"isp"`
	Ip       string `json:"ip"`
}

func ip2Regin(req *Ip2ReginReq) (*IpInfo, error) {
	if req.Ip == "" {
		return nil, errors.New("invalid ip")
	}
	info, err := ip2Region.MemorySearch(req.Ip)
	if err != nil {
		return nil, err
	}

	return &IpInfo{
		CityId:   info.CityId,
		Country:  info.Country,
		Region:   info.Region,
		Province: info.Province,
		City:     info.City,
		ISP:      info.ISP,
		Ip:       req.Ip,
	}, err
}
