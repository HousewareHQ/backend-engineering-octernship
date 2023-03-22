import {Navigate} from 'react-router-dom'
import AuthApi from '../contexts/AuthApi';
import React from 'react';
import { ValidateSession } from '../helpers/AuthenticateUser';


export const ProtectedPage= ({component:Component })=>{
    const {auth,setAuth} = React.useContext(AuthApi);
    const [authLoading,setAuthLoading] = React.useState(true);
    const [users,setUsers] = React.useState([]);
    React.useEffect(()=>{
            async function fetchRes(){
                const res = await ValidateSession(); //TO CHECK WHETHER SESSION IS VALID OR NOT
                setAuthLoading(false)   //SETTING LOADING FALSE
                setAuth(res[0]);//UPDATING ,AUTH STATE == isSessionValid
                setUsers(res[1]); //UPDATING,users list
               
                /* res[0] -> isValidSession(boolean)
                res[1] -> usersList(array)
                */
            }        
            fetchRes();
    },[setAuth]);

    if(authLoading) return <div>Loading..</div>
    return(
        auth?
        (<Component data={users}/>) //render component if auth success
        :
        <Navigate to="/login" replace={true}/> //redirectReplace if auth fails
        


    );
}