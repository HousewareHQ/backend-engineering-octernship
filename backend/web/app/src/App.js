import React from 'react';
import{
  BrowserRouter,

 
}  from 'react-router-dom' 
import './App.css';
import AuthApi from './contexts/AuthApi';
import { Routing } from './routes/routes';

function App() {
  const[auth,setAuth] = React.useState(false); //will be used to determine is session valid and user is loggedin
  const value = React.useMemo(  
    () => ({ auth, setAuth }),
    [auth, setAuth ],
  );

  return (
    <div className="App">
      <AuthApi.Provider value={value}>
        <BrowserRouter>
         <Routing/>
        </BrowserRouter>
        </AuthApi.Provider>
    </div>
  );
}

export default App;
