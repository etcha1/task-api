# task-api

Tech Stack

**Language**: Go  
**Database**: PostgreSQL  
**HTTP**: net/http  
**Router**: chi  
**SQL**: pgx  
**Authentication**: JWT  
**Password hashing**: bcrypt  
**Configuration**: environment variables  
**Migration tool**: golang-migrate  
**Testing**: Go's testing package  

## Run project
docker compose up -d

## Run SQL migrations
docker run -v {{ migration dir }}:/migrations --network host migrate/migrate
    -path=/migrations/ -database postgres://127.0.0.1:5432/taskmanager up 2

## Api doc
http://localhost:3000/docs/index.html