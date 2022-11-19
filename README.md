# go-todo

A simple To Do list written in Go.

Originally this was written as a CLI application but I have extended it to also include a web API.

The primary focus is now the backend web API using chi at `cli/web_server`, the example CLI application can be found at `cmd/cli`

To utilise the the backend I've put together a simple frontend using Svelte.

## Running

The full web application utilising the backend and frontend can be run using:
```bash
docker compose up -d
```
It is then accessed at http://localhost:8080

## Ideas for improvements
- Add the ability to have multiple task lists, such as one for work and one for personal.
- Add documentation for the API using swagger.