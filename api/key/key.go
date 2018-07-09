package key

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/dwburke/prefs/storage"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

//var ldb *leveldb.DB

func SetupRoutes(r *gin.Engine) {

	r.GET("/prefs/:context/:key", GetKey)
	//r.POST("/prefs/:context/:key", SetKey)
}

func TranslateKey(template string, p *gin.Params) (string, error) {
	re1, err := regexp.Compile(`{(.*?)}`)
	if err != nil {
		return "", err
	}

	result_slice := re1.FindAllStringSubmatch(template, -1)

	for _, str := range result_slice {
		needed_value, found := p.Get(str[1])

		if found == false {
			return "", errors.New(fmt.Sprintf("Value for %s not found", str[0]))
		}

		var re2 = regexp.MustCompile(str[0])
		template = re2.ReplaceAllString(template, needed_value)
	}

	return template, nil
}

func GetKey(c *gin.Context) {
	search := viper.GetStringSlice("prefs.search")

	st, err := storage.New()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
	}

	var return_value string
	var return_key string

	for _, search_key := range search {
		trans_key, err := TranslateKey(search_key, &c.Params)

		if err != nil || trans_key == "" {
			continue
		}

		val, err := st.Get(trans_key)

		if err == nil && val != "" {
			return_value = val
			return_key = trans_key
			break
		}
	}

	if return_key != "" && return_value != "" {
		c.JSON(200, gin.H{
			"key":   return_key,
			"value": return_value,
		})
	} else {
		c.JSON(404, gin.H{
			"error": "Not Found",
		})
	}
}

//func SetKey(c *gin.Context) {
//param_context := c.Param("context")
//param_key := c.Param("key")

//if ldb == nil {
//var err error
//ldb, err = leveldb.OpenFile("data", nil)
//if err != nil {
//c.JSON(500, gin.H{"error": err})
//return
//}
//defer ldb.Close()
//}

//c.JSON(200, gin.H{
//"context": param_context,
//"key":     param_key,
//"value":   "bar",
//})
//}
