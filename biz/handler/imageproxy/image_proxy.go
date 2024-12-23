package imageproxy

import (
	"affine-worker-go/biz/common/myutils"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func ImageProxy(ctx context.Context, c *app.RequestContext) {
	rawURL := c.Query("url")
	if rawURL == "" {
		c.JSON(consts.StatusBadRequest, utils.H{
			"msg": `Missing "url" parameter`,
		})
		return
	}
	targetURL, err := myutils.FixURL(rawURL)
	if err != nil || targetURL == nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"msg": "Invalid URL",
		})
		return
	}
}
