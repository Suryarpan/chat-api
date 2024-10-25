-- name: CreateMessage :one
INSERT INTO message_meta (
    from_pvt_id, to_pvt_id, mssg_status, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: CreateMessageType :one
INSERT INTO message_type_meta (
    mssg_id, mssg_type, attach_mssg_id
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: CreateMessageText :one
INSERT INTO message_text (
    mssg_id, mssg_body
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetMessageById :one
SELECT mm.*, mtm.mssg_type, mtm.attach_mssg_id, mt.mssg_body
FROM message_meta mm
JOIN message_type_meta mtm ON mtm.mssg_id = mm.mssg_id
JOIN message_text mt ON mt.mssg_id = mm.mssg_id
WHERE mm.mssg_id = $1;

-- name: GetMessageByIdPublic :one
SELECT mm.mssg_id, fu.user_id as from_user_id, tu.user_id as to_user_id, mm.mssg_status, mm.created_at, mm.updated_at, mtm.mssg_type, mtm.attach_mssg_id, mt.mssg_body
FROM message_meta mm
JOIN message_type_meta mtm ON mtm.mssg_id = mm.mssg_id
JOIN message_text mt ON mt.mssg_id = mm.mssg_id
JOIN users fu ON fu.pvt_id = mm.from_pvt_id
JOIN users tu ON tu.pvt_id = mm.to_pvt_id
WHERE mm.mssg_id = $1;
