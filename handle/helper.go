package handle

import (
	"io/ioutil"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ErrorJSON struct {
	Success bool  `json:"success"`
	Info string  `json:"info"`
}

func responseHTML(c *gin.Context, code int, path string) {
	data, e := ioutil.ReadFile(path)
	if e != nil {
		panic(e)
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", data)
}

func getError(success bool, info string) (ErrorJSON) {
	return ErrorJSON{
		Success: success,
		Info: info,
	}
}

func setTokenCookie(c *gin.Context, token string) {
	c.SetCookie("token", token, 102400, "/", "", false, false)
}

func clearTokenCookie(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, false)
}