package key

import (
	//"encoding/json"
	//"io/ioutil"
	//"os"
	"errors"
	"fmt"
	"regexp"
	//"reflect"
	//"strconv"
	//"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	//"github.com/syndtr/goleveldb/leveldb"
)

//var ldb *leveldb.DB

func SetupRoutes(r *gin.Engine) {

	r.GET("/prefs/:context/:key", GetKey)
	r.GET("/prefs/:context/:app/:name/:key", GetKey)

	//r.GET("/prefs/:context/:key", GetKey)
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
	test_results := map[string]string{
		"dev.someapp.foo":        "bar",
		"dev.someapp.123456.foo": "bar2",
	}

	search := viper.GetStringSlice("search")

	var return_value string
	var return_key string

	for _, search_key := range search {
		fmt.Println("Search Key:", search_key)

		trans_key, err := TranslateKey(search_key, &c.Params)

		if err != nil || trans_key == "" {
			continue
		}

		if val, ok := test_results[trans_key]; ok {
			return_value = val
			return_key = trans_key
			fmt.Println("Return Key:", return_key)
			fmt.Println("Return Value:", return_value)
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
