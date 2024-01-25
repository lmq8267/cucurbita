package web

import (
	"net/http"
	"strconv"
	"time"

	"github.com/foolin/goview"
	"github.com/gin-gonic/gin"
	"github.com/lanthora/cucurbita/candy"
	"github.com/lanthora/cucurbita/storage"
)

func Index(c *gin.Context) {
	online := int64(0)
	daily := int64(0)
	weekly := int64(0)
	all := int64(0)
	domain := int64(0)
	user := int64(0)

	currentUser := c.MustGet("user").(*User)
	if currentUser.Role == "admin" {
		storage.Model(&candy.Device{}).Where("online = true").Count(&online)
		storage.Model(&candy.Device{}).Where("online = true").Or("conn_updated_at > ?", time.Now().AddDate(0, 0, -1)).Count(&daily)
		storage.Model(&candy.Device{}).Where("online = true").Or("conn_updated_at > ?", time.Now().AddDate(0, 0, -7)).Count(&weekly)
		storage.Model(&candy.Device{}).Count(&all)
		storage.Model(&candy.Domain{}).Count(&domain)
		storage.Model(&User{}).Count(&user)
	} else {
		storage.Model(&candy.Device{}).Where("online = true AND username = ?", currentUser.Name).Count(&online)
		storage.Model(&candy.Device{}).Where("online = true AND username = ?", currentUser.Name).Or("conn_updated_at > ? AND username = ?", time.Now().AddDate(0, 0, -1), currentUser.Name).Count(&daily)
		storage.Model(&candy.Device{}).Where("online = true AND username = ?", currentUser.Name).Or("conn_updated_at > ? AND username = ?", time.Now().AddDate(0, 0, -7), currentUser.Name).Count(&weekly)
		storage.Model(&candy.Device{}).Where("username = ?", currentUser.Name).Count(&all)
		storage.Model(&candy.Domain{}).Where("username = ?", currentUser.Name).Count(&domain)
		storage.Model(&User{}).Where("name = ?", currentUser.Name).Or("inviter = ?", currentUser.Name).Count(&user)
	}

	type Card struct {
		Href  string
		Title string
		Value string
	}
	cards := []Card{}
	cards = append(cards, Card{Href: "/device?active=online", Title: "当前在线设备", Value: strconv.FormatInt(online, 10)})
	cards = append(cards, Card{Href: "/device?active=daily", Title: "每日活跃设备", Value: strconv.FormatInt(daily, 10)})
	cards = append(cards, Card{Href: "/device?active=weekly", Title: "每周活跃设备", Value: strconv.FormatInt(weekly, 10)})
	cards = append(cards, Card{Href: "/device", Title: "全部设备", Value: strconv.FormatInt(all, 10)})
	cards = append(cards, Card{Href: "/domain", Title: "网络", Value: strconv.FormatInt(domain, 10)})
	cards = append(cards, Card{Href: "/user", Title: "用户", Value: strconv.FormatInt(user, 10)})
	c.HTML(http.StatusOK, "index.html", goview.M{
		"cards": cards,
	})
}

func Favicon(c *gin.Context) {
	buffer, err := views.ReadFile("views/favicon.ico")
	if err != nil {
		c.Status(http.StatusNotFound)
	} else {
		c.Data(http.StatusOK, "image/x-icon", buffer)
	}
}
