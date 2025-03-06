package models

type Channel struct {
	ChannelName string              `json:"channel_name"`
	Devices     []KepwareDeviceData `json:"devices"`
}

type KepwareDeviceData struct {
	DeviceName string `json:"device_name"`
	Tags       []Tag  `json:"tags"`
}

type Tag struct {
	TagName string `json:"TagName"`
	TagId   string `json:"TagId"`
}
