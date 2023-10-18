import React, {useState, useEffect} from 'react';
import Cookies from 'js-cookie';
import AuthBanner from "./component/AuthBanner";
import Header from './component/Header';
import {TokenDetails} from './component/TokenDetails';
import {DecodedAccessToken} from './component/DecodedAccessToken';
import {ProtectedResource} from './component/ProtectedResource';
import {getCookieValue} from "./utils/cookieUtils";
import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap/dist/js/bootstrap.bundle.min.js';
import config from './config';

function App() {
    const [hasSession, setHasSession] = useState(false);
    const [tokenData, setTokenData] = useState(null);
    const [protectedResourceData, setProtectedResourceData] = useState(null);

    useEffect(() => {
        const sessionName = Cookies.get('session');
        if (sessionName) {
            setHasSession(true);
        } else {
            console.log("Session not found.");
        }

        const loginUrl = `${config.authClientUrl}/tokenData`;
        const handleGetTokenData = async () => {
            try {
                const accessToken = getCookieValue('access_token');
                const headers = {
                    'Authorization': `Bearer ${accessToken}`
                };

                const response = await fetch(loginUrl, {
                    method: 'GET',
                    headers: headers
                });

                if (response.ok) {
                    const tokenData = await response.json();
                    setTokenData(tokenData);
                } else {
                    console.error('Error fetching token data:', response.status);
                }
            } catch (error) {
                console.error('An error occurred while fetching token data:', error);
            }
        };

        handleGetTokenData();
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
        const refreshTokenUrl = `${config.authClientUrl}/refreshToken`;
        const accessToken = getCookieValue('access_token');
        const headers = {
            'Authorization': `Bearer ${accessToken}`
        };

        fetch(refreshTokenUrl, {
            method: 'GET',
            headers: headers
        })
            .then(response => {
                if (response.ok) {
                    return response.json();
                } else {
                    throw new Error('Failed to refresh token');
                }
            })
            .then(data => {
                document.cookie = `access_token=${data.access_token}`;
                window.location.href = window.location.href;
            })
            .catch(error => {
                console.error(error);
            });
    };

    const handleGetProtectedResource = () => {
        const getProtectedResourceUrl = `${config.authClientUrl}/services`;
        const accessToken = getCookieValue('access_token');
        const headers = {
            'Authorization': `Bearer ${accessToken}`
        };

        fetch(getProtectedResourceUrl, {
            method: 'GET',
            headers: headers
        })
            .then(response => {
                if (response.ok) {
                    return response.json();
                } else {
                    throw new Error('Failed to get protected resource');
                }
            })
            .then(data => {
                setProtectedResourceData(data);
            })
            .catch(error => {
                console.error(error);
            });
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
                            <TokenDetails tokenData={tokenData}/>
                            <DecodedAccessToken tokenData={tokenData} />
                            <ProtectedResource resourceDetails={protectedResourceData ? protectedResourceData.services : null}/>
                        </>
                    )}
                </div>
            </div>
        </div>
    );
}

export default App;
