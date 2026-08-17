import "./styles/index.css"
import { useEffect, useState } from "react";
import Authentication from "./components/Authentication.jsx";
import { BrowserRouter as Router, Routes, Route, Navigate } from "react-router-dom"
import Events from "./pages/Events.jsx"
import CreateEvent from "./pages/CreateEvent.jsx"
import Admin from "./pages/Admin.jsx"

function App() {
  const [user, setUser] = useState(null);
  const [loadingUser, setLoadingUser] = useState(true);

  const handleLogout = async () => {
    try {
      await fetch("http://localhost:8080/v1/logout", {
        method: "POST",
        credentials: "include",
      });
    } catch (error) {
      console.error("Logout failed", error);
    }
    setUser(null);
  };

  const loadUser = async () => {
    try {
      const response = await fetch("http://localhost:8080/v1/me", {
        credentials: "include",
      });

      if (!response.ok) {
        throw new Error("Unauthorized");
      }

      const userData = await response.json();
      setUser(userData);
      return userData;
    } catch (error) {
      console.error("Failed to load current user", error);
      setUser(null);
    } finally {
      setLoadingUser(false);
    }
  };

  useEffect(() => {
    loadUser();
  }, []);

  return (
    <>
      <Router>
        <Routes>
          <Route path="/" element={<div className="website"><Authentication setUser={setUser} /></div>} />
          <Route
            path="/events"
            element={
              loadingUser ? (
                <div className="website full-page">Loading...</div>
              ) : user ? (
                <div className="website full-page"> 
                  <Events user={user} onLogout={handleLogout} onUserReload={loadUser} />
                </div>
              ) : (
                <Navigate to="/" replace />
              )
            }
          />
          <Route
            path="/events/new"
            element={
              loadingUser ? (
                <div className="website full-page">Loading...</div>
              ) : user ? (
                <div className="website full-page">
                  <CreateEvent user={user} />
                </div>
              ) : (
                <Navigate to="/" replace />
              )
            }
          />
          <Route
            path="/events/edit/:id"
            element={
              loadingUser ? (
                <div className="website full-page">Loading...</div>
              ) : user ? (
                <div className="website full-page">
                  <CreateEvent user={user} />
                </div>
              ) : (
                <Navigate to="/" replace />
              )
            }
          />
          <Route
            path="/admin"
            element={
              loadingUser ? (
                <div className="website full-page">Loading...</div>
              ) : user ? (
                <div className="website full-page">
                  <Admin user={user} onUserReload={loadUser} />
                </div>
              ) : (
                <Navigate to="/" replace />
              )
            }
          />
        </Routes>
      </Router>
    </>
  );
}

export default App;
