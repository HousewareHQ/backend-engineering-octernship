package AppConstant

import "time"

var USER_COLLECTION = "users"
var DB_NAME = "houseware"
var TOKEN_EXPIRY = time.Second * time.Duration(3)        //for jwt token expiryAT format
var REFRESH_TOKEN_EXPIRY = time.Hour * time.Duration(24) //for jwt token expiryAT format

var TOKEN_COOKIE_EXPIRY = 3600              //1hr          //for cookie token expiryAT format
var REFRESH_TOKEN_COOKIE_EXPIRY = 24 * 3600 //for cookie token expiryAT format
