# Hotel reservation

## Resources
### Mongodb driver
Documentation
```
https://www.mongodb.com/docs/drivers/go/current/quick-start/
```

### Use go mod to initialize a Project
```
mkdir go-quickstart
cd go-quickstart
go mod init go-quickstart
```
Add MongoDB as a Dependency
```
go get go.mongodb.org/mongo-driver/mongo
```

## gofiber
install
```
go get github.com/gofiber/fiber/v2
```

## Connect to Mongo
```
mongosh
```

run the commands 
```
use DATABASE-NAME
show collections
```