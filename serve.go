package plugin

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"math"
	"os"
	"time"
)

var parentPid int

func init() {
	parentPid = os.Getppid()
}

type ServeOpts struct {
	FactoryFunc Factory
	Logger      hclog.Logger
}

func checkParentAlive() {
	go func() {
		for {
			if os.Getppid() != parentPid {
				fmt.Println("parent no alive, exit")
				os.Exit(0)
			}
			_, err := os.FindProcess(parentPid)
			if err != nil {
				fmt.Println("parent no alive, exit")
				os.Exit(0)
			}
			time.Sleep(5 * time.Second)
		}
	}()
}

func Serve(opts *ServeOpts) error {
	checkParentAlive()

	logger := opts.Logger
	if logger == nil {
		logger = hclog.New(&hclog.LoggerOptions{
			Level:           hclog.Debug,
			Output:          hclog.DefaultOutput,
			JSONFormat:      true,
			IncludeLocation: false,
		})
	}

	pluginSets := map[int]plugin.PluginSet{
		1: {
			PluginName: &DriverGRPCPlugin{
				Factory: opts.FactoryFunc,
				Logger:  logger,
			},
		},
	}
	serveOpts := &plugin.ServeConfig{
		HandshakeConfig:  Handshake,
		VersionedPlugins: pluginSets,
		Logger:           logger,
		GRPCServer: func(opts []grpc.ServerOption) *grpc.Server {
			opts = append(opts, grpc.MaxRecvMsgSize(math.MaxInt32))
			opts = append(opts, grpc.MaxSendMsgSize(math.MaxInt32))
			return plugin.DefaultGRPCServer(opts)
		},
	}

	plugin.Serve(serveOpts)

	return nil
}
