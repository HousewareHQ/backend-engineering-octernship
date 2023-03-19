import React from "react"
import { useNavigate } from "react-router-dom"
import { DestroySession } from "../helpers/AuthenticateUser"
import { CreateUserForm } from "./CreateUserForm";
import { User } from "./User";


export const Dashbaord= (props)=>{
    const navigate = useNavigate(); //navigator
    const handleLogout = async (e)=>{
        e.preventDefault();//to avoid reloading page

        const res = await DestroySession(); //Destroy session and delete token cookies
        if(res){
            navigate("/login",{replace:true}); //redirect to login (WITH REPLACE)
        }
    }

    //Returns user as list item
    const userItems = props.data.map((user)=>{
        return (<li key={user["username"]}>
           <User id={user["ID"]} username={user["username"]} org={user["org"]} usertype={user["usertype"]}/>
           </li>)
    });

  

    return(
        <div className="dashboard-container">
        <h1 style={{display:"inline-block"}}>USERS LIST</h1> 
        <button style={{display:"inline-block"}} onClick={handleLogout}>Logout</button>
        <br/>
        <div className="users-container">
        <ul className="users-list">
            {/* render all users in unordered list */}
            {userItems} 
        </ul>
        </div>
        <CreateUserForm/>
        </div>
    )
}
