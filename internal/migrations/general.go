package migrations

import "database/sql"

func RunMigrations(db *sql.DB) error {

	_, err := db.Exec(`
	CREATE table IF NOT EXISTS links(
	    id      serial        constraint links_pk            primary key,
		user_id integer not null,    
		link    text    not null, 
		deleted_at timestamp
	); 
	CREATE unique index IF NOT EXISTS links_link_uindex on links (link);
`)

	return err
}

func DropTables(db *sql.DB) error {
	_, err := db.Exec("DROP table IF EXISTS links")

	return err
}
