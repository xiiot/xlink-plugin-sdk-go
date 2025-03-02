package models

type Subscription struct {
	AppName    string `json:"app_name"`
	DriverName string `json:"driver_name"`
	GroupName  string `json:"group_name"`
	Params     string `json:"params"`
	Node       Node   `gorm:"foreignKey:AppName;references:Name"`
	Group      Group  `gorm:"foreignKey:DriverName,GroupName;references:DriverName,Name"`
}
