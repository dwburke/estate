package key

import (
	"errors"
	"fmt"
	//"reflect"
	"regexp"

	"github.com/dwburke/estate/storage"
	"github.com/dwburke/estate/storage/common"

	"github.com/dwburke/go-tools"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	ErrNotEnoughValues = errors.New("estate: not enough values")
	//ErrAllValuesNotUsed = errors.New("estate: not found")
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/estate/:context/:key", GetKey)
	r.POST("/estate/:context/:key", SetKey)
	r.DELETE("/estate/:context/:key", DeleteKey)
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

// Matching all params is deliberately not required, we want "best/first match"
func GetKey(c *gin.Context) {
	search := viper.GetStringSlice("estate.search")

	st, err := storage.New()
	if err != nil {
		c.JSON(500, gin.H{
			"error": "error opening database",
		})
		return
	}
	defer st.Close()

	var return_value string
	var return_key string

	params := tools.AllGinParams(c)

	for _, search_key := range search {
		trans_key, err := TranslateKey(search_key, &params)

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
		if err == common.ErrNotFound {
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
		if val != nil {
			return_value = string(val)
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

	params := tools.AllGinParams(c)

	param_value, ok := params.Get("value")

	if !ok {
		c.JSON(500, gin.H{
			"error": "'value' required",
		})
		return
	}

	var return_key string

	search := viper.GetStringSlice("estate.search")

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
	st, err := storage.New()
	if err != nil {
		c.JSON(500, gin.H{
			"error": "error opening database",
		})
		return
	}
	defer st.Close()

	err = st.Set(return_key, []byte(param_value))
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
	search := viper.GetStringSlice("estate.search")

	st, err := storage.New()
	if err != nil {
		c.JSON(500, gin.H{
			"error": "error opening database",
		})
		return
	}
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
