package main

const UserSchema = `
CREATE TABLE IF NOT EXISTS users(
	username TEXT NOT NULL PRIMARY KEY,
	password TEXT NOT NULL);
`
