package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"effective-mobile-song-library/internal/model"
)

type SongsRepository struct {
	db *sql.DB
}

func NewSongsRepository(db *sql.DB) *SongsRepository {
	return &SongsRepository{db: db}
}

func (sr *SongsRepository) Get(id uint64) (*model.SongInfo, error) {
	query := `
	SELECT song_id, group, song, release_date, text, link
	FROM songs
	WHERE song_id=$1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var songInfo model.SongInfo

	err := sr.db.QueryRowContext(ctx, query, id).Scan(
		&songInfo.ID,
		&songInfo.Group,
		&songInfo.Song,
		&songInfo.ReleaseDate,
		&songInfo.Text,
		&songInfo.Link,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &songInfo, nil
}

func (sr *SongsRepository) GetAll(filters model.SongFilters) ([]*model.SongInfo, error) {
	query := `
	SELECT song_id, group, song, release_date, song_text, link
	FROM songs
	WHERE ($1 = '' OR LOWER(group)=LOWER($1))
	AND ($2 = '' OR LOWER(song)=LOWER($2))
	AND ($3 = '' OR release_date LIKE '%' || $3 || '%')
	AND (
		$4 = '' OR
		EXISTS (
			SELECT 1
			FROM unnest(song_text) as verse
			WHERE verse LIKE '%' || $4 || '%'
		)
	)
	AND (link=$5 OR $5 = '')
	ORDER BY song_id ASC
	LIMIT $5 OFFSET $6`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		filters.Group,
		filters.Song,
		filters.ReleaseDate,
		filters.Text,
		filters.Link,
		filters.PageSize,
		(filters.Page - 1) * filters.PageSize,
	}

	rows, err := sr.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	songs := []*model.SongInfo{}

	for rows.Next() {
		var songInfo model.SongInfo
		err := rows.Scan(
			&songInfo.ID,
			&songInfo.Group,
			&songInfo.Song,
			&songInfo.ReleaseDate,
			&songInfo.Text,
			&songInfo.Link,
		)
		if err != nil {
			return nil, err
		}

		songs = append(songs, &songInfo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return songs, nil
}

func (sr *SongsRepository) GetText(filters model.SongTextFilters) (*string, error) {
	query := `
	SELECT song_text[$2]
	FROM songs
	WHERE (song_id = $1)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		filters.ID,
		filters.Verse,
	}

	rows, err := sr.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var verse string

	for rows.Next() {
		err := rows.Scan(&verse)
		if err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &verse, nil
}

func (sr *SongsRepository) Insert(songInfo *model.SongInfo) error {
	query := `
		INSERT INTO songs (group, song, release_date, text, link)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING song_id`

	args := []any{
		songInfo.Group,
		songInfo.Song,
		songInfo.ReleaseDate,
		songInfo.Text,
		songInfo.Link,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return sr.db.QueryRowContext(ctx, query, args...).Scan(&songInfo.ID)
}

func (sr *SongsRepository) Update(song *model.SongInfo) error {
	query := `
		UPDATE songs
		SET group = $2, song = $3, release_date = $4, link = $5, song_text = $6
		WHERE song_id = $1`

	args := []any{
		song.ID,
		song.Group,
		song.Song,
		song.ReleaseDate,
		song.Link,
		arrayIntoPlaceholders(song.Text, 6),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := sr.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (sr *SongsRepository) Delete(id uint64) error {
	query := `
	DELETE FROM songs
	WHERE song_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := sr.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
