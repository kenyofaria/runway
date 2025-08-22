import React from 'react';
import './AppList.css';

const AppList = ({ apps, isLoading, error, onSelectApp, onFetchApps }) => {
    return (
        <div className="app-list-container">
            <h2>App List</h2>
            <button onClick={onFetchApps} disabled={isLoading}>
                {isLoading ? 'Loading...' : 'Get Apps'}
            </button>

            {error && <p className="error-message">Error: {error}</p>}

            {apps.length > 0 && (
                <div className="app-cards">
                    {apps.map((app) => (
                        <div key={app.id} className="app-card">
                            <img src={app.artwork_url} alt={app.name} className="app-image" />
                            <div className="app-info">
                                <h3>{app.name}</h3>
                                <p>by {app.artistName}</p>
                                <p>Released: {app.releaseDate}</p>
                            </div>
                            <button onClick={() => onSelectApp(app.id, app.name)}>
                                View Reviews
                            </button>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
};

export default AppList;