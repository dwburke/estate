package key

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"

	"github.com/dwburke/prefs/storage"
	"github.com/dwburke/prefs/storage/memory"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

//var ldb *leveldb.DB

var (
	ErrNotEnoughValues = errors.New("prefs: not enough values")
	//ErrAllValuesNotUsed = errors.New("prefs: not found")
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/prefs/:context/:key", GetKey)
	r.POST("/prefs/:context/:key", SetKey)
	r.DELETE("/prefs/:context/:key", DeleteKey)
}

// take key template and try to replace everything with param's passed to api
// TODO - optionally require matching all params, return appro error
func TranslateKey(template string, p *gin.Params) (string, error) {
	re1, err := regexp.Compile(`{(.*?)}`)
	if err != nil {
		return "", err
	}

	result_slice := re1.FindAllStringSubmatch(template, -1)

	for _, str := range result_slice {
		needed_value, found := p.Get(str[1])

		if found == false {
			return "", ErrNotEnoughValues
		}

		var re2 = regexp.MustCompile(str[0])
		template = re2.ReplaceAllString(template, needed_value)
	}

	return template, nil
}

// Matching all params is deliberately not required, we want "best match"
func GetKey(c *gin.Context) {
	search := viper.GetStringSlice("prefs.search")

	var st storage.Storage
	st = memory.New()
	defer st.Close()

	var return_value string
	var return_key string

	for _, search_key := range search {
		trans_key, err := TranslateKey(search_key, &c.Params)

		if err == ErrNotEnoughValues {
			continue
		}

		if err != nil {
			c.JSON(500, gin.H{
				"error": err,
			})
			return
		}

		val, err := st.Get(trans_key)

		// no rows, try next key
		if err == storage.ErrNotFound {
			continue
		}

		// unexpected error, fail
		if err != nil {
			c.JSON(500, gin.H{
				"error": err,
			})
			return
		}

		// value? found it!
		if val != "" {
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
		return
	} else {
		c.JSON(404, gin.H{
			"error": "Not Found",
		})
		return
	}
}

// TODO - matching all params is required
func SetKey(c *gin.Context) {
	type E struct {
		Data string
	}

	post_data := &E{}
	c.Bind(post_data)
	fmt.Println("post_data:", post_data)

	fmt.Println("body:", reflect.TypeOf(c.Request.Body))
	param_interface, ok := c.Get("value")
	fmt.Println("ok:", ok)
	fmt.Println(param_interface)

	var param_value string
	param_value = param_interface.(string)

	if param_value == "" {
		c.JSON(500, gin.H{
			"error": "'value' required",
		})
		return
	}

	var return_key string

	search := viper.GetStringSlice("prefs.search")

	// find matching key to set based on passed in variables
	for _, search_key := range search {
		trans_key, err := TranslateKey(search_key, &c.Params)

		// TODO - be more specific about what errors are ok here
		if err != nil || trans_key == "" {
			continue
		}

		return_key = trans_key
		break
	}

	if return_key == "" {
		c.JSON(404, gin.H{
			"error": "key not found",
		})
		return
	}

	// storage object
	var st storage.Storage
	st = memory.New()
	defer st.Close()

	err := st.Set(return_key, param_value)
	if err != nil {
		c.JSON(500, gin.H{
			"error": fmt.Sprintf("Error setting value: %s", err),
		})
		return
	}

	c.JSON(200, gin.H{
		"key":   return_key,
		"value": param_value,
	})
}

// TODO - matching all params is required
func DeleteKey(c *gin.Context) {
	search := viper.GetStringSlice("prefs.search")

	var st storage.Storage
	st = memory.New()
	defer st.Close()

	var return_key string

	for _, search_key := range search {
		trans_key, err := TranslateKey(search_key, &c.Params)

		// TODO - be more specific about what errors are ok here
		if err != nil || trans_key == "" {
			continue
		}

		return_key = trans_key

		err = st.Delete(trans_key)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err,
			})
			return
		}
	}

	if return_key != "" {
		c.JSON(200, gin.H{
			"key": return_key,
		})
		return
	} else {
		c.JSON(404, gin.H{
			"error": "Not Found",
		})
		return
	}
}
