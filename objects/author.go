package objects

import (
	"fmt"
	"os"
	"time"
)

const TimeFormat = "-0700"

func GitTimestamp(t time.Time) string {
	return fmt.Sprintf("%d %s", t.Unix(), t.Format(TimeFormat))
}

func Author(user, email string) string {
	return fmt.Sprintf("%s <%s>", user, email)
}

func AuthorAtHost(user, host string) string {
	email := user + "@" + host
	return Author(user, email)
}

func DefaultAuthor() string {
	user := os.Getenv("USER")
	if user == "" {
		user = "Anonymous Coward"
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	return AuthorAtHost(user, host)
}

func AuthorLine(author string, t time.Time) string {
	return fmt.Sprintf("%s %s", author, GitTimestamp(t))
}
