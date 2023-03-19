import React from "react";
import { DeleteUser } from "../helpers/ApiCalls";


export const User= (props)=>{
    const handleDeleteUser= async (uid)=>{
        const [statusCode,statusText] = await DeleteUser(uid);
        alert(`${statusCode}:${statusText}`);
        //Refreshing whole page for simplicity purpose to see updated userlist
        window.location.reload(false);
    }
    return(
        <div className="user-container">
            <h3 >{props.id} {props.username}  {props.org}  {props.usertype}  <button style={{display:"inline"}} onClick={()=>{handleDeleteUser(props.id)}}>DELETE</button></h3>
            

        </div>
    )
}