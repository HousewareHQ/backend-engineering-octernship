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
        <>
            <td>{props.username}  </td>
            <td>{props.org}  </td>
            <td>{props.usertype}  </td>
            <td>{props.createdon}</td>
            <td>{props.updatedon}</td>

            <td><button style={{display:"inline"}} onClick={()=>{handleDeleteUser(props.id)}}>DELETE</button></td>
            
         </>
    )
}