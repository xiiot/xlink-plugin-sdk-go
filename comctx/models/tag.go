package models

type Tag struct {
	DriverName  string  `json:"driver_name"`
	GroupName   string  `json:"group_name"`
	Name        string  `json:"name"`
	Address     string  `json:"address"`
	Attribute   int64   `json:"attribute"`
	Precision   int64   `json:"precision"`
	Decimal     float32 `json:"decimal"`
	Bias        float32 `json:"bias"`
	Type        int64   `json:"type"`
	Description string  `json:"description"`
	Value       string  `json:"value"`
	//Node        Node    `gorm:"foreignKey:DriverName;references:Name"`
	//Group       Group   `gorm:"foreignKey:DriverName,GroupName;references:DriverName,Name"`
}
