import React, { useState, useEffect } from 'react';
import Cookies from 'js-cookie';
import Header from './component/Header';
import {TokenDetails} from './component/ProtectedResource';
import {DecodedAccessToken} from './component/ProtectedResource';
import {ProtectedResource} from './component/ProtectedResource';
import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap/dist/js/bootstrap.bundle.min.js';

import {
  dummyAccessToken,
  dummyExpiresIn,
  dummyRefreshToken,
  dummyScope,
  dummyTokenType,
  dummyDecodedToken,
  dummyResourceDetails,
} from './dummyData';
import AuthBanner from "./component/AuthBanner";

function App() {
  const [hasSession, setHasSession] = useState(false);

  useEffect(() => {
    const sessionName = Cookies.get('session-name');
    if (sessionName) {
      setHasSession(true);
    }
  }, []);

  const handleLogin = () => {
    // Implement login here
  };


  const handleLogout = () => {
    Cookies.remove('session-name');
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
        <AuthBanner hasSession={hasSession} onGetProtectedResource={handleGetProtectedResource} />
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
                  <DecodedAccessToken decodedToken={dummyDecodedToken} />
                  <ProtectedResource resourceDetails={dummyResourceDetails} />
                </>
            )}
          </div>
        </div>
      </div>
  );
}

export default App;
