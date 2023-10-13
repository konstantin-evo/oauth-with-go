import React, {useState, useEffect} from 'react';
import Cookies from 'js-cookie';
import AuthBanner from "./component/AuthBanner";
import Header from './component/Header';
import {TokenDetails} from './component/TokenDetails';
import {DecodedAccessToken} from './component/DecodedAccessToken';
import {ProtectedResource} from './component/ProtectedResource';
import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap/dist/js/bootstrap.bundle.min.js';
import config from './config';

import {
    dummyResourceDetails,
} from './dummyData';

function App() {
    const [hasSession, setHasSession] = useState(false);

    useEffect(() => {
        const sessionName = Cookies.get('session');
        if (sessionName) {
            setHasSession(true);
        } else {
            console.log("Session not found.");
        }
    }, []);

    const handleLogin = async () => {
        const loginUrl = `${config.authClientUrl}/login`;
        window.location.href = loginUrl;
    };

    const handleLogout = () => {
        Cookies.remove('session');
        setHasSession(false);
    };

    const handleRefreshToken = () => {
        // Implement token refresh logic here
    };

    const handleGetProtectedResource = () => {
        // Implement fetching protected resource here
    };

    return (
        <div className="App">
            <Header
                hasSession={hasSession}
                onLogin={handleLogin}
                onLogout={handleLogout}
                onRefreshToken={handleRefreshToken}
            />
            <AuthBanner hasSession={hasSession} onGetProtectedResource={handleGetProtectedResource}/>
            <div className="container mt-4">
                <div className="row">
                    {hasSession && (
                        <>
                            <TokenDetails/>
                            <DecodedAccessToken/>
                            <ProtectedResource resourceDetails={dummyResourceDetails}/>
                        </>
                    )}
                </div>
            </div>
        </div>
    );
}

export default App;
