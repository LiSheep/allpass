package handle

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/LiSheep/allpass/model"
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"encoding/json"
	"errors"
)


func initPasswordHandler(router *gin.Engine) {
	router.GET("/password/listAPI", listPasswordAPI)
	router.POST("/password/add", addPasswordAPI)
	router.POST("/password/remove", removePasswordAPI)
	router.POST("/password/update", updatePasswordAPI)
	router.POST("/password/updateAll", updateAllPasswordAPI)
}


func listPasswordAPI(c *gin.Context) {
	userInterface, ex := c.Get("user")
	if !ex {
		panic("get user fail")
	}
	user := userInterface.(*model.User)
	c.JSON(http.StatusOK, user.Passwords)
}


func addPasswordAPI(c *gin.Context) {
	site := c.PostForm("site")
	pass := c.PostForm("password")
	username := c.PostForm("username")
	userInterface, ex := c.Get("user")
	if !ex {
		panic("get user fail")
	}
	user := userInterface.(*model.User)
	if user.FindPassword(site) != nil {
		c.JSON(http.StatusOK, getError(false, "pass desc exist"))
		return
	}
	err := user.AddPassword(site, username, pass)
	if  err != nil {
		c.JSON(http.StatusOK, getError(false, ""));
		return
	}
	c.JSON(http.StatusOK, getError(true, ""))
}

func removePasswordAPI(c *gin.Context) {
	site := c.PostForm("site")
	userInterface, ex := c.Get("user")
	if !ex {
		panic("get user fail")
	}
	user := userInterface.(*model.User)
	err := user.RemovePassword(site)
	if err != nil {
		c.JSON(http.StatusOK, getError(false, ""));
		return
	}
	c.JSON(http.StatusOK, getError(true, ""))
}

func updatePasswordAPI(c *gin.Context) {
	id := c.PostForm("id")
	rid := bson.ObjectIdHex(id)
	site := c.PostForm("site")
	sec := c.PostForm("secret")
	username := c.PostForm("username")
	fmt.Println(rid, site, sec)
	userInterface, ex := c.Get("user")
	if !ex {
		panic("get user fail")
	}
	user := userInterface.(*model.User)
	err := user.UpdatePassword(rid, site, username, sec)
	if err != nil {
		c.JSON(http.StatusOK, getError(false, ""));
		return
	}
	c.JSON(http.StatusOK, getError(true, ""))
}

func updateAllPassword(user *model.User, datas string) error {
	var newPass model.Passwords
	json.Unmarshal([]byte(datas), &newPass)
	if len(user.Passwords) != len(newPass) {
		return errors.New("password not match")
	}
	for i, p := range(user.Passwords) {
		for _, np := range(newPass) {
			if p.Id == np.Id {
				user.Passwords[i].OldSecret = p.Secret
				user.Passwords[i].Secret = np.Secret
				break
			}
		}
	}
	return user.UpdateAllPasswords(user.Passwords)
}

func updateAllPasswordAPI(c *gin.Context) {
	datas, _ := c.GetPostForm("data")
	userInterface, ex := c.Get("user")
	if !ex {
		panic("get user fail")
	}
	user := userInterface.(*model.User)

	err := updateAllPassword(user, datas)
	if err != nil {
		c.JSON(http.StatusOK, getError(false, err.Error()))
	} else {
		c.JSON(http.StatusOK, getError(true, ""))
	}
}
