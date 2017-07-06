package handle

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"net/http"
	"github.com/robbert229/jwt"
	"github.com/LiSheep/allpass/model"
	"gopkg.in/mgo.v2"
)

func validateToken(token string) (bool) {
	if token == "" {
		return false
	}
	err := algorithm.Validate(token)
	if err == nil {
		return true
	} else {
		return false
	}
}

func authMiddleware(c *gin.Context) {
	fmt.Println("in middleware")
	token := c.Request.Header.Get("Authentication")
	if token != "" {
		if validateToken(token) {
			fmt.Println("valid token...")
			claims, err := algorithm.Decode(token)
			if err != nil {
				c.Abort()
				c.Redirect(http.StatusOK, "/user/login")
				fmt.Println("invalid claims")
				return
			}
			name, err := claims.Get("name")
			if err != nil {
				panic(err)
			}
			user, err := model.UserFind(name)
			if err != nil {
				c.Abort()
				c.Redirect(http.StatusOK, "/user/login")
				return
			}
			fmt.Println("login success", user.Name)
			c.Set("user", user)
			c.Next()
			return
		}
	}
	c.Redirect(http.StatusFound, "/user/login")
	c.Abort()
}

func checkLogin(c *gin.Context) {
	fmt.Println("checklogin")
	token, err := c.Cookie("token")
	if err == nil && token != "" {
		if validateToken(token) {
			c.Redirect(http.StatusFound, "/password")
			c.Abort()
		} else {
			c.Next()
		}
	}
}

func encodeUserPassword(password string) []byte {
	return cret.EncodeHmacSha512([]byte(password))
}

var algorithm jwt.Algorithm
func initUserHandler(router *gin.Engine) {
	algorithm = jwt.HmacSha256("ThisIsTheSecret")
	router.POST("/user/updatePass", userUpdatePass)
	checkAuth := router.Group("/")
	checkAuth.Use(checkLogin)
	{
		checkAuth.GET("/user/login", userLogin)
		checkAuth.POST("/user/login", userLoginAPI)
		checkAuth.GET("/user/register", userRegister)
		checkAuth.POST("/user/register", userRegisterAPI)
	}
	router.Use(authMiddleware)
	router.GET("/user", userMain)
	router.GET("/user/logout", userLogout)
}

func userMain(c *gin.Context) {
	c.String(http.StatusOK, "hello world")
}

func userLogin(c *gin.Context) {
	responseHTML(c, http.StatusOK, "./views/user/login.html")
}

func userLogout(c *gin.Context) {
	c.Redirect(http.StatusFound, "/user/login")
}

func userLoginAPI(c *gin.Context) {
	name := c.PostForm("name")
	password := c.PostForm("password")
	res, err := model.UserValidate(name, encodeUserPassword(password))
	if err != nil || !res  {
		c.JSON(http.StatusOK, getError(false, "server error"))
		return
	}
	claims := jwt.NewClaim()
	claims.Set("name", name)
	token, err := algorithm.Encode(claims)
	if err != nil {
		c.JSON(http.StatusOK, getError(false, "server error"))
		return
	}
	//setTokenCookie(c, token)
	c.JSON(http.StatusOK, getError(true, token))
}

func userRegister(c *gin.Context) {
	responseHTML(c, http.StatusOK, "./views/user/register.html")
}

func userRegisterAPI(c *gin.Context) {
	name := c.PostForm("name")
	pass := c.PostForm("password")
	_, err := model.UserFind(name)
	if err!= nil && err != mgo.ErrNotFound {
		c.JSON(http.StatusOK, getError(false, err.Error()))
		return
	}
	if err == nil {
		c.JSON(http.StatusOK, getError(false, "user exist"))
		return
	}

	if len(name) < 3 {
		c.JSON(http.StatusOK, getError(false, "user name too short"))
		return
	}
	if len(name) > 10 {
		c.JSON(http.StatusOK, getError(false, "user name too long"))
		return
	}
 	if len(pass) < 6 {
		c.JSON(http.StatusOK, getError(false, "password too short"))
		return
	}
	if (len(pass) > 30) {
		c.JSON(http.StatusOK, getError(false, "password too long"))
		return
	}
	err = model.UserAdd(name, encodeUserPassword(pass))
	if err != nil {
		c.JSON(http.StatusOK, getError(false, err.Error()))
		return
	}
	c.JSON(http.StatusOK, getError(true, ""))
}

func userUpdatePass(c *gin.Context) {
	name := c.PostForm("name")
	password := c.PostForm("password")
	res, err := model.UserValidate(name, encodeUserPassword(password))
	if err != nil || !res  {
		c.JSON(http.StatusOK, getError(false, "error"))
		return
	}
	newPassword := c.PostForm("newPassword")
	if len(newPassword) < 6 {
		c.JSON(http.StatusOK, getError(false, "password too short"))
		return
	}
	if (len(newPassword) > 30) {
		c.JSON(http.StatusOK, getError(false, "password too long"))
		return
	}
	userInterface, ex := c.Get("user")
	if !ex {
		fmt.Println("get user fail")
	}
	err = model.UserPassUpdate(name, encodeUserPassword(newPassword))
	if err != nil {
		c.JSON(http.StatusOK, getError(false, "server error"))
		return
	}
	user := userInterface.(*model.User)
	datas := c.PostForm("datas")
	err = updateAllPassword(user, datas)
	if err != nil {
		c.JSON(http.StatusOK, getError(false, err.Error()))
		return
	}

	c.JSON(http.StatusOK, getError(true, ""))
}