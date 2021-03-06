package link

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

	row := repo.dbConnection.QueryRowContext(ctx, "SELECT id, user_id, link, deleted_at FROM links WHERE id = $1 LIMIT 1", number)

	var dbLink Link

	err = row.Scan(&dbLink.ID, &dbLink.UserID, &dbLink.URL, &dbLink.DeletedAt)

	if err != nil {
		return "", err
	}

	if dbLink.DeletedAt.Valid {
		return dbLink.URL, &ModelDeletedError{msg: fmt.Sprintf("Deleted on:" + dbLink.DeletedAt.Time.String())}
	}

	return dbLink.URL, nil
}

func (repo *DatabaseRepository) AddLink(link string, userID uint32) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	var id int

	if err := repo.dbConnection.QueryRowContext(ctx, "INSERT INTO public.links(user_id, link) VALUES ($1, $2) RETURNING id", userID, link).Scan(&id); err != nil {
		var pgError *pq.Error

		if errors.As(err, &pgError) {

			if pgError.Code == pgerrcode.UniqueViolation {
				linkID, queryErr := repo.getLinkIDByURL(ctx, link)

				if queryErr != nil {
					return 0, queryErr
				}

				return linkID, err
			}
		}
		return 0, err
	}

	return id, nil
}

func (repo *DatabaseRepository) getLinkIDByURL(ctx context.Context, url string) (int, error) {
	row := repo.dbConnection.QueryRowContext(ctx, "SELECT id, user_id, link FROM links WHERE link = $1 LIMIT 1", url)

	var dbLink Link

	err := row.Scan(&dbLink.ID, &dbLink.UserID, &dbLink.URL)

	if err != nil {
		return 0, err
	}

	return dbLink.ID, nil

}

func (repo *DatabaseRepository) GetAllLinks() (map[string]*Link, error) {
	links := make(map[string]*Link)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	rows, err := repo.dbConnection.QueryContext(ctx, "SELECT id, user_id, link FROM public.links")

	if err != nil || rows.Err() != nil {
		return links, err
	}

	// ?????????????????????? ?????????????????? ?????????? ?????????????????? ??????????????
	defer rows.Close()

	for rows.Next() {
		var link Link
		_ = rows.Scan(&link.ID, &link.UserID, &link.URL)

		links[string(rune(link.ID))] = &link
	}

	return links, nil
}

func (repo *DatabaseRepository) DeleteLinks(ids []string) error {
	linksToDelete := make([]int, 0, len(ids))

	for _, id := range ids {
		i, scanErr := strconv.Atoi(id)

		if scanErr != nil {
			return fmt.Errorf("error parsing ids: %v", scanErr)
		}

		linksToDelete = append(linksToDelete, i)
	}

	ctx := context.Background()

	tx, err := repo.dbConnection.Begin()
	if err != nil {
		return err
	}
	// ?????? 1.1 ??? ???????? ?????????????????? ????????????, ???????????????????? ??????????????????
	defer tx.Rollback()

	// ?????? 2 ??? ?????????????? ????????????????????
	stmt, err := tx.PrepareContext(ctx, `
					UPDATE links
					SET deleted_at = $1
					WHERE id = $2
				`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	for _, linkID := range linksToDelete {
		if _, updateErr := stmt.ExecContext(ctx, time.Now(), linkID); updateErr != nil {
			return fmt.Errorf("error deleting the links: %v", updateErr)
		}
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return fmt.Errorf("error deleting the links: %v", commitErr)
	}

	return nil
}
