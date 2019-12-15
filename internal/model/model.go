package model

import "time"

type UserChat struct {
	ID        int       `json:"id" db:"id"`
	UserName  string    `json:"user_name" db:"user_name"`
	Name      string    `json:"name" db:"name"`
	ChatID    int64     `json:"chat_id" db:"chat_id"`
	Warning   bool      `json:"warning" db:"warning"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Coordinates - historical user positions
type Coordinates struct {
	ID            int       `json:"id" db:"id"`
	UserChatID    int       `json:"user_chat_id" db:"user_chat_id"`
	HotLocationID int       `json:"hot_location_id" db:"hot_location_id"`
	Longitude     float64   `json:"longitude" db:"longitude"`
	Latitude      float64   `json:"latitude" db:"latitude"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// HotLocation - meeting rectangle (pointX - north-west angle, pointY - south-east angle)
type HotLocation struct {
	ID        int       `json:"id" db:"id"`
	EventDate time.Time `json:"event_date" db:"event_date"`
	Name      string    `json:"name" db:"name"`
	PointXLat float64   `json:"point_x_lat" db:"point_x_lat"`
	PointXLon float64   `json:"point_x_lon" db:"point_x_lon"`

	PointYLat float64   `json:"point_y_lat" db:"point_y_lat"`
	PointYLon float64   `json:"point_y_lon" db:"point_y_lon"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type News struct {
	ID            int       `json:"id" db:"id"`
	Text          string    `json:"text" db:"text"`
	HotLocationID int       `json:"hot_location_id" db:"hot_location_id"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type Contacts struct {
	ID         int    `json:"id" db:"id"`
	Contact    string `json:"contact" db:"contact"`
	UserChatID int    `json:"user_chat_id" db:"user_chat_id"`
}
