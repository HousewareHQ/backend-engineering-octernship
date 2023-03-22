# To start testing API:
Step 1:Start server
```
go run main.go
```
>Visit http://localhost:8080/api/v1/swagger/index.html

### **Way 1**:
Step1 -Add below line to your .env file
```
MONGODB_CREDURL ="mongodb+srv://publicuser:publicuser@cluster0.gdwf4t5.mongodb.net/houseware?retryWrites=true&w=majority"
```

Step 2-<br>
Use below credentials to access API for testing <br>
**Admin credentials**<br>
1.'user1','abcde'<br>
-ORCN(org1)<br>
2.'user2','12345'<br>
-DreamFist(org2)

# **OR**

### **WAY 2**:
Step 1:
Create and setup "houseware" MongoDB database in your own cluster
i.e. using your MongoDB Account.
<br>
<br>
Step 2:
Further,To start API testing create a super user (first-user).
Use this user to create-other and access API.

**How to create super user?**

Make a POST request with appropriate body.
As following:

**POST**:http://localhost:8080/api/v1/super/create-user/

**BODY**:(Note:usertype must be admin only)
```
{
    "username":"example",
    "password":"examplepass",
    "usertype":"ADMIN",
    "org":"exampleorg"

}
```
- - -
#  API Documentation
### OpenAPI/Swagger Generated Documentation

Step 1:Start Server
```
cd backend/api/server
go run main.go
```
Step 2:
>Visit:http://localhost:8080/api/v1/swagger/index.html

# API Architecture
![api architecture](https://i.ibb.co/ncTLk3J/authbackend-drawio.png)


# 
