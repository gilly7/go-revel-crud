package models

import (
	"context"
	"fmt"
	"go-revel-crud/app/db"
	"go-revel-crud/app/helpers"
	"strings"
)

const (
	createPostSQL = `insert into posts (username, title, content, created_at) values ($1, $2, $3, $4) returning id`
	getPostSQL    = `select id, username, title, content, created_at, updated_at from posts`
	getPostByID   = getPostSQL + ` where id=$1`
	updatePostSQL = `update posts set (username, title, content, updated_at) = ($1, $2, $3, $4) where id = $5`
	countPostSQL  = `select count(id) from posts`
	deletePostSQL = `delete from posts where id=$1`
)

type (
	Post struct {
		SequentialIdentifier
		Username string `json:"username"`
		Title    string `json:"title"`
		Content  string `json:"content"`
		Timestamps
	}
)

func (p *Post) All(
	ctx context.Context,
	db db.SQLOperations,
	filter *Filter,
) ([]*Post, error) {
	posts := make([]*Post, 0)

	query, args := p.buildQuery(
		getPostSQL,
		filter,
	)

	rows, err := db.QueryContext(
		ctx,
		query,
		args...,
	)
	defer rows.Close()
	if err != nil {
		return posts, err
	}

	for rows.Next() {
		var post Post
		err = rows.Scan(
			&post.ID,
			&post.Username,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return posts, err
		}
		posts = append(posts, &post)
	}

	return posts, err
}

func (q *Post) Count(
	ctx context.Context,
	db db.SQLOperations,
	filter *Filter,
) (int, error) {
	query, args := q.buildQuery(
		countPostSQL,
		&Filter{
			Term: filter.Term,
		},
	)
	var recordsCount int
	err := db.QueryRowContext(ctx, query, args...).Scan(&recordsCount)
	return recordsCount, err
}

func (p *Post) Delete(
	ctx context.Context,
	db db.SQLOperations,
) (int64, error) {
	res, err := db.ExecContext(ctx, deletePostSQL, p.ID)
	if err != nil {
		return 0, err
	}

	rowsDeleted, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsDeleted, nil
}

func (p *Post) ByID(
	ctx context.Context,
	db db.SQLOperations,
	id int64,
) (*Post, error) {
	row := db.QueryRowContext(ctx, getPostByID, id)
	return p.scan(row)
}

func (p *Post) Save(
	ctx context.Context,
	db db.SQLOperations,
) error {
	p.Timestamps.Touch()

	var err error
	if p.IsNew() {
		err := db.QueryRowContext(
			ctx,
			createPostSQL,
			p.Username,
			p.Title,
			p.Content,
			p.Timestamps.CreatedAt,
		).Scan(&p.ID)
		return err
	}
	_, err = db.ExecContext(
		ctx,
		updatePostSQL,
		p.Username,
		p.Title,
		p.Content,
		p.Timestamps.UpdatedAt,
		p.ID,
	)
	return err
}

func (*Post) scan(
	row db.RowScanner,
) (*Post, error) {
	var p Post
	err := row.Scan(
		&p.ID,
		&p.Username,
		&p.Title,
		&p.Content,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	return &p, err
}

func (p *Post) buildQuery(
	query string,
	filter *Filter,
) (string, []interface{}) {
	conditions := make([]string, 0)
	args := make([]interface{}, 0)
	placeholder := helpers.NewPlaceholder()

	if filter.Term != "" {
		likeStmt := make([]string, 0)
		columns := []string{"username", "title", "content"}
		for _, col := range columns {
			search := fmt.Sprintf(" (lower(%s) like '%%' || $%d || '%%')", col, placeholder.Touch())
			likeStmt = append(likeStmt, search)
			args = append(args, filter.Term)
		}
		conditions = append(conditions, fmt.Sprintf("(%s)", strings.Join(likeStmt, " or")))
	}

	if len(conditions) > 0 {
		query += " where" + strings.Join(conditions, " and")
	}

	if filter.Per > 0 && filter.Page > 0 {
		query += fmt.Sprintf(" order by id desc limit $%d offset $%d", placeholder.Touch(), placeholder.Touch())
		args = append(args, filter.Per, (filter.Page-1)*filter.Per)
	}

	return query, args
}
