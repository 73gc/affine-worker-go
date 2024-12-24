package imageproxy

import (
	"affine-worker-go/biz/common/headerutil"
	"affine-worker-go/biz/common/httpclient"
	"affine-worker-go/biz/common/myutils"
	"context"
	"io"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
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
	req, err := http.NewRequest(string(c.Request.Method()), targetURL.String(), nil)
	if err != nil {
		hlog.Error("image proxy error: ", err)
		c.JSON(consts.StatusInternalServerError, utils.H{
			"msg": "Failed to fetch image",
		})
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, err := httpclient.Do(req)
	if err != nil {
		hlog.Error("image proxy error: ", err)
		c.JSON(consts.StatusInternalServerError, utils.H{
			"msg": "Failed to fetch image",
		})
	}
	defer resp.Body.Close()
	if contentType := resp.Header.Get("Content-Type"); contentType != "" {
		c.SetContentType(contentType)
	}
	if contentDisposition := resp.Header.Get("Content-Dispositon"); contentDisposition != "" {
		c.Response.Header.Set("Content-Disposition", contentDisposition)
	}

	if _, err := io.Copy(c.Response.BodyWriter(), resp.Body); err != nil {
		hlog.Error("image proxy error: ", err)
		c.JSON(consts.StatusInternalServerError, utils.H{
			"msg": "Failed to fetch image",
		})
	}
}
