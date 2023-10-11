import React, {useState, useEffect} from 'react';
import Cookies from 'js-cookie';
import AuthBanner from "./component/AuthBanner";
import Header from './component/Header';
import {TokenDetails} from './component/ProtectedResource';
import {DecodedAccessToken} from './component/ProtectedResource';
import {ProtectedResource} from './component/ProtectedResource';
import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap/dist/js/bootstrap.bundle.min.js';
import config from './config';

import {
    dummyAccessToken,
    dummyExpiresIn,
    dummyRefreshToken,
    dummyScope,
    dummyTokenType,
    dummyDecodedToken,
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
                            <TokenDetails
                                accessToken={dummyAccessToken}
                                expiresIn={dummyExpiresIn}
                                refreshToken={dummyRefreshToken}
                                scope={dummyScope}
                                tokenType={dummyTokenType}
                            />
                            <DecodedAccessToken decodedToken={dummyDecodedToken}/>
                            <ProtectedResource resourceDetails={dummyResourceDetails}/>
                        </>
                    )}
                </div>
            </div>
        </div>
    );
}

export default App;
