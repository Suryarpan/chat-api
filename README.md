# Chat Api (Golang)

This is an experimental chat api built in go. Following features are available in this.

1. [x] New user registration
1. [x] User Authentication using JWT
1. [ ] CRUD on user
1. [ ] Blocklist for users
1. [ ] CRUD on messages
1. [ ] CRUD on User groups
1. [ ] CRUD on group messages
1. [ ] Admin related operations

## Tech Stack

The server is built with [go-chi](https://github.com/go-chi/chi). For backend
storage as of now Postgres is used. There are plans to introduce NoSQL DBs for
message storage.
