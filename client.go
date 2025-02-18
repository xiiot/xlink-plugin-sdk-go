package plugin

import (
	"errors"
	"github.com/hashicorp/go-plugin"
	"os/exec"
	"sync"
)

var PluginFactories = make(map[string]*Client)
var PluginLock sync.Mutex

type Client struct {
	pluginClient *plugin.Client
	path         string
	name         string
	service      *gRPCClient
	enable       bool
	on           bool

	sync.Mutex
}

func NewClient(name, path string) (*Client, error) {
	c := new(Client)
	c.enable = true
	c.path = path
	c.name = name

	err := c.Check()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) Check() error {
	c.Lock()
	defer c.Unlock()

	if !c.enable {
		return errors.New("plugin " + c.name + " disabled")
	}
	// grpc连接正常,直接返回
	if c.pluginClient != nil && !c.pluginClient.Exited() {
		return nil
	}
	c.on = false

	plugClient := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: Handshake,
		Plugins:         PluginMap,
		Cmd:             exec.Command("sh", "-c", c.path),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolGRPC},
	})

	rpcClient, err := plugClient.Client()
	if err != nil {
		return err
	}

	raw, err := rpcClient.Dispense(PluginName)
	if err != nil {
		return err
	}

	c.pluginClient = plugClient
	c.service = raw.(*gRPCClient)

	c.on = true

	return nil
}

func (c *Client) GetDriverInfo(req *Request) (*Response, error) {
	if err := c.Check(); err != nil {
		return nil, err
	}
	return c.service.GetDriverInfo(req)
}

func (c *Client) SetConfig(req *Request) (*Response, error) {
	if err := c.Check(); err != nil {
		return nil, err
	}
	return c.service.SetConfig(req)
}

func (c *Client) Setup(config *BackendConfig) (*Response, error) {
	if err := c.Check(); err != nil {
		return nil, err
	}
	return c.service.Setup(config)
}

func (c *Client) Start(req *Request) (*Response, error) {
	if err := c.Check(); err != nil {
		return nil, err
	}
	return c.service.Start(req)
}

func (c *Client) Restart(req *Request) (*Response, error) {
	if err := c.Check(); err != nil {
		return nil, err
	}
	return c.service.Restart(req)
}

func (c *Client) Stop(req *Request) (*Response, error) {
	if err := c.Check(); err != nil {
		return nil, err
	}
	return c.service.Stop(req)
}

func (c *Client) Set(req *Request) (*Response, error) {
	if err := c.Check(); err != nil {
		return nil, err
	}
	return c.service.Set(req)
}

func (c *Client) Get(req *Request) (*Response, error) {
	if err := c.Check(); err != nil {
		return nil, err
	}
	return c.service.Get(req)
}

func (c *Client) Disable() error {
	c.Lock()
	c.enable = false
	c.on = false
	client := c.pluginClient
	c.Unlock()

	if client != nil {
		client.Kill()
	}
	return nil
}

func (c *Client) Open() error {
	c.Lock()
	c.enable = true
	c.Unlock()

	return c.Check()
}

func (c *Client) Status() (enable, on bool) {
	c.Lock()
	defer c.Unlock()

	if !c.enable {
		return c.enable, c.on
	}
	if c.pluginClient != nil && !c.pluginClient.Exited() {
		c.on = true
	} else {
		c.on = false
	}
	return c.enable, c.on
}

type DriverConfig struct {
	Name string
	Path string
}

func RegisterPlugin(driver DriverConfig) (*Client, error) {
	PluginLock.Lock()
	defer PluginLock.Unlock()
	if c, ok := PluginFactories[driver.Name]; ok {
		return c, nil
	}
	c, err := NewClient(driver.Name, driver.Path)
	if err != nil {
		return nil, err
	}
	PluginFactories[driver.Name] = c

	return c, nil
}

func GetPlugin(name string) (*Client, error) {
	PluginLock.Lock()
	defer PluginLock.Unlock()
	if c, ok := PluginFactories[name]; ok {
		return c, nil
	}
	return nil, errors.New("plugin not found")
}

func ClosePlugin(name string) error {
	PluginLock.Lock()
	defer PluginLock.Unlock()
	if c, ok := PluginFactories[name]; ok {
		return c.Disable()
	}
	return errors.New("plugin not found")
}
