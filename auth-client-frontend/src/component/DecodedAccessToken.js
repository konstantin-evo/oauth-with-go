import React from "react";

function DecodedAccessToken({tokenData}) {

    if (!tokenData) {
        return (
            <div className="col-md-4 mb-4">
                <div className="card h-100">
                    <div className="card-header">
                        Decoded Access Token
                    </div>
                    <div className="card-body overflow-auto">Loading token details...</div>
                </div>
            </div>
        );
    }

    const decodedToken = decodeAccessToken(tokenData.access_token);

    function decodeAccessToken(accessToken) {
        try {
            const payloadBase64 = accessToken.split('.')[1];
            const decodedPayload = JSON.parse(atob(payloadBase64));
            return decodedPayload;
        } catch (error) {
            console.error('Error decoding AccessToken:', error);
            return null;
        }
    }

    if (!decodedToken) {
        return (
            <div className="col-md-4 mb-4">
                <div className="card h-100">
                    <div className="card-header">
                        Decoded Access Token
                    </div>
                    <div className="card-body overflow-auto">Error decoding token.</div>
                </div>
            </div>
        );
    }

    return (
        <div className="col-md-4 mb-4">
            <div className="card h-100">
                <div className="card-header">
                    Decoded Access Token
                </div>
                <div className="card-body overflow-auto">
                    {Object.keys(decodedToken).map((key) => (
                        <>
                            <h6 className="card-subtitle mb-2 text-muted">{key}</h6>
                            <p className="card-text">{decodedToken[key]}</p>
                        </>
                    ))}
                </div>
            </div>
        </div>
    );
}

export {DecodedAccessToken};