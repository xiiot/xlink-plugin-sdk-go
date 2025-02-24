package comctx

import (
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"time"
)

type Context struct {
	RequestID string
	StartTime time.Time
	Logger    *zap.Logger
}

func NewTraceContext(requestID string, logger *zap.Logger) *Context {
	c := &Context{requestID, time.Now(), logger}
	c.SetRequestID(requestID)
	c.SetTraceStartTime()
	return c
}

func (c *Context) SetRequestID(requestID string) {
	if requestID == "" {
		requestID = uuid.NewV4().String()
	}
	c.Logger = c.Logger.With(zap.Any("requestID", requestID))
	c.RequestID = requestID
}

func (c *Context) SetTraceStartTime() {
	t := time.Now()
	c.StartTime = t
	c.Logger = c.Logger.With(zap.Any("startTime", t))
}

func (c *Context) TraceCostTime(event string) {
	c.Logger = c.Logger.With(zap.Any(event+" cost time", time.Since(c.StartTime)))
}
