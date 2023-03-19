import {API_URL,API_GETUSERS_ENDPOINT, API_LOGIN_ENDPOINT, API_LOGOUT_ENDPOINT} from "../constants/Constants"

/*returns true if session is valid*/
export  async function ValidateSession(){
    //Make a api call to validate token and check if session is still valid or not
   const res =await fetch(API_URL+API_GETUSERS_ENDPOINT,{
        method:"GET",
        headers:{
            "Content-Type":"application/json"
        },
        mode:"cors",
        credentials:"include",


     });
     const users = await res.json();
     //returns true if session valid and also returns all users in current user's org
    return [res.status===200 ,users];
}

export async function AuthenticateUser  (username,password) {
        const res=await fetch(API_URL+API_LOGIN_ENDPOINT,{
        method:"POST",
        headers:{
            "Content-Type":"application/json"
        },
        credentials:"include",
        
        mode:"cors",
        body:JSON.stringify({
            "username":username,
            "password":password,
        }),

    });
    return res.status===200

}

export async function DestroySession(){ //destroy's session by expiring cookies i.e. logging out
    const res = await fetch(API_URL+API_LOGOUT_ENDPOINT,{
        method:"POST",
        credentials:"include",
        mode:"cors",


    });
    return res.status ===200

}