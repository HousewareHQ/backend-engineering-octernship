# Instructions:
## Setting Up BACKEND
Make sure golang environment is setted up.[Setup golang](https://go.dev/doc/install)<br>
## STEP #1:
**Clone Repository**
```
git clone https://github.com/HousewareHQ/houseware---backend-engineering-octernship-thepranays.git
```

## STEP #2:
**INSTALL REQUIRED MODULES**
<br>
Run following commands in terminal
```
cd houseware---backend-engineering-octernship-thepranays/backend

go mod tidy
go mod vendor
```

## STEP #3:
**SETUP DATABASE**
<br>
For simplicity purpose,we are going to use [MongoDB atlas](https://www.mongodb.com/atlas/database)
to create cloud database.
<br>
Step 1:[Create Account](https://account.mongodb.com/account/register)
<br>
Step 2:[Deploy cluster](https://www.mongodb.com/docs/atlas/tutorial/deploy-free-tier-cluster/)
<br>
Step 3:[Setup Network Access](https://www.mongodb.com/docs/atlas/security/add-ip-address-to-list/)
<br>
Step 4: Create database named "houseware" [Setup Database Access](https://www.mongodb.com/docs/atlas/tutorial/create-mongodb-user-for-cluster/)
<br> 
Step 5:[Connecting To Database](https://www.mongodb.com/docs/atlas/tutorial/connect-to-your-cluster/)




## STEP #4:

**CREATE .ENV FILE**

Example of .env 
```
MONGODB_CREDURL = "mongodb+srv://{username}:{password}@cluster0.abc123.mongodb.net/houseware?retryWrites=true&w=majority"

PORT = "8080"
```

## STEP #5:
**Starting Backend-API Server**
```
cd /houseware---backend-engineering-octernship-thepranays/backend/api/server

go run .
```
---
___
## Setting Up FRONTEND

Make sure Node.js envirnoment is setted up.[Setup Node.js](https://nodejs.org/en/download)

**Starting React.js client:**
```
cd /houseware---backend-engineering-octernship-thepranays/backend/web/app

npm install

npm start
``` 

- - -
## Tech Stack:
Gin-gonic (Golang API Framework)
<br>
![gin-gonic logo](https://preview.redd.it/3dto8z3ma7671.png?width=960&crop=smart&auto=webp&v=enabled&s=6b6fa77f1355b4dbeccd2637c5ee2967d92aab58)
React.js
<br>
![reactjs logo](https://www.datocms-assets.com/45470/1631110818-logo-react-js.png)
MongoDB
<br>
![mongodb logo](https://g.foolcdn.com/art/companylogos/square/mdb.png)

___
# More detailed information:
**/backend/api/README.md**:Information on API endpoints,architecture and database design. 

**/backend/web/app/README.md**:Information on Frontend Client.


