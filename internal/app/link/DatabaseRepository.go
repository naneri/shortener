package link

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"strconv"
	"time"
)

type DatabaseRepository struct {
	dbConnection *sql.DB
}

func InitDatabaseRepository(db *sql.DB) (*DatabaseRepository, error) {
	dbRepo := DatabaseRepository{dbConnection: db}

	return &dbRepo, nil
}

func (repo *DatabaseRepository) GetLink(urlID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	number, err := strconv.ParseUint(urlID, 10, 32)

	if err != nil {
		return "", errors.New("unable to parse integer: " + urlID)
	}

	row := repo.dbConnection.QueryRowContext(ctx, "SELECT id, user_id, link FROM links WHERE id = $1 LIMIT 1", number)

	var dbLink Link

	err = row.Scan(&dbLink.ID, &dbLink.UserID, &dbLink.URL)

	if err != nil {
		return "", err
	}

	return dbLink.URL, nil
}

func (repo *DatabaseRepository) AddLink(link string, userID uint32) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	var id int
	if err := repo.dbConnection.QueryRowContext(ctx, "INSERT INTO public.links(user_id, link) VALUES ($1, $2) RETURNING id", userID, link).Scan(&id); err != nil {
		var pgError pq.Error

		if errors.As(err, &pgError) {
			if pgError.Code == pgerrcode.UniqueViolation {
				linkId, queryErr := repo.getLinkIdByUrl(ctx, link)

				if queryErr != nil {
					return 0, queryErr
				}
				return linkId, err
			}
		}

		return 0, err
	}

	return id, nil
}

func (repo *DatabaseRepository) getLinkIdByUrl(ctx context.Context, url string) (int, error) {
	row := repo.dbConnection.QueryRowContext(ctx, "SELECT id, user_id, link FROM links WHERE link = $1 LIMIT 1", url)

	var dbLink Link

	err := row.Scan(&dbLink.ID, &dbLink.UserID, &dbLink.URL)

	if err != nil {
		return 0, err
	}

	return dbLink.ID, nil

}

func (repo *DatabaseRepository) GetAllLinks() map[string]*Link {
	links := make(map[string]*Link)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	rows, _ := repo.dbConnection.QueryContext(ctx, "SELECT id, user_id, link FROM public.links")

	// обязательно закрываем перед возвратом функции
	defer rows.Close()

	for rows.Next() {
		var link Link
		_ = rows.Scan(&link.ID, &link.UserID, &link.URL)

		links[string(rune(link.ID))] = &link
	}

	return links
}
