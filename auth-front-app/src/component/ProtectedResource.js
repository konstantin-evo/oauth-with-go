import React from 'react';

function ProtectedResource({resourceDetails}) {
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

export {ProtectedResource};
