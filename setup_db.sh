#!/bin/bash
sqlite3 urlshortener.db <<EOF
CREATE TABLE IF NOT EXISTS urls (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    key TEXT UNIQUE NOT NULL,
    original_url TEXT NOT NULL
);
EOF
echo "Database urlshortener.db created with table 'urls'."
