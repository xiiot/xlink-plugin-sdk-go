package models

type Node struct {
	Name          string         `json:"name"`
	Type          int64          `json:"type"`
	State         int64          `json:"state"`
	PluginName    string         `json:"plugin_name"`
	Groups        []Group        `gorm:"foreignKey:DriverName;references:Name"`
	Setting       Setting        `gorm:"foreignKey:NodeName;references:Name"`
	Subscriptions []Subscription `gorm:"foreignKey:AppName;references:Name"`
}
