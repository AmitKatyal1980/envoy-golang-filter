package main

import (

	xds "github.com/cncf/xds/go/xds/type/v3"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/envoyproxy/envoy/contrib/golang/common/go/api"
	"github.com/envoyproxy/envoy/contrib/golang/filters/http/source/go/pkg/http"
)

const Name = "sophosFilter"

func init() {
	http.RegisterHttpFilterConfigFactoryAndParser(Name, ConfigFactory, &parser{})
}

type config struct {
	echoBody string
	// other fields
	tenantId string
	connectorId string
	regionId    string
	gatewayFQDN string
}

type parser struct {
}

func (p *parser) Parse(any *anypb.Any) (interface{}, error) {
	configStruct := &xds.TypedStruct{}
	if err := any.UnmarshalTo(configStruct); err != nil {
		return nil, err
	}

	v := configStruct.Value
	conf := &config{}
        conf.tenantId, _ = v.AsMap()["tenantId"].(string)
        conf.connectorId, _ = v.AsMap()["connectorId"].(string)
        conf.regionId, _ = v.AsMap()["regionId"].(string)
        conf.gatewayFQDN, _ = v.AsMap()["gatewayFQDN"].(string)
	return conf, nil
}

func (p *parser) Merge(parent interface{}, child interface{}) interface{} {
	parentConfig := parent.(*config)
	childConfig := child.(*config)

	// copy one, do not update parentConfig directly.
	newConfig := *parentConfig
	if childConfig.echoBody != "" {
		newConfig.echoBody = childConfig.echoBody
	}
	return &newConfig
}

func ConfigFactory(c interface{}) api.StreamFilterFactory {
	conf, ok := c.(*config)
	if !ok {
		panic("unexpected config type")
	}

	return func(callbacks api.FilterCallbackHandler) api.StreamFilter {
		return &filter{
			callbacks: callbacks,
			config:    conf,
		}
	}
}

func main() {}
