import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, useNavigate, useParams } from 'react-router-dom';
import AppList from './components/AppList';
import AppReviews from './components/AppReviews';
import './App.css';

// Component to handle the app reviews route
const AppReviewsRoute = ({ apps, setApps, isLoadingApps, setIsLoadingApps, errorApps, setErrorApps }) => {
  const { appId } = useParams();
  const navigate = useNavigate();
  const [selectedHours, setSelectedHours] = useState(24);

  // Find the app name from the apps array
  const selectedApp = apps.find(app => app.id === appId);
  const appName = selectedApp ? selectedApp.name : 'Unknown App';

  const handleBack = () => {
    navigate('/');
  };

  // If we don't have the app data and we have an appId, we might need to fetch apps
  useEffect(() => {
    if (appId && apps.length === 0 && !isLoadingApps) {
      // Optionally fetch apps if not already loaded
      // This ensures we have app names even when navigating directly to a review URL
    }
  }, [appId, apps.length, isLoadingApps]);

  return (
      <AppReviews
          appId={appId}
          appName={appName}
          onBack={handleBack}
          selectedHours={selectedHours}
          setSelectedHours={setSelectedHours}
      />
  );
};

// Component to handle the app list route
const AppListRoute = ({ apps, isLoadingApps, errorApps, onFetchApps }) => {
  const navigate = useNavigate();

  const handleSelectApp = (appId, appName) => {
    navigate(`/app/${appId}`);
  };

  return (
      <AppList
          apps={apps}
          isLoading={isLoadingApps}
          error={errorApps}
          onSelectApp={handleSelectApp}
          onFetchApps={onFetchApps}
      />
  );
};

const App = () => {
  const [apps, setApps] = useState([]);
  const [isLoadingApps, setIsLoadingApps] = useState(false);
  const [errorApps, setErrorApps] = useState(null);

  const fetchApps = async () => {
    setIsLoadingApps(true);
    setErrorApps(null);
    try {
      const response = await fetch(`${process.env.REACT_APP_API_URL}/app/list`);
      if (!response.ok) {
        throw new Error('Network response was not ok');
      }
      const data = await response.json();
      setApps(data);
    } catch (err) {
      setErrorApps(err.message);
    } finally {
      setIsLoadingApps(false);
    }
  };

  return (
      <Router>
        <div className="App">
          <header className="App-header">
            <h1>App Showcase</h1>
          </header>
          <main>
            <Routes>
              <Route
                  path="/"
                  element={
                    <AppListRoute
                        apps={apps}
                        isLoadingApps={isLoadingApps}
                        errorApps={errorApps}
                        onFetchApps={fetchApps}
                    />
                  }
              />
              <Route
                  path="/app/:appId"
                  element={
                    <AppReviewsRoute
                        apps={apps}
                        setApps={setApps}
                        isLoadingApps={isLoadingApps}
                        setIsLoadingApps={setIsLoadingApps}
                        errorApps={errorApps}
                        setErrorApps={setErrorApps}
                    />
                  }
              />
            </Routes>
          </main>
        </div>
      </Router>
  );
};

export default App;