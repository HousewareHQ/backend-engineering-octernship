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
        return (<tr key={user["username"]}>
           <User id={user["ID"]} username={user["username"]} org={user["org"]} usertype={user["usertype"]} createdon={user["createdon"]} updatedon={user["updatedon"]}/>
           </tr>)
    });

  

    return(
        <div className="dashboard-container">
        <h1 style={{display:"inline-block"}}>USERS LIST</h1> 
        <button style={{display:"inline-block",marginLeft:"25px"}} onClick={handleLogout}>Logout</button>
        <br/>
        <div className="users-container">
        <table>
            <thead>
            <tr>
                <th>Name</th>
                <th>Organisation</th>
                <th>User-Type</th>
                <th>Created On</th>
                <th>Updated On</th>
            </tr>
            </thead>
            <tbody>
            {userItems}
            </tbody>
        </table>
       
        </div>
        <CreateUserForm/>
        </div>
    )
}
