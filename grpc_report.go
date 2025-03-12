package plugin

import (
	"context"
	"github.com/xiiot/xlink-plugin-sdk-go/proto"
)

var _ Report = &gRPCReportClient{}

type gRPCReportClient struct {
	client proto.ReportClient
}

func (m *gRPCReportClient) Post(req *Request) (*Response, error) {
	res, err := m.client.Post(context.Background(), &proto.RequestArgs{
		Request:   req.Req,
		RequestId: req.RequestID,
	})
	if err != nil {
		return nil, err
	}
	return &Response{Data: res.Data, RequestID: res.RequestId}, nil
}

func (m *gRPCReportClient) State(req *Request) (*Response, error) {
	res, err := m.client.State(context.Background(), &proto.RequestArgs{
		Request:   req.Req,
		RequestId: req.RequestID,
	})
	if err != nil {
		return nil, err
	}
	return &Response{Data: res.Data, RequestID: res.RequestId}, nil
}

type gRPCReportServer struct {
	proto.UnimplementedReportServer
	Impl Report
}

func (m *gRPCReportServer) Post(_ context.Context, req *proto.RequestArgs) (*proto.ResponseResult, error) {
	res, err := m.Impl.Post(&Request{
		Req:       req.Request,
		RequestID: req.RequestId,
	})
	if err != nil {
		return &proto.ResponseResult{}, err
	}
	return &proto.ResponseResult{Data: res.Data, RequestId: res.RequestID}, nil
}

func (m *gRPCReportServer) State(_ context.Context, req *proto.RequestArgs) (resp *proto.ResponseResult, err error) {
	res, err := m.Impl.State(&Request{
		Req:       req.Request,
		RequestID: req.RequestId,
	})
	if err != nil {
		return &proto.ResponseResult{}, err
	}
	return &proto.ResponseResult{Data: res.Data, RequestId: res.RequestID}, nil
}
