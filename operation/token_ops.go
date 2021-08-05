package operation

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/lib/pq"
)

type Token struct {
	token     string
	frequency int
}

// GenerateToken generate tokens in lowercase
// It returns the generated in bytes
func GenerateToken() []byte {
	letters := []byte("abcdefghijklmnopqrstuvwxyz")
	bytes := make([]byte, 10)
	for i := 0; i < 10; i++ {
		bytes[i] = letters[rand.Int63()%int64(len(letters))]
	}
	return bytes
}

type TokenReadWriter interface {
	WriteTokenToFile(path string, writes_range int, token func() []byte) error
	TokenToStdout(path string) error
	TokenToDatabase(path string, db *sql.DB) error
}

type tokenWriter struct {
}

func NewToken() *tokenWriter {
	return &tokenWriter{}
}

func (tr *tokenWriter) readTokenFromFile(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (tr *tokenWriter) TokenToStdout(path string) error {

	file, err := tr.readTokenFromFile(path)

	if err != nil {
		return err
	}

	defer file.Close()

	w := bufio.NewWriter(os.Stdout)
	scanner := bufio.NewScanner(file)

	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	for scanner.Scan() {
		fmt.Fprintln(w, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "shouldn't see an error scanning a string")
	}

	return nil

}

func (tr *tokenWriter) TokenToDatabase(path string, db *sql.DB) error {

	tokenList := []interface{}{} // A slice to store the tokens temporary

	file, err := tr.readTokenFromFile(path)

	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	for scanner.Scan() {

		tokenList = append(tokenList, scanner.Text()) // append the token being read from the file to the tokenList slice

	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "shouldn't see an error scanning a string")
	}
	tx, err := db.Begin() // start a db transaction

	if err != nil {
		return err
	}

	_, err = tx.Exec(`CREATE TEMP TABLE token_temp
	 ON COMMIT DROP
	AS SELECT * FROM token
	WITH NO DATA`) // create a temporary database table

	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(pq.CopyIn("token_temp", "token")) // passed in pq.CopyIn for bulk imports

	if err != nil {

		return err
	}

	for _, token := range tokenList {
		_, err = stmt.Exec(token) // insert token to the temporary table

		if err != nil {
			tx.Rollback()
			return err
		}

	}

	_, err = stmt.Exec()

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = tx.Exec(`INSERT INTO token (token, frequency)
	SELECT token, COUNT(*) - 1 FROM token_temp
	GROUP BY token
 `) // insert into the database table from the temporary table
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()

	if err != nil {
		return err
	}

	ctx := context.Background()

	var nonUniqeTokens []Token

	query := "SELECT token, frequency FROM token WHERE frequency > $1"

	rows, err := db.QueryContext(ctx, query, 0) // quey for tokens with the frequency greater than 0

	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var token Token
		if err := rows.Scan(&token.token, &token.frequency); err != nil {
			log.Fatal(err)
		}
		nonUniqeTokens = append(nonUniqeTokens, token)
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, nonUniqeTokens)

	return nil

}

func (tr *tokenWriter) WriteTokenToFile(path string, writes_range int, token func() []byte) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		return err
	}

	defer f.Close()

	w := bufio.NewWriter(f)

	defer w.Flush()

	defer f.Sync()

	for i := 0; i < writes_range; i++ {
		data := token()
		data[len(data)-1] = '\n'
		w.Write(data)
	}

	return nil
}
