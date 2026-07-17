# task-api

# Run project
docker compose up -d

# Run SQL migrations
docker run -v {{ migration dir }}:/migrations --network host migrate/migrate
    -path=/migrations/ -database postgres://127.0.0.1:5432/taskmanager up 2