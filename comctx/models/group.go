package models

type Group struct {
	DriverName string `json:"driver_name"`
	Name       string `json:"name"`
	OldName    string `gorm:"-" json:"old_name"`
	Interval   int64  `json:"interval"`
	//Node          Node           `gorm:"foreignKey:DriverName;references:Name"`
	Tags          []Tag          `gorm:"foreignKey:DriverName,GroupName;references:DriverName,Name"`
	Subscriptions []Subscription `gorm:"foreignKey:DriverName,GroupName;references:DriverName,Name"`
}
