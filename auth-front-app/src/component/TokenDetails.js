import React from "react";

function TokenDetails({tokenData}) {

    if (!tokenData) {
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
                    <p className="card-text">{tokenData.access_token}</p>

                    <h6 className="card-subtitle mb-2 text-muted">ExpiresIn</h6>
                    <p className="card-text">{tokenData.expires_in}</p>

                    <h6 className="card-subtitle mb-2 text-muted">RefreshToken</h6>
                    <p className="card-text">{tokenData.refresh_token}</p>

                    <h6 className="card-subtitle mb-2 text-muted">Scope</h6>
                    <p className="card-text">{tokenData.scope}</p>

                    <h6 className="card-subtitle mb-2 text-muted">TokenType</h6>
                    <p className="card-text">{tokenData.token_type}</p>
                </div>
            </div>
        </div>
    );
}

export {TokenDetails};
