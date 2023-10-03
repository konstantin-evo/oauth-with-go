import React from 'react';

function DecodedAccessToken({ decodedToken }) {
    return (
        <div className="col-md-4 mb-4">
            <div className="card h-100">
                <div className="card-header">
                    Decoded Access Token
                </div>
                <div className="card-body overflow-auto">
                    <h6 className="card-subtitle mb-2 text-muted">acr</h6>
                    <p className="card-text">{decodedToken.acr}</p>
                    <h6 className="card-subtitle mb-2 text-muted">allowed-origins</h6>
                    <p className="card-text">{decodedToken['allowed-origins']}</p>
                    <h6 className="card-subtitle mb-2 text-muted">aud</h6>
                    <p className="card-text">{decodedToken.aud}</p>
                    <h6 className="card-subtitle mb-2 text-muted">auth_time</h6>
                    <p className="card-text">{decodedToken.auth_time}</p>
                    <h6 className="card-subtitle mb-2 text-muted">azp</h6>
                    <p className="card-text">{decodedToken.azp}</p>
                    <h6 className="card-subtitle mb-2 text-muted">email_verified</h6>
                    <p className="card-text">{decodedToken.email_verified}</p>
                    <h6 className="card-subtitle mb-2 text-muted">exp</h6>
                    <p className="card-text">{decodedToken.exp}</p>
                    <h6 className="card-subtitle mb-2 text-muted">iat</h6>
                    <p className="card-text">{decodedToken.iat}</p>
                    <h6 className="card-subtitle mb-2 text-muted">iss</h6>
                    <p className="card-text">{decodedToken.iss}</p>
                    <h6 className="card-subtitle mb-2 text-muted">jti</h6>
                    <p className="card-text">{decodedToken.jti}</p>
                    <h6 className="card-subtitle mb-2 text-muted">scope</h6>
                    <p className="card-text">{decodedToken.scope}</p>
                    <h6 className="card-subtitle mb-2 text-muted">session_state</h6>
                    <p className="card-text">{decodedToken.session_state}</p>
                    <h6 className="card-subtitle mb-2 text-muted">sid</h6>
                    <p className="card-text">{decodedToken.sid}</p>
                    <h6 className="card-subtitle mb-2 text-muted">sub</h6>
                    <p className="card-text">{decodedToken.sub}</p>
                    <h6 className="card-subtitle mb-2 text-muted">typ</h6>
                    <p className="card-text">{decodedToken.typ}</p>
                </div>
            </div>
        </div>
    );
}

function ProtectedResource({ resourceDetails }) {
    return (
        <div className="col-md-4 mb-4">
            <div className="card h-100">
                <div className="card-header">
                    Protected Resource
                </div>
                <div className="card-body">
                    <p className="card-text">{resourceDetails}</p>
                </div>
            </div>
        </div>
    );
}

function TokenDetails({ accessToken, expiresIn, refreshToken, scope, tokenType }) {
    return (
        <div className="col-md-4 mb-4">
            <div className="card h-100">
                <div className="card-header">
                    Token Details
                </div>
                <div className="card-body overflow-auto">
                    <h6 className="card-subtitle mb-2 text-muted">AccessToken</h6>
                    <p className="card-text">{accessToken}</p>

                    <h6 className="card-subtitle mb-2 text-muted">ExpiresIn</h6>
                    <p className="card-text">{expiresIn}</p>

                    <h6 className="card-subtitle mb-2 text-muted">RefreshToken</h6>
                    <p className="card-text">{refreshToken}</p>

                    <h6 className="card-subtitle mb-2 text-muted">Scope</h6>
                    <p className="card-text">{scope}</p>

                    <h6 className="card-subtitle mb-2 text-muted">TokenType</h6>
                    <p className="card-text">{tokenType}</p>
                </div>
            </div>
        </div>
    );
}

export { DecodedAccessToken, ProtectedResource, TokenDetails };
