package comctx

import "time"

const KeyDeviceName = "deviceName"
const KeyGroupName = "groupName"

const (
	DeviceInactivated = 0
	DeviceOnline      = 1
	DeviceOffline     = 2
	DeviceUnknown     = 3
)

type ReportProperty struct {
	Time  time.Time `json:"time,omitempty"`
	Value any       `json:"value,omitempty"`
}
