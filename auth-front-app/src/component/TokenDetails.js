import React, {useEffect, useState} from "react";
import {getCookieValue} from "../utils/cookieUtils";
import config from "../config";
import {DecodedAccessToken} from "./DecodedAccessToken";

function TokenDetails() {
    const [accessToken, setAccessToken] = useState('');
    const [expiresIn, setExpiresIn] = useState('');
    const [refreshToken, setRefreshToken] = useState('');
    const [scope, setScope] = useState('');
    const [tokenType, setTokenType] = useState('');
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        const loginUrl = `${config.authClientUrl}/tokenData`;
        const handleGetTokenDetails = async () => {
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
                    setAccessToken(tokenData.access_token);
                    setExpiresIn(tokenData.expires_in);
                    setRefreshToken(tokenData.refresh_token);
                    setScope(tokenData.scope);
                    setTokenType(tokenData.token_type);

                    DecodedAccessToken({ accessToken: tokenData.access_token });
                } else {
                    console.error('Error fetching token data:', response.status);
                }
            } catch (error) {
                console.error('An error occurred while fetching token data:', error);
            } finally {
                setIsLoading(false);
            }
        };

        handleGetTokenDetails();
    }, []);

    if (isLoading) {
        return (
            <div className="col-md-4 mb-4">
                <div className="card-header">
                    Token Details
                </div>
                <div className="card h-100">
                    <div className="card-body overflow-auto">Loading token details...</div>
                </div>
            </div>
        );
    }

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


export {TokenDetails};
