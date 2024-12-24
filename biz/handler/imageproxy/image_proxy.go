package imageproxy

import (
	"affine-worker-go/biz/common/headerutil"
	"affine-worker-go/biz/common/myutils"
	"context"
	"errors"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func ImageProxy(ctx context.Context, c *app.RequestContext) {
	rawURL := c.Query("url")
	if rawURL == "" {
		hlog.Warn(errors.New(`missing "url" parameter`))
		c.JSON(consts.StatusBadRequest, utils.H{
			"msg": `Missing "url" parameter`,
		})
		return
	}
	targetURL, err := myutils.FixURL(rawURL)
	if err != nil || targetURL == nil {
		hlog.Error(err)
		c.JSON(consts.StatusBadRequest, utils.H{
			"msg": "Invalid URL",
		})
		return
	}
	headers, err := headerutil.CloneHeaders(&c.Request.Header)
	if err != nil {
		hlog.Error("invalid headers:", err)
		c.JSON(consts.StatusBadRequest, utils.H{
			"msg": "Failed to fetch image",
		})
		return
	}
	c.JSON(consts.StatusOK, headers)
}
