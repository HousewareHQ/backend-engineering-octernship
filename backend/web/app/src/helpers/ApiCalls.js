import { API_GETUSERS_ENDPOINT,API_CREATE_USER_ENDPOINT, API_URL, API_DELETE_USER_ENDPOINT } from "../constants/Constants";


export async function GetAllUsers(){
    //fetch all users belonging to one organization 
    const result =await fetch(API_URL+API_GETUSERS_ENDPOINT,{
        method:"GET",
        headers:{
            "Content-Type":"application/json"
        },
      
        mode:"cors",
        credentials:"include",


     });
     return await result.json();
 }

 export async function CreateUser(username,pass,usertype){
    //Create user if admin
    const result = await fetch(API_URL+API_CREATE_USER_ENDPOINT,{
        method:"POST",
        headers:{
            "Content-Type":"application/json"
        },
        body:JSON.stringify({
            "username":username,
            "password":pass,
            "usertype":usertype,            
        }),
        mode:"cors",
        credentials:"include",
    });
    return [result.status,result.statusText]
 }
 export async function DeleteUser(uid){
    //delete user if admin
    const result = await fetch(API_URL+API_DELETE_USER_ENDPOINT+`/${uid}`,{
        method:"DELETE",
       
        mode:"cors",
        credentials:"include",
    });
    return [result.status,result.statusText]
 }