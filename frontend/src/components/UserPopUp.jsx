import { useEffect, useRef, useState } from "react";
import "../styles/UserPopUp.css";

function UserPopUp({ name, email, onLogout, organizer, onOrganizerChange }) {
  const [isOpen, setIsOpen] = useState(false);
  const popupRef = useRef(null);

  useEffect(() => {
    const handleClickOutside = (event) => {
      if (popupRef.current && !popupRef.current.contains(event.target)) {
        setIsOpen(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const handleLogout = () => {
    setIsOpen(false);
    onLogout?.();
  };

  const handleRoleChange = async () => {
  try {
    const response = await fetch("http://localhost:8080/v1/register-organizer", {
      method: "POST",
      credentials: "include",
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.message || "Could not switch to organizer");
    }

    const organizerResponse = await fetch(
      `http://localhost:8080/v1/organizers/${data.organizerId}`
    );

    if (!organizerResponse.ok) {
      throw new Error("Could not load organizer details");
    }

    onOrganizerChange(await organizerResponse.json());
    setIsOpen(false);
    
    
  } catch (error) {
      console.error("Server error:", error);
  }
  
};

  const handleDemote = async () => {
    try {
      const response = await fetch("http://localhost:8080/v1/demote-organizer", {
        method: "POST",
        credentials: "include",
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || "Could not switch back to user");
      }

      onOrganizerChange(null);
      setIsOpen(false);
    }
    catch {
      console.error("Server error:", error);
    }
  }

  return (
    <div className="user-popup" ref={popupRef}>
      <button
        className="user-popup-trigger"
        onClick={() => setIsOpen((prev) => !prev)}
        aria-expanded={isOpen}
      >
        <span className="user-avatar">{(name || "U").charAt(0).toUpperCase()}</span>
        <span className="user-name"> {organizer ? organizer.name : name || "Profile"}</span>
      </button>

      {isOpen && (
        <div className="user-popup-menu" role="menu">
          <div className="user-popup-info">
            <p className="user-popup-name">{name || "User"}</p>
            <p className="user-popup-email">{email || "No email provided"}</p>
            {organizer ? (
              <button 
              className="user-popup-trigger"
              onClick={handleDemote}
              >
                Switch to User
              </button>
            ) : (
              <button
              className="user-popup-trigger"
              onClick={handleRoleChange}
              >
                Switch to Organizer
              </button>
            )}
          </div>

          <button className="user-popup-logout" onClick={handleLogout}>
            Log Out
          </button>
        </div>
      )}
    </div>
  );
}

export default UserPopUp;