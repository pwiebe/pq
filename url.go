package pq

import (
	"fmt"
	nurl "net/url"
	//"sort"
	"strings"
)

// ParseURL no longer needs to be used by clients of this library since supplying a URL as a
// connection string to sql.Open() is now supported:
//
//	sql.Open("postgres", "postgres://bob:secret@1.2.3.4:5432/mydb?sslmode=verify-full")
//
// It remains exported here for backwards-compatibility.
//
// ParseURL converts a url to a connection string for driver.Open.
// Example:
//
//	"postgres://bob:secret@1.2.3.4:5432/mydb?sslmode=verify-full"
//
// converts to:
//
//	"user=bob password=secret host=1.2.3.4 port=5432 dbname=mydb sslmode=verify-full"
//
// A minimal example:
//
//	"postgres://"
//
// This will be blank, causing driver.Open to use all of the defaults

func ParseURL2Map(url string) (map[string]string, error) {
	m := make(map[string]string)
	u, err := nurl.Parse(url)
	if err != nil {
		return m, err
	}

	if u.Scheme != "postgres" {
		return m, fmt.Errorf("invalid connection protocol: %s", u.Scheme)
	}

	if u.User != nil {
		v := u.User.Username()
		m["user"] = v

		v, _ = u.User.Password()
		m["password"] = v
	}

	i := strings.Index(u.Host, ":")
	if i < 0 {
		m["host"] = u.Host
	} else {
		m["host"] = u.Host[:i]
		m["port"] = u.Host[i+1:]
	}

	if u.Path != "" {
		m["dbname"] = u.Path[1:]
	}

	q := u.Query()
	for k := range q {
		m[k] = q.Get(k)
	}

	return m, err
}

func ParsedMap2String(config map[string]string) string {
	var kvs []string
	accrue := func(k, v string) {
		if v != "" {
			kvs = append(kvs, k+"="+v)
		}
	}
	for k, v := range config {
		accrue(k, v)
	}
	return strings.Join(kvs, " ")
}

func ParseURL(url string) (string, error) {
	m, err := ParseURL2Map(url)
	if err != nil {
		return "", err
	}
	return ParsedMap2String(m), nil
}
