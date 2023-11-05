import React from 'react';
import lockOpenIcon from '../img/lock-open.svg';
import lockClosedIcon from '../img/lock-closed.svg';

function AuthBanner({ hasSession, onGetProtectedResource }) {
    return (
        <div className="jumbotron text-center">
            <img className="p-3" src={hasSession ? lockOpenIcon : lockClosedIcon} alt={hasSession ? 'Open Lock Icon' : 'Closed Lock Icon'} style={{ width: '18%' }}/>
            <h1 className="display-4">Welcome to the Auth App</h1>
            <p className="lead">Explore the Keycloak based authentication solution for educational purposes.</p>

            {hasSession && (
                <button id="get-protected-resource" className="btn btn-primary btn-lg mt-4" onClick={onGetProtectedResource}>
                    Get protected resource
                </button>
            )}
        </div>
    );
}

export default AuthBanner;
