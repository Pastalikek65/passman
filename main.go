package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/Pastalikek65/passman/store"
	"github.com/Pastalikek65/passman/tui"
)

var version = "0.1.0"

func defaultDB() string {
	if xdg := os.Getenv("XDG_CACHE_HOME"); xdg != "" {
		return filepath.Join(xdg, "passman", "pass.db")
	}
	home, _ := os.UserHomeDir()
	if home == "" { home = os.Getenv("HOME") }
	if home == "" { home = "/tmp" }
	return filepath.Join(home, ".cache", "passman", "pass.db")
}

func getPassphrase() string {
	if p := os.Getenv("PASSMAN_MASTER"); p != "" { return p }
	fmt.Fprint(os.Stderr, "Enter master passphrase: ")
	var p string
	fmt.Scanln(&p)
	return p
}

func main() {
	root := &cobra.Command{Use: "passman", Short: "Age-encrypted password manager for Termux", Version: version}
	root.AddCommand(&cobra.Command{
		Use: "add [service] [username]", Args: cobra.ExactArgs(2), Short: "add entry",
		RunE: func(cmd *cobra.Command, args []string) error {
			pass := getPassphrase()
			if pass == "" { return fmt.Errorf("passphrase required") }
			s, _ := store.Open(defaultDB(), pass)
			defer s.Close()
			fmt.Print("Password: ")
			var pw string
			fmt.Scanln(&pw)
			id, err := s.Add(args[0], args[1], pw)
			if err != nil { return err }
			fmt.Println(id)
			return nil
		},
	})
	root.AddCommand(&cobra.Command{
		Use: "get [service]", Args: cobra.ExactArgs(1), Short: "get password",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, _ := store.Open(defaultDB(), getPassphrase())
			defer s.Close()
			e, err := s.Get(args[0])
			if err != nil { return err }
			fmt.Println(e.Password)
			return nil
		},
	})
	root.AddCommand(&cobra.Command{
		Use: "list", Short: "list services",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, _ := store.Open(defaultDB(), getPassphrase())
			defer s.Close()
			ents,_:=s.List()
			for _,e:=range ents { fmt.Printf("%s\t%s\n", e.Service, e.Username) }
			return nil
		},
	})
	root.AddCommand(&cobra.Command{
		Use: "tui", Short: "open TUI",
		RunE: func(cmd *cobra.Command, args []string) error {
			s,_:=store.Open(defaultDB(), getPassphrase())
			defer s.Close()
			m:=tui.New(s)
			p:=tea.NewProgram(m, tea.WithAltScreen())
			_,err:=p.Run()
			return err
		},
	})
	root.AddCommand(&cobra.Command{
		Use: "delete [service]", Args: cobra.ExactArgs(1), Short: "delete entry",
		RunE: func(cmd *cobra.Command, args []string) error {
			s,_:=store.Open(defaultDB(), getPassphrase())
			defer s.Close()
			return s.Delete(args[0])
		},
	})
	if err:=root.Execute(); err!=nil { os.Exit(1) }
}
