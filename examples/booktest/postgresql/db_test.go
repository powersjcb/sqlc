package booktest

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgtype"
	"testing"
	"time"

	"github.com/kyleconroy/sqlc/internal/sqltest"
)

func TestBooks(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second * 10)

	tx, cleanup := sqltest.PGXSQL(ctx, t, []string{"schema.sql"})
	defer cleanup()

	dq := New(tx)

	// create an author
	a, err := dq.CreateAuthor(ctx, "Unknown Master")
	if err != nil {
		t.Fatal(err)
	}

	// create transaction
	tx, err = tx.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}

	tq := dq.WithTx(tx)

	// save first book
	var now pgtype.Timestamp
	err = now.Set(time.Now())
	if err != nil {
		t.Fatal(err)
	}
	var tags1 pgtype.VarcharArray
	err = tags1.Set([]string{})
	if err != nil {
		t.Fatal(err)
	}
	_, err = tq.CreateBook(ctx, CreateBookParams{
		AuthorID:  a.AuthorID,
		Isbn:      "1",
		Title:     "my book title",
		Booktype:  BookTypeFICTION,
		Year:      2016,
		Available: now,
		Tags:      tags1,
	})
	if err != nil {
		t.Fatal(err)
	}

	// save second book
	var tags2 pgtype.VarcharArray
	err = tags2.Set([]string{"test", "tag2"})
	if err != nil {
		t.Fatal(err)
	}

	b1, err := tq.CreateBook(ctx, CreateBookParams{
		AuthorID:  a.AuthorID,
		Isbn:      "2",
		Title:     "the second book",
		Booktype:  BookTypeFICTION,
		Year:      2016,
		Available: now,
		Tags:      tags2,
	})
	if err != nil {
		t.Fatal(err)
	}

	// update the title and tags
	err = tq.UpdateBook(ctx, UpdateBookParams{
		BookID: b1.BookID,
		Title:  "changed second title",
		Tags:   []string{"cool", "disastor"},
	})
	if err != nil {
		t.Fatal(err)
	}

	var tags3 pgtype.VarcharArray
	err  = tags3.Set([]string{"cool"})
	if err != nil {
		t.Fatal(err)
	}

	// save third book
	_, err = tq.CreateBook(ctx, CreateBookParams{
		AuthorID:  a.AuthorID,
		Isbn:      "3",
		Title:     "the third book",
		Booktype:  BookTypeFICTION,
		Year:      2001,
		Available: now,
		Tags:      tags3,
	})
	if err != nil {
		t.Fatal(err)
	}
	var tags4 pgtype.VarcharArray
	err = tags4.Set([]string{"other"})
	if err != nil {
		t.Fatal(err)
	}
	// save fourth book
	b3, err := tq.CreateBook(ctx, CreateBookParams{
		AuthorID:  a.AuthorID,
		Isbn:      "4",
		Title:     "4th place finisher",
		Booktype:  BookTypeNONFICTION,
		Year:      2011,
		Available: now,
		Tags:      tags4,
	})
	if err != nil {
		t.Fatal(err)
	}

	// tx commit
	err = tx.Commit(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// upsert, changing ISBN and title
	err = dq.UpdateBookISBN(ctx, UpdateBookISBNParams{
		BookID: b3.BookID,
		Isbn:   "NEW ISBN",
		Title:  "never ever gonna finish, a quatrain",
		Tags:   []string{"someother"},
	})
	if err != nil {
		t.Fatal(err)
	}

	// retrieve first book
	books0, err := dq.BooksByTitleYear(ctx, BooksByTitleYearParams{
		Title: "my book title",
		Year:  2016,
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, book := range books0 {
		t.Logf("Book %d (%s): %s available: %s\n", book.BookID, book.Booktype, book.Title, book.Available.Time.Format(time.RFC822Z))
		author, err := dq.GetAuthor(ctx, book.AuthorID)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Book %d author: %s\n", book.BookID, author.Name)
	}

	// find a book with either "cool" or "other" tag
	t.Logf("---------\nTag search results:\n")
	res, err := dq.BooksByTags(ctx, []string{"cool", "other", "someother"})
	if err != nil {
		t.Fatal(err)
	}
	for _, ab := range res {
		t.Logf("Book %d: '%s', Author: '%s', ISBN: '%s' Tags: '%v'\n", ab.BookID, ab.Title, ab.Name, ab.Isbn, ab.Tags)
	}

	// TODO: call say_hello(varchar)

	// get book 4 and delete
	b5, err := dq.GetBook(ctx, b3.BookID)
	if err != nil {
		t.Fatal(err)
	}
	if err := dq.DeleteBook(ctx, b5.BookID); err != nil {
		t.Fatal(err)
	}

	// lookup empty books result
	books, err := dq.BooksByTitleYear(ctx, BooksByTitleYearParams{
		Title: "Unpublished Book",
		Year:  -1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(books) != 0 {
		t.Fatal("books should be empty")
	}

	// check correct encoding type
	data, err := json.Marshal(&books)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "null" {
		t.Fatalf("json.Marshal should encode null got: %s", string(data))
	}
}
