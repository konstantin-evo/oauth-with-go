import React from 'react';

function ProtectedResource({resourceDetails}) {
    return (
        <div className="col-md-4 mb-4">
            <div className="card h-100">
                <div className="card-header">
                    Protected Resource
                </div>
                <div className="card-body overflow-auto">
                    <h6 className="card-subtitle mb-2 text-muted">
                        {resourceDetails ? 'Services received from a protected server' : ''}
                    </h6>
                    {resourceDetails && Object.keys(resourceDetails).map((key) => (
                        <div key={key}>
                            <p className="card-text">{resourceDetails[key]}</p>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
}

export {ProtectedResource};
