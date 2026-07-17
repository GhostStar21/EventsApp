import { useEffect, useRef, useState } from "react";
import "../styles/UserPopUp.css";

function UserPopUp({ name, email, onLogout }) {
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

  return (
    <div className="user-popup" ref={popupRef}>
      <button
        className="user-popup-trigger"
        onClick={() => setIsOpen((prev) => !prev)}
        aria-expanded={isOpen}
      >
        <span className="user-avatar">{(name || "U").charAt(0).toUpperCase()}</span>
        <span className="user-name">{name || "Profile"}</span>
      </button>

      {isOpen && (
        <div className="user-popup-menu" role="menu">
          <div className="user-popup-info">
            <p className="user-popup-name">{name || "User"}</p>
            <p className="user-popup-email">{email || "No email provided"}</p>
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