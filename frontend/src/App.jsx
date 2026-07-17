import "./styles/index.css"
import { useEffect, useState } from "react";
import Authentication from "./components/Authentication.jsx";
import { BrowserRouter as Router, Routes, Route, Navigate } from "react-router-dom"
import Events from "./pages/Events.jsx"

function App() {
  const [user, setUser] = useState(null);
  const [loadingUser, setLoadingUser] = useState(true);

  const handleLogout = () => {
    localStorage.removeItem("token");
    setUser(null);
  };

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (!token) {
      setLoadingUser(false);
      return;
    }

    const loadUser = async () => {
      try {
        const response = await fetch("http://localhost:8080/v1/me", {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (!response.ok) {
          throw new Error("Unauthorized");
        }

        const userData = await response.json();
        setUser(userData);
      } catch (error) {
        console.error("Failed to load current user", error);
        localStorage.removeItem("token");
        setUser(null);
      } finally {
        setLoadingUser(false);
      }
    };

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
                <div className="website full-page"><Events user={user} onLogout={handleLogout} /></div>
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