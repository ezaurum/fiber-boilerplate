# GoFiber Docker Boilerplate
## Casbin
## GORM

## TODO
- [ ] Excel
- [ ] Remember me with cookie
- [ ] Login with cookie
- [ ] Redis session
- [ ] Session slide
- [ ] CORS
- [x] Websocket
- [ ] OpenAPI
- [x] ID generator
- [x] Default database model


## Development

### Start the application 


```bash
go run app.go
```

### Use local container

```
# Shows all commands
make help

# Clean packages
make clean-packages

# Generate go.mod & go.sum files
make requirements

# Generate docker image
make build

# Generate docker image with no cache
make build-no-cache

# Run the projec in a local container
make up

# Run local container in background
make up-silent

# Run local container in background with prefork
make up-silent-prefork

# Stop container
make stop

# Start container
make start
```

## Production

```bash
docker build -t gofiber .
docker run -d -p 3000:3000 gofiber ./app -prod
```

Go to http://localhost:3000:


![Go Fiber Docker Boilerplate](./go_fiber_boilerplate.gif)
