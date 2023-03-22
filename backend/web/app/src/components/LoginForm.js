import React,{useState} from 'react'
import {useNavigate } from 'react-router-dom';
import AuthApi from '../contexts/AuthApi';
import { AuthenticateUser } from '../helpers/AuthenticateUser';
export const LoginForm = (props)=>{
    const [username,setUsername] =useState('');
    const [pass,setPass] = useState('');
    const {auth,setAuth} = React.useContext(AuthApi) //To Auth API access context
    const navigate = useNavigate();//To redirect on login

    const handleLogin= async (e)=>  {
        e.preventDefault();//to avoid reloading page
        const isLoggedin =await AuthenticateUser(username,pass)
        setAuth(isLoggedin)
       
        
        
    }
    React.useEffect(()=>{ //When auth changes execute code inside 
        if(auth) navigate("/dashboard",{replace:true})  // if authenticated then redirectAndReplace to dashboard
        // eslint-disable-next-line react-hooks/exhaustive-deps
    },[auth]);

    return (
        <div className="div-login-container">
            <form className="login-form" onSubmit={handleLogin}>
                <div>
                    <label htmlFor="username">Username</label>
                    <br/>
                    <input value={username} type="text" onChange={(e)=>setUsername(e.target.value)} name="username" id="username"/>
                </div>
                <div>
                    <label htmlFor="password">Password</label>
                    <br/>

                    <input value={pass} type="password" onChange={(e)=>setPass(e.target.value)} name="password" id="password"/>    
                </div>
                <button type="submit">Login</button>
            </form>
        </div>
    )
}
