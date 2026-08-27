package store

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"filippo.io/age"
	_ "modernc.org/sqlite"
)

type Entry struct {
	ID       int64
	Service  string
	Username string
	Password string // decrypted
	Created  string
}

type Store struct {
	db         *sql.DB
	passphrase string
	path       string
}

func Open(path, passphrase string) (*Store, error) {
	if passphrase == "" {
		return nil, fmt.Errorf("passphrase required (set PASSMAN_MASTER)")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", path)
	if err != nil { return nil, err }
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil { db.Close(); return nil, err }
	if _, err := db.Exec("PRAGMA busy_timeout=5000"); err != nil { db.Close(); return nil, err }
	_ = os.Chmod(path, 0600)
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS entries(id INTEGER PRIMARY KEY AUTOINCREMENT, service TEXT UNIQUE, username TEXT, password TEXT, created TEXT)`); err != nil {
		return nil, err
	}
	return &Store{db: db, passphrase: passphrase, path: path}, nil
}

func (s *Store) Close() error { return s.db.Close() }

func encrypt(passphrase, plaintext string) (string, error) {
	recipient, err := age.NewScryptRecipient(passphrase)
	if err != nil { return "", err }
	var buf bytes.Buffer
	w, err := age.Encrypt(&buf, recipient)
	if err != nil { return "", err }
	if _, err := io.WriteString(w, plaintext); err != nil { return "", err }
	if err := w.Close(); err != nil { return "", err }
	return buf.String(), nil
}

func decrypt(passphrase, ciphertext string) (string, error) {
	identity, err := age.NewScryptIdentity(passphrase)
	if err != nil { return "", err }
	r, err := age.Decrypt(stringsReader(ciphertext), identity)
	if err != nil { return "", err }
	b, err := io.ReadAll(r)
	return string(b), err
}

func stringsReader(s string) io.Reader { return bytes.NewBufferString(s) }

func (s *Store) Add(service, username, password string) (int64, error) {
	if service == "" || password == "" {
		return 0, fmt.Errorf("service and password required")
	}
	ct, err := encrypt(s.passphrase, password)
	if err != nil { return 0, err }
	res, err := s.db.Exec("INSERT OR REPLACE INTO entries(service, username, password, created) VALUES(?,?,?, datetime('now'))", service, username, ct)
	if err != nil { return 0, err }
	return res.LastInsertId()
}

func (s *Store) Get(service string) (*Entry, error) {
	var e Entry
	var ct string
	err := s.db.QueryRow("SELECT id, service, username, password, created FROM entries WHERE service=?", service).Scan(&e.ID, &e.Service, &e.Username, &ct, &e.Created)
	if err != nil { return nil, err }
	pw, err := decrypt(s.passphrase, ct)
	if err != nil { return nil, err }
	e.Password = pw
	return &e, nil
}

func (s *Store) List() ([]Entry, error) {
	rows, err := s.db.Query("SELECT id, service, username, password, created FROM entries ORDER BY service")
	if err != nil { return nil, err }
	defer rows.Close()
	var out []Entry
	for rows.Next() {
		var e Entry
		var ct string
		if err := rows.Scan(&e.ID, &e.Service, &e.Username, &ct, &e.Created); err != nil { return nil, err }
		// don't decrypt all for list (for speed), just service/username
		e.Password = ""
		out = append(out, e)
	}
	return out, rows.Err()
}

func (s *Store) Delete(service string) error {
	_, err := s.db.Exec("DELETE FROM entries WHERE service=?", service)
	return err
}
