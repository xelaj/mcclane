package model

import "time"

type UserChat struct {
	ID       int
	UserName string
	ChatID   int64
	Warning  bool
}

// Coordinates - historical user positions
type Coordinates struct {
	ID            int
	UserChatID    int
	HotLocationID int
	Longitude     float64
	Latitude      float64
	CreatedAt     time.Time
}

// HotLocation - meeting rectangle (pointX - north-west angle, pointY - south-east angle)
type HotLocation struct {
	ID        int
	EventDate time.Time
	Name      string
	PointXLat float64
	PointXLng float64

	PointYLat float64
	PointYLng float64
}
