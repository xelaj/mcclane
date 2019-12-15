package pg

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/xelaj/mcclane/internal/model"
	"log"
)

var (
	NotFoundErr = fmt.Errorf("not found")
)

type DB struct {
	conn *sqlx.DB
}

func NewDB(conn *sqlx.DB) *DB {
	return &DB{
		conn: conn,
	}
}

func (db *DB) AddUserChat(ctx context.Context, uc model.UserChat) (int64, error) {
	res, err := db.conn.ExecContext(ctx, `INSERT INTO user_chat (user_name, chat_id) VALUES ($1, $2)`, uc.UserName, uc.ChatID)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (db *DB) UpdateUserChatName(ctx context.Context, chatID int64, name string) error {
	_, err := db.conn.ExecContext(ctx, `UPDATE user_chat SET name = $1 WHERE chatt_id = $2`, name, chatID)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) AddCoordinates(ctx context.Context, uc model.Coordinates) (int64, error) {
	res, err := db.conn.ExecContext(ctx, `INSERT INTO coordinates (latitude, longitude, user_chat_id, hot_location_id) VALUES ($1, $2, $3, $4)`, uc.Latitude, uc.Longitude, uc.UserChatID, uc.HotLocationID)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (db *DB) AddHotLocation(ctx context.Context, uc model.HotLocation) (int64, error) {
	return 0, nil
}

func (db *DB) GetUserChat(ctx context.Context, id int) (*model.UserChat, error) {
	rr, err := db.conn.QueryxContext(ctx, `SELECT * FROM user_chat WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}
	if rr.Next() {
		res := new(model.UserChat)
		err := rr.StructScan(res)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
	return nil, NotFoundErr
}

func (db *DB) GetCoordinates(ctx context.Context, id int) (*model.Coordinates, error) {
	rr, err := db.conn.QueryxContext(ctx, `SELECT * FROM coordinates WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}
	if rr.Next() {
		res := new(model.Coordinates)
		err := rr.StructScan(res)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
	return nil, NotFoundErr
}

func (db *DB) GetHotLocation(ctx context.Context, id int) (*model.HotLocation, error) {
	rr, err := db.conn.QueryxContext(ctx, `SELECT * FROM hot_location WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}
	if rr.Next() {
		res := new(model.HotLocation)
		err := rr.StructScan(res)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
	return nil, NotFoundErr
}

func (db *DB) GetUserChatByChatID(ctx context.Context, id int64) (*model.UserChat, error) {
	rr, err := db.conn.QueryxContext(ctx, `SELECT * FROM user_chat WHERE chat_id=$1 ORDER BY id DESC LIMIT 1`, id)
	if err != nil {
		return nil, err
	}
	if rr.Next() {
		res := new(model.UserChat)
		err := rr.StructScan(res)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
	return nil, NotFoundErr
}

func (db *DB) GetLastCoordinatesByUserChatID(ctx context.Context, ucID int) (*model.Coordinates, error) {
	rr, err := db.conn.QueryxContext(ctx, `SELECT * FROM coordinates WHERE user_chat_id=$1 ORDER BY created_at DESC LIMIT 1`, ucID)
	if err != nil {
		return nil, err
	}
	if rr.Next() {
		res := new(model.Coordinates)
		err := rr.StructScan(res)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
	return nil, NotFoundErr
}

func (db *DB) GetLastCoordinatesByHotLocationID(ctx context.Context, hlID int) ([]*model.Coordinates, error) {
	rr, err := db.conn.QueryxContext(ctx, `SELECT * FROM coordinates WHERE hot_location_id=$1`, hlID)
	if err != nil {
		return nil, err
	}
	result := make([]*model.Coordinates, 0)
	for rr.Next() {
		res := new(model.Coordinates)
		err := rr.StructScan(res)
		if err != nil {
			log.Println("")
			continue
		}
		result = append(result, res)
	}
	return result, nil
}

func (db *DB) GetHotLocationByPoint(ctx context.Context, lat, lon float64) (*model.HotLocation, error) {
	rr, err := db.conn.QueryxContext(ctx, `SELECT * FROM hot_location WHERE point_x_lat > $1 AND point_y_lat < $1 AND point_x_lon < $2 AND point_y_lon > $2`, lat, lon)
	if err != nil {
		return nil, err
	}
	if rr.Next() {
		res := new(model.HotLocation)
		err := rr.StructScan(res)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
	return nil, NotFoundErr
}

func (db *DB) AddNews(ctx context.Context, uc model.News) (int64, error) {
	res, err := db.conn.ExecContext(ctx, `INSERT INTO news (text, hot_location_id) VALUES ($1, $2)`, uc.Text, uc.HotLocationID)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (db *DB) GetLastNews(ctx context.Context, id int) ([]*model.News, error) {
	rr, err := db.conn.QueryxContext(ctx, `SELECT * FROM news WHERE id > $1`, id)
	if err != nil {
		return nil, err
	}
	res := make([]*model.News, 0)
	for rr.Next() {
		n := new(model.News)
		err := rr.StructScan(res)
		if err != nil {
			continue
		}
		res = append(res, n)
	}
	return res, nil
}
