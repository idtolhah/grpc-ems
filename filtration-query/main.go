package main

import (
	"database/sql"
	"time"
)

var (
	timeout   = 10 * time.Second
	db_client *sql.DB
)
