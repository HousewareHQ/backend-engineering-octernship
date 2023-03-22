import React from 'react';
import{
    Route,  
    Routes,
    
}  from 'react-router-dom';
import { Dashbaord } from '../components/Dashboard';
import {LoginForm} from '../components/LoginForm';
import { ProtectedPage } from '../components/ProtectedPage';





export const Routing = ()=>{
 
    return (
        <Routes>
            <Route path="/" element={<ProtectedPage component={Dashbaord} />}/>
            <Route path="/login" element={<LoginForm/>}/>
            <Route path="/dashboard" element={<ProtectedPage component={Dashbaord} />}/> 
        </Routes>
    )
}
