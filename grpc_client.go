package plugin

import (
	"context"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/xiiot/xlink-plugin-sdk-go/proto"
	"google.golang.org/grpc"
	"net"
)

var _ Driver = &gRPCClient{}

type gRPCClient struct {
	client proto.DriverClient
	broker *plugin.GRPCBroker
	logger hclog.Logger
}

func (c *gRPCClient) GetDriverInfo(req *Request) (*Response, error) {
	res, err := c.client.GetDriverInfo(context.Background(), &proto.RequestArgs{
		Request: req.Req,
	})
	if err != nil {
		return nil, err
	}
	return &Response{
		Data: res.Data,
	}, nil
}

func (c *gRPCClient) SetConfig(req *Request) (*Response, error) {
	res, err := c.client.SetConfig(context.Background(), &proto.RequestArgs{
		Request: req.Req,
	})
	if err != nil {
		return nil, err
	}
	return &Response{
		Data: res.Data,
	}, nil
}

func (c *gRPCClient) Setup(config *BackendConfig) (*Response, error) {
	reportImpl := config.ReportSvc
	report := &gRPCReportServer{
		Impl: reportImpl,
	}

	// 直接通过 grpc 启动 report server
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, err
	}
	brokerID := lis.Addr().(*net.TCPAddr).Port
	grpcServer := grpc.NewServer()
	proto.RegisterReportServer(grpcServer, report)
	go func() {
		er := grpcServer.Serve(lis)
		if er != nil {
			c.logger.Error("failed to start grpc report server", er)
		}
	}()

	res, err := c.client.Setup(context.Background(), &proto.RequestArgs{
		Request:  config.DriverName,
		PluginId: uint32(brokerID),
	})
	if err != nil {
		return nil, err
	}

	return &Response{
		Data: res.Data,
	}, nil
}

func (c *gRPCClient) Start(req *Request) (*Response, error) {
	_, err := c.client.Start(context.Background(), &proto.RequestArgs{
		Request: req.Req,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (c *gRPCClient) Restart(req *Request) (*Response, error) {
	_, err := c.client.Restart(context.Background(), &proto.RequestArgs{
		Request: req.Req,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (c *gRPCClient) Stop(req *Request) (*Response, error) {
	_, err := c.client.Stop(context.Background(), &proto.RequestArgs{
		Request: req.Req,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (c *gRPCClient) Get(req *Request) (*Response, error) {
	_, err := c.client.Get(context.Background(), &proto.RequestArgs{
		Request: req.Req,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (c *gRPCClient) Set(req *Request) (*Response, error) {
	_, err := c.client.Set(context.Background(), &proto.RequestArgs{
		Request: req.Req,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}
