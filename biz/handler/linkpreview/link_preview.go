package linkpreview

import (
	"affine-worker-go/biz/common/httpclient"
	"affine-worker-go/biz/common/myutils"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func LinkPreview(ctx context.Context, c *app.RequestContext) {
	hlog.Info("received request: ", string(c.GetHeader("Origin")), string(c.Method()))
	body := make(map[string]interface{})
	data, err := c.Body()
	if err != nil {
		hlog.Error("invalid request body: ", err)
		c.JSON(consts.StatusBadRequest, utils.H{
			"msg": "Invalid request body",
		})
		return
	}
	err = json.Unmarshal(data, &body)
	if err != nil {
		hlog.Error("invalid request body: ", err)
		c.JSON(consts.StatusBadRequest, utils.H{
			"msg": "Invalid request body",
		})
		return
	}
	var targetURL *url.URL
	if _, exist := body["url"]; exist {
		if rawURL, ok := body["url"].(string); ok {
			if targetURL, err = myutils.FixURL(rawURL); err != nil {
				hlog.Error("invalid url")
				c.JSON(consts.StatusBadRequest, utils.H{
					"msg": "Invalid URL",
				})
				return
			}
		} else {
			hlog.Error("invalid url")
			c.JSON(consts.StatusBadRequest, utils.H{
				"msg": "Invalid URL",
			})
			return
		}
	} else {
		hlog.Error("invalid url")
		c.JSON(consts.StatusBadRequest, utils.H{
			"msg": "Invalid URL",
		})
		return
	}
	hlog.Info("processing request url: ", targetURL)
	req := &http.Request{
		Method: http.MethodGet,
		URL:    targetURL,
		// Header: myutils.CloneHeaders(&c.Request.Header),
	}
	resp, err := httpclient.Do(req)
	if err != nil {
		hlog.Error("link preview error: ", err)
		c.JSON(consts.StatusInternalServerError, utils.H{
			"msg": "Internal Server Error",
		})
		return
	}
	defer resp.Body.Close()
	respBody := map[string]interface{}{
		"url":      resp.Request.URL.String(),
		"images":   []string{},
		"videos":   []string{},
		"favicons": []string{},
	}
	if doc, err := html.Parse(resp.Body); err != nil {
		hlog.Error("link preview error: ", err)
		c.JSON(consts.StatusInternalServerError, utils.H{
			"msg": "Internal Server Error",
		})
		return
	} else {
		nodes := []*html.Node{doc}
		for {
			if len(nodes) == 0 {
				break
			}
			currentNode := nodes[0]
			nodes = nodes[1:]
			for chi := currentNode.FirstChild; chi != nil; chi = chi.NextSibling {
				nodes = append(nodes, chi)
			}
			if currentNode.Type != html.ElementNode {
				continue
			}
			if currentNode.DataAtom == atom.Meta {
				var property, content string
				for _, attr := range currentNode.Attr {
					if attr.Key == "property" || attr.Key == "name" {
						property = attr.Val
					}
					if attr.Key == "content" {
						content = attr.Val
					}
				}
				if property == "" || content == "" {
					continue
				}
				switch strings.ToLower(property) {
				case "og:title":
					respBody["title"] = content
				case "og:site_name":
					respBody["siteName"] = content
				case "og:description":
					respBody["description"] = content
				case "og:image":
					respBody["images"] = appendURL(content, respBody["images"].([]string))
				case "og:video":
					respBody["videos"] = appendURL(content, respBody["vides"].([]string))
				case "og:type":
					respBody["mediaType"] = content
				case "description":
					if _, exist := respBody["description"]; !exist {
						respBody["description"] = content
					}
				}
			}
			if currentNode.DataAtom == atom.Link {
				var rel, href string
				for _, attr := range currentNode.Attr {
					if attr.Key == "rel" {
						rel = attr.Val
					}
					if attr.Key == "href" {
						href = attr.Val
					}
				}
				if strings.Contains(rel, "icon") {
					respBody["favicons"] = appendURL(href, respBody["favicons"].([]string))
				}
			}
			if currentNode.DataAtom == atom.Title {
				if chi := currentNode.FirstChild; chi != nil && chi.Type == html.TextNode {
					if _, exist := respBody["title"]; !exist {
						respBody["title"] = chi.Data
					}
				}
			}
			if currentNode.DataAtom == atom.Img {
				var imgSrc string
				for _, attr := range currentNode.Attr {
					if attr.Key == "src" {
						imgSrc = attr.Val
						break
					}
				}
				if imgSrc != "" {
					respBody["images"] = appendURL(imgSrc, respBody["images"].([]string))
				}
			}
			if currentNode.DataAtom == atom.Video {
				var videoSrc string
				for _, attr := range currentNode.Attr {
					if attr.Key == "src" {
						videoSrc = attr.Val
						break
					}
				}
				if videoSrc != "" {
					respBody["vides"] = appendURL(videoSrc, respBody["videos"].([]string))
				}
			}
		}
	}
	if len(respBody["favicons"].([]string)) <= 0 {
		faviconUrl := targetURL.Scheme + "://" + targetURL.Host + "/" + "favicon.ico"
		respBody["favicons"] = appendURL(faviconUrl, respBody["favicons"].([]string))
	}
	c.JSON(consts.StatusOK, respBody)
}

func appendURL(rawURL string, urls []string) []string {
	if rawURL != "" {
		fixedURL, _ := myutils.FixURL(rawURL)
		if fixedURL != nil {
			return append(urls, fixedURL.String())
		}
	}
	return urls
}
