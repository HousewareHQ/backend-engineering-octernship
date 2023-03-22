import React, { useState } from "react";
import { CreateUser } from "../helpers/ApiCalls";


export const CreateUserForm= ()=>{
    const[username,setUsername] = useState("");
    const[pass,setPass] = useState("");
    const[usertype,setUsertype] = useState("");
    const [statusText,setStatusText] = useState("");  
    const [statusCode,setStatusCode] = useState();

    const handleCreateUserButton = async (e)=>{
        e.preventDefault();//to avoid reloading page
        const [code,msg]=await CreateUser(username,pass,usertype); //returns status code and statusText
        setStatusCode(code);
        setStatusText(":"+msg);
      

        if(code===200){
            /*
                Instead of listening to changes to userlist and then re-rendering only effecting components,
                We will just refresh the page after showing alert for simplicity purpose
            */
           alert(`${code}:${msg}`);
           window.location.reload(false);

        }

    }
    
   
    return(
       <div className="create-user-container">
                
            <form className="create-user-form">
                <label htmlFor="username">USERNAME:</label>
                <input value={username} type="text" onChange={(e)=>{setUsername(e.target.value)}} name="username" id="username"/>
                <label htmlFor="password">PASSWORD:</label>
                <input value={pass} type="password" onChange={(e)=>{setPass(e.target.value)}} name="password" id="password"/>
                <label htmlFor="usertype">USERTYPE:</label>
                <select value={usertype} name="usertype" id="usertype" onChange={(e)=>{setUsertype(e.target.value)}}>
                    <option value="" disabled>Select Usertype</option>
                    <option value="USER">User</option>
                    <option value="ADMIN">Admin</option>
    
                </select>      
                <button id="button-createuser" onClick={handleCreateUserButton} type="submit">Create User</button>      

            </form>
            <div className="create-complete-container">
                <h1>{statusCode} {statusText}</h1>
            </div>
        </div>
       
    )
}