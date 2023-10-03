import React from 'react';

function Header({ hasSession, onLogin, onLogout, onRefreshToken }) {
    return (
        <nav className="navbar navbar-expand-lg navbar-light bg-light">
                <div className="collapse navbar-collapse">
                    <ul className="navbar-nav ml-auto">
                        {hasSession ? (
                            <>
                                <li className="nav-item mr-1">
                                    <button className="btn btn-secondary" onClick={onRefreshToken}>Refresh Token</button>
                                </li>
                                <li className="nav-item">
                                    <button className="btn btn-secondary" onClick={onLogout}>Logout</button>
                                </li>
                            </>
                        ) : (
                            <li className="nav-item">
                                <button className="btn btn-secondary" onClick={onLogin}>Login</button>
                            </li>
                        )}
                    </ul>
                </div>
        </nav>
    );
}

export default Header;
