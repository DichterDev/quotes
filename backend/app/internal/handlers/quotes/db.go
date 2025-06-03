package quotes

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/segmentio/ksuid"
	. "github.com/segmentio/ksuid"
)

func logError(err error) {
	if err != nil {
		var pgErr pgconn.PgError
		if errors.As(err, pgErr) {
			fmt.Println(pgErr.Message)
			fmt.Println(pgErr.Code)
		}
	}
}

func initDB() (*pgx.Conn, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://user:password@localhost:5432/test_db"
		log.Println("DATABASE_URL environment variable not set. Using default connection string.")
	}

	conn, err := pgx.Connect(context.Background(), connStr)
	logError(err)

	defer conn.Close(context.Background())

	return conn, err
}

/*
	table collection {
		id char(27) [pk]
		name varchar(255)
		desc text
	}
*/
func insertCollection(collection Collection) (KSUID, error) {
	conn, _ := initDB()
	uuid := ksuid.New()

	_, err := conn.Exec(context.Background(), `INSERT INTO collections (id, name, description) VALUES($1, $2, $3)`, uuid, collection.Name, collection.Description)
	logError(err)

	for _, author := range collection.Authors {

		id := insertAuthor(author)
		relCollectionAuthor(uuid, id)
	}

	return uuid, err
}

/*
	table author {
	  id char(27) [pk]
	  name varchar(255)
	}
*/
func insertAuthor(author Author) KSUID {
	conn, _ := initDB()
	uuid := ksuid.New()

	_, err := conn.Exec(context.Background(), `INSERT INTO author (id, name) VALUES($1, $2)`, uuid, author.Name)
	logError(err)

	return uuid
}

/*
	table quote {
	  id char(27) [pk]
	  text text
	  author char(27) [fk]
	}
*/
func insertQuote(quote Quote) KSUID {
	conn, _ := initDB()
	uuid := ksuid.New()

	_, err := conn.Exec(context.Background(), `INSERT INTO quote (id, text, author) VALUES($1, $2)`, uuid, quote.Quote, quote.Author)
	logError(err)

	relCollectionQuote(quote.Collection, uuid)

	return uuid
}

/*
	table rel_collection_author {
	  id char(27) [pk]
	  collection char(27) [fk]
	  author char(27) [fk]
	}
*/
func relCollectionAuthor(collectionId KSUID, authorId KSUID) KSUID {
	conn, _ := initDB()
	uuid := ksuid.New()

	_, err := conn.Exec(context.Background(), `INSERT INTO rel_collection_author (id, collection, author) VALUES($1, $2, $3)`, uuid, collectionId, authorId)
	logError(err)

	return uuid
}

/*
	table rel_collection_quote {
	  id char(27) [pk]
	  collection char(27) [fk]
	  quote char(27) [fk]
	}
*/
func relCollectionQuote(collectionId KSUID, quoteId KSUID) KSUID {
	conn, _ := initDB()
	uuid := ksuid.New()

	_, err := conn.Exec(context.Background(), `INSERT INTO rel_collection_author (id, collection, quote) VALUES($1, $2, $3)`, uuid, collectionId, quoteId)
	logError(err)

	return uuid
}

func getCollection(id KSUID) Collection {
	conn, _ := initDB()

	var collection Collection

	err := conn.QueryRow(context.Background(), `SELECT name, desc from collection WHERE id = $1`, id).Scan(&collection.Name, &collection.Description)
	logError(err)

	rows, _ := conn.Query(
		context.Background(),
		`
		SELECT a.name
		FROM author a
		JOIN rel_collection_author rca ON a.id = rca.author
		WHERE rca.collection = $1
		`,
		id,
	)

	for rows.Next() {
		var author Author
		err := rows.Scan(&author.Name)
		logError(err)

		collection.Authors = append(collection.Authors, author)
	}

	return collection
}

func getAuthor(id KSUID) Author {
	conn, _ := initDB()

	var author Author

	err := conn.QueryRow(context.Background(), `SELECT name from collection WHERE id = $1`, id).Scan(&author.Name)
	logError(err)

	return author
}
