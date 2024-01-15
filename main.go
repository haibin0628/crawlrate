package main

import (
	"crawlrate/models"
	"net/http"

	"github.com/jasonlvhit/gocron"

	"github.com/haibin0628/galaxylib"
	"github.com/labstack/echo"
)

func main() {

	galaxylib.DefaultGalaxyConfig.InitConfig()
	galaxylib.DefaultGalaxyLog.ConfigLogger()

	crawlData()

	e := echo.New()

	e.GET("/currency/:convert", func(c echo.Context) error {
		convertor := c.Param("convert")
		currency := &models.Currency{
			Convertor: convertor,
		}
		ret := currency.Get()

		data := &struct {
			Convert string
			Rate    float64
			Date    string
		}{
			Convert: ret.Convertor,
			Rate:    ret.Rate,
			Date:    ret.CrawlTime,
		}

		return c.JSON(http.StatusOK, data)
	})

	e.Start("0.0.0.0:1299")

}

func crawlData() {

	crawlTime := galaxylib.GalaxyCfgFile.MustValue("data", "crawlTime")

	go func() {

		gocron.Every(1).Day().At(crawlTime).Do(func() {
			c := &models.Currency{}

			c.FromRemote()
		})

		<-gocron.Start()
	}()
}
