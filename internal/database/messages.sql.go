// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: messages.sql

package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

const createMessage = `-- name: CreateMessage :one
INSERT INTO message_meta (
    from_pvt_id, to_pvt_id, mssg_status, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING mssg_id, from_pvt_id, to_pvt_id, mssg_status, created_at, updated_at
`

type CreateMessageParams struct {
	FromPvtID  int32         `json:"from_pvt_id"`
	ToPvtID    int32         `json:"to_pvt_id"`
	MssgStatus MessageStatus `json:"mssg_status"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

func (q *Queries) CreateMessage(ctx context.Context, arg CreateMessageParams) (MessageMetum, error) {
	row := q.db.QueryRow(ctx, createMessage,
		arg.FromPvtID,
		arg.ToPvtID,
		arg.MssgStatus,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i MessageMetum
	err := row.Scan(
		&i.MssgID,
		&i.FromPvtID,
		&i.ToPvtID,
		&i.MssgStatus,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createMessageText = `-- name: CreateMessageText :one
INSERT INTO message_text (
    mssg_id, mssg_body
) VALUES (
    $1, $2
) RETURNING mssg_id, mssg_body
`

type CreateMessageTextParams struct {
	MssgID   int64  `json:"mssg_id"`
	MssgBody string `json:"mssg_body"`
}

func (q *Queries) CreateMessageText(ctx context.Context, arg CreateMessageTextParams) (MessageText, error) {
	row := q.db.QueryRow(ctx, createMessageText, arg.MssgID, arg.MssgBody)
	var i MessageText
	err := row.Scan(&i.MssgID, &i.MssgBody)
	return i, err
}

const createMessageType = `-- name: CreateMessageType :one
INSERT INTO message_type_meta (
    mssg_id, mssg_type, attach_mssg_id
) VALUES (
    $1, $2, $3
) RETURNING mssg_id, mssg_type, attach_mssg_id
`

type CreateMessageTypeParams struct {
	MssgID       int64       `json:"mssg_id"`
	MssgType     MessageType `json:"mssg_type"`
	AttachMssgID pgtype.Int8 `json:"attach_mssg_id"`
}

func (q *Queries) CreateMessageType(ctx context.Context, arg CreateMessageTypeParams) (MessageTypeMetum, error) {
	row := q.db.QueryRow(ctx, createMessageType, arg.MssgID, arg.MssgType, arg.AttachMssgID)
	var i MessageTypeMetum
	err := row.Scan(&i.MssgID, &i.MssgType, &i.AttachMssgID)
	return i, err
}

const getMessageById = `-- name: GetMessageById :one
SELECT mm.mssg_id, mm.from_pvt_id, mm.to_pvt_id, mm.mssg_status, mm.created_at, mm.updated_at, mtm.mssg_type, mtm.attach_mssg_id, mt.mssg_body
FROM message_meta mm
JOIN message_type_meta mtm ON mtm.mssg_id = mm.mssg_id
JOIN message_text mt ON mt.mssg_id = mm.mssg_id
WHERE mm.mssg_id = $1
`

type GetMessageByIdRow struct {
	MssgID       int64         `json:"mssg_id"`
	FromPvtID    int32         `json:"from_pvt_id"`
	ToPvtID      int32         `json:"to_pvt_id"`
	MssgStatus   MessageStatus `json:"mssg_status"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	MssgType     MessageType   `json:"mssg_type"`
	AttachMssgID pgtype.Int8   `json:"attach_mssg_id"`
	MssgBody     string        `json:"mssg_body"`
}

func (q *Queries) GetMessageById(ctx context.Context, mssgID int64) (GetMessageByIdRow, error) {
	row := q.db.QueryRow(ctx, getMessageById, mssgID)
	var i GetMessageByIdRow
	err := row.Scan(
		&i.MssgID,
		&i.FromPvtID,
		&i.ToPvtID,
		&i.MssgStatus,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.MssgType,
		&i.AttachMssgID,
		&i.MssgBody,
	)
	return i, err
}

const getMessageByIdPublic = `-- name: GetMessageByIdPublic :one
SELECT mm.mssg_id, fu.user_id as from_user_id, tu.user_id as to_user_id, mm.mssg_status, mm.created_at, mm.updated_at, mtm.mssg_type, mtm.attach_mssg_id, mt.mssg_body
FROM message_meta mm
JOIN message_type_meta mtm ON mtm.mssg_id = mm.mssg_id
JOIN message_text mt ON mt.mssg_id = mm.mssg_id
JOIN users fu ON fu.pvt_id = mm.from_pvt_id
JOIN users tu ON tu.pvt_id = mm.to_pvt_id
WHERE mm.mssg_id = $1
`

type GetMessageByIdPublicRow struct {
	MssgID       int64         `json:"mssg_id"`
	FromUserID   pgtype.UUID   `json:"from_user_id"`
	ToUserID     pgtype.UUID   `json:"to_user_id"`
	MssgStatus   MessageStatus `json:"mssg_status"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	MssgType     MessageType   `json:"mssg_type"`
	AttachMssgID pgtype.Int8   `json:"attach_mssg_id"`
	MssgBody     string        `json:"mssg_body"`
}

func (q *Queries) GetMessageByIdPublic(ctx context.Context, mssgID int64) (GetMessageByIdPublicRow, error) {
	row := q.db.QueryRow(ctx, getMessageByIdPublic, mssgID)
	var i GetMessageByIdPublicRow
	err := row.Scan(
		&i.MssgID,
		&i.FromUserID,
		&i.ToUserID,
		&i.MssgStatus,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.MssgType,
		&i.AttachMssgID,
		&i.MssgBody,
	)
	return i, err
}
