# Golang API build

## Setup & Running Code
- First clone repo in your PC / System.
- install go package dependencies in your pc, if needed.
- Run : go run main.go  : to start code in your PC
- Test : go test -v  : to test unit testing functions
- Database and collection in DB, will auto created at first successful run of project
- API: You can test API's by opening, postman json file in your postman, and then after running login API please set token in API's, else api will throw error. As CRUD operation API's are secured using JWT token

## API

- http://localhost:9000/user : PUT : its used to create user 
- http://localhost:9000/user : POST : its used to login user to application and get token
- ALL ROUTES BELOW THIS ARE PROTECTED AND NEED TOKEN, WHICH IS SENT IN RESPONSE OF LOGIN
- http://localhost:9000/update/:username : POST : its used to update user token
- http://localhost:9000/user/:username : GET : its used to get specified user in route
- http://localhost:9000/delete : POST : its used to delete user 
- http://localhost:9000/users : GET : its used to get all users from application
- http://localhost:9000/filterbycoordinates : POST : its used to get users to given coordinates sent via request body