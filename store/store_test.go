package store

import (
	"testing"
)

func TestAddAndGet(t *testing.T) {
	s, _ := Open(t.TempDir()+"/pass.db", "testpass123")
	defer s.Close()
	id, err := s.Add("github", "alice", "s3cr3t")
	if err != nil { t.Fatal(err) }
	if id == 0 { t.Fatal("no id") }
	e, err := s.Get("github")
	if err != nil { t.Fatal(err) }
	if e.Password != "s3cr3t" { t.Fatalf("want s3cr3t got %q", e.Password) }
	if e.Username != "alice" { t.Fatalf("want alice got %q", e.Username) }
}

func TestList(t *testing.T) {
	s, _ := Open(t.TempDir()+"/pass.db", "testpass123")
	defer s.Close()
	s.Add("a", "u1", "p1")
	s.Add("b", "u2", "p2")
	list, _ := s.List()
	if len(list) != 2 { t.Fatalf("want 2 got %d", len(list)) }
	if list[0].Service != "a" { t.Fatalf("want a got %s", list[0].Service) }
}

func TestDelete(t *testing.T) {
	s, _ := Open(t.TempDir()+"/pass.db", "testpass123")
	defer s.Close()
	s.Add("x", "u", "p")
	if err := s.Delete("x"); err != nil { t.Fatal(err) }
	if _, err := s.Get("x"); err == nil { t.Fatal("want error after delete") }
}

func TestWrongPassphrase(t *testing.T) {
	s, _ := Open(t.TempDir()+"/pass.db", "correct")
	defer s.Close()
	s.Add("svc", "u", "secret")
	s2, _ := Open(s.path, "wrong")
	defer s2.Close()
	// Get with wrong passphrase should fail to decrypt
	if _, err := s2.Get("svc"); err == nil { t.Fatal("want decrypt error with wrong pass") }
}
