package migrations

import "database/sql"

func RunMigrations(db *sql.DB) error {

	_, err := db.Exec(" create table IF NOT EXISTS links(id      serial        constraint links_pk            primary key,    user_id integer not null,    link    text    not null);")

	return err
}

func DropTables(db *sql.DB) error {
	_, err := db.Exec("DROP table IF EXISTS links")

	return err
}
