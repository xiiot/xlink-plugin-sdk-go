package plugin

import (
	"context"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/xiiot/xlink-plugin-sdk-go/proto"
	"google.golang.org/grpc"
)

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	// This isn't required when using VersionedPlugins
	ProtocolVersion:  1,
	MagicCookieKey:   "DRIVER_PLUGIN",
	MagicCookieValue: "ON",
}

const PluginName = "driver"

// PluginMap is the map of plugins we can dispense.
var PluginMap = map[string]plugin.Plugin{
	PluginName: &DriverGRPCPlugin{},
}

type Response struct {
	Data      string `json:"data"`
	RequestID string `json:"requestID"`
}

type Request struct {
	BrokerID  uint32 `json:"brokerID"`
	Req       string `json:"req"`
	RequestID string `json:"requestID"`
}

type Driver interface {
	// GetDriverInfo 获取驱动信息
	GetDriverInfo(req *Request) (*Response, error)
	// SetConfig 配置驱动，目前只配置了驱动的配置文件路径
	SetConfig(req *Request) (*Response, error)
	//UpdateConfig 更新驱动配置
	UpdateConfig(req *Request) (*Response, error)
	// Setup 宿主进程上报接口传递，必须调用下述逻辑，其余可用户自定义
	Setup(config *BackendConfig) (*Response, error)
	// Start 驱动采集启动，用户自定义实现
	Start(req *Request) (*Response, error)
	// Restart 驱动重启，用户自定义实现
	Restart(req *Request) (*Response, error)
	// Stop 驱动停止，用户自定义实现
	Stop(req *Request) (*Response, error)

	// Get 召测，用户自定义实现
	Get(req *Request) (*Response, error)
	// Set 置数，用户自定义实现
	Set(req *Request) (*Response, error)
}

// Report 驱动 --> 宿主
type Report interface {
	Post(req *Request) (*Response, error)
	State(req *Request) (*Response, error)
}

type BackendConfig struct {
	DriverName string
	ReportSvc  Report
	Logger     hclog.Logger
}

type Factory func(context.Context, *BackendConfig) (Driver, error)

type DriverGRPCPlugin struct {
	plugin.Plugin
	Factory Factory
	Logger  hclog.Logger
}

func (p *DriverGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterDriverServer(s, &gRPCServer{
		broker:  broker,
		factory: p.Factory,
		logger:  p.Logger,
	})
	return nil
}

func (p *DriverGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (any, error) {
	return &gRPCClient{
		client: proto.NewDriverClient(c),
		broker: broker,
		logger: p.Logger,
	}, nil
}

var _ plugin.GRPCPlugin = &DriverGRPCPlugin{}
