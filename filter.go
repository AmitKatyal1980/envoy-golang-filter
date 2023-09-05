package main

import (
	"fmt"
	"strconv"
        "regexp"
	"github.com/envoyproxy/envoy/contrib/golang/common/go/api"
)

var UpdateUpstreamBody = "upstream response body updated by the sophosFilter plugin"

type filter struct {
	api.PassThroughStreamFilter

	callbacks api.FilterCallbackHandler
	path      string
	config    *config
}

func (f *filter) sendLocalReplyInternal() api.StatusType {
	body := fmt.Sprintf("%s, %s %s %s path: %s\r\n", f.config.echoBody, f.config.tenantId, f.config.connectorId, f.config.regionId, f.path)
	f.callbacks.SendLocalReply(200, body, nil, 0, "")
	return api.LocalReply
}

// Callbacks which are called in request path
func (f *filter) DecodeHeaders(header api.RequestHeaderMap, endStream bool) api.StatusType {
	f.path, _ = header.Get(":path")
	if f.config.gatewayFQDN == "" {
	  header.Set("x-sophos-tid", f.config.tenantId)
	  header.Set("x-sophos-cid", f.config.connectorId)
	  header.Set("x-sophos-rid", f.config.regionId)
        } else {
	  host, _ := header.Get(":authority")
	  cookie, _ := header.Get("cookie")
          if host != f.config.gatewayFQDN && cookie != "" {
	      regexPattern := `_oauth2_proxy_[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}=[\w|\-=_]*;?`
	      re := regexp.MustCompile(regexPattern)
	      cleanedCookie := re.ReplaceAllString(cookie, "")
	      header.Set("cookie", cleanedCookie)
	  }
	}

	return api.Continue
}

/*
The callbacks can be implemented on demand

func (f *filter) DecodeData(buffer api.BufferInstance, endStream bool) api.StatusType {
	return api.Continue
}

func (f *filter) DecodeTrailers(trailers api.RequestTrailerMap) api.StatusType {
	return api.Continue
}
*/

func (f *filter) EncodeHeaders(header api.ResponseHeaderMap, endStream bool) api.StatusType {
	if f.path == "/update_upstream_response" {
		header.Set("Content-Length", strconv.Itoa(len(UpdateUpstreamBody)))
	}
	header.Set("Rsp-Header-From-Go", "bar-test")
	//header.Set("x-sophos-tid", f.config.tenantId)
	//header.Set("x-sophos-cid", f.config.connectorId)
	//header.Set("x-sophos-rid", f.config.regionId)
	return api.Continue
}

// Callbacks which are called in response path
func (f *filter) EncodeData(buffer api.BufferInstance, endStream bool) api.StatusType {
	if f.path == "/update_upstream_response" {
		if endStream {
			buffer.SetString(UpdateUpstreamBody)
		} else {
			// TODO implement buffer->Drain, buffer.SetString means buffer->Drain(buffer.Len())
			buffer.SetString("")
		}
	}
	return api.Continue
}

/*
The callbacks can be implemented on demand

func (f *filter) EncodeTrailers(trailers api.ResponseTrailerMap) api.StatusType {
	return api.Continue
}

func (f *filter) OnDestroy(reason api.DestroyReason) {
}
*/
