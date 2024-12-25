package linkpreview

import (
	"context"
	"encoding/json"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func LinkPreview(ctx context.Context, c *app.RequestContext) {
	hlog.Info("received request: ", string(c.GetHeader("Origin")), string(c.Method()))
	body := make(map[string]interface{})
	data, err := c.Body()
	if err != nil {
		hlog.Error("invalid request body")
		c.JSON(consts.StatusBadRequest, utils.H{
			"msg": "Invalid request body",
		})
		return
	}
	err = json.Unmarshal(data, &body)
	if err != nil {
		hlog.Error("invalid request body")
		c.JSON(consts.StatusBadRequest, utils.H{
			"msg": "Invalid request body",
		})
		return
	}
	var targetURL string
	if _, exist := body["url"]; exist {
		if rawURL, ok := body["url"].(string); ok {
			targetURL = rawURL
		} else {
			hlog.Error("invalid url")
			c.JSON(consts.StatusBadRequest, utils.H{
				"msg": "Invalid URL",
			})
		}
	} else {
		hlog.Error("invalid url")
		c.JSON(consts.StatusBadRequest, utils.H{
			"msg": "Invalid URL",
		})
	}
	hlog.Info("processing request url: ", targetURL)
	// req, err := http.NewRequest()
}
