import { useEffect, useRef, useState } from "react";
import "../styles/UserPopUp.css";

/**
 * Showcases the User icon, presenting either the ability to switch to an organizer or create an organizer, 
 * thereby giving permission to create events.
 * @param {name} - The name of the user logged in.
 * @param {email} - The email of the user logged in. 
 * @param {role} - The role of the user, granting them different rights.
 * @param {isOrganizerMember} - A boolean that tells us if the user belongs to an organizer or not. 
 * @param {onLogout} - A function  
 * @param {organizer} - The organizer to which the user is potentially a member of.
 * @param {onOrganizerChange} - A function granting access to a user (upon authentication) to create events as an organizer.
 * @param {onUserReload} - A function that loads the user that is performing actions.
 * @returns 
 */
function UserPopUp( {name,
  email,
  role,
  isOrganizerMember,
  onLogout,
  organizer,
  onOrganizerChange,
  onUserReload,
}) {
  const [isOpen, setIsOpen] = useState(false);
  const [showCreatePrompt, setShowCreatePrompt] = useState(false);
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [formData, setFormData] = useState({
    name: "",
    orgNumber: "",
    type: "Personal",
  });
  const [statusMessage, setStatusMessage] = useState("");
  const [isErrorMessage, setIsErrorMessage] = useState(false);
  const popupRef = useRef(null);

  // If the user clicks outside the popup, close the popup.
  useEffect(() => {
    const handleClickOutside = (event) => {
      if (popupRef.current && !popupRef.current.contains(event.target)) {
        setIsOpen(false);
        setShowCreatePrompt(false);
        setShowCreateForm(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const resetModalState = () => {
    setShowCreatePrompt(false);
    setShowCreateForm(false);
    setStatusMessage("");
    setIsErrorMessage(false);
  };

  const handleLogout = () => {
    setIsOpen(false);
    resetModalState();
    onLogout?.();
  };

  const isCurrentOrganizerMode = role === "ORGANIZER";
  const hasOrganizerMembership = Boolean(isOrganizerMember || isCurrentOrganizerMode);
  const displayName = isCurrentOrganizerMode && organizer ? organizer.name : name || "Profile";

  // Responsible for switching role between user and organizer.
  const handleRoleChange = async () => {
    if (!hasOrganizerMembership) {
      setShowCreatePrompt(true);
      setShowCreateForm(false);
      setStatusMessage("");
      setIsErrorMessage(false);
      return;
    }

    try {
      const response = await fetch("http://localhost:8080/v1/promote-organizer", {
        method: "POST",
        credentials: "include",
      });

      const data = await response.json().catch(() => null);

      if (!response.ok) {
        throw new Error(data?.message || "Could not switch to organizer");
      }

      const organizerResponse = await fetch(
        `http://localhost:8080/v1/organizers/${data.organizerId}`,
        { credentials: "include" }
      );

      if (!organizerResponse.ok) {
        throw new Error("Could not load organizer details");
      }

      onOrganizerChange(await organizerResponse.json());
      await onUserReload?.();
      setIsOpen(false);
      resetModalState();
    } catch (error) {
      console.error("Server error:", error);
      setStatusMessage(error.message || "Could not switch to organizer");
      setIsErrorMessage(true);
    }
  };

  const handleCreateOrganizer = async (event) => {
    event.preventDefault();
    setStatusMessage("");
    setIsErrorMessage(false);

    try {
      const response = await fetch("http://localhost:8080/v1/promote-organizer", {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          name: formData.name,
          orgNumber: Number(formData.orgNumber),
          type: formData.type,
        }),
      });

      const data = await response.json().catch(() => null);

      if (!response.ok) {
        throw new Error(data?.message || "Could not create organizer");
      }

      const organizerResponse = await fetch(
        `http://localhost:8080/v1/organizers/${data.organizerId}`,
        { credentials: "include" }
      );

      if (!organizerResponse.ok) {
        throw new Error("Could not load organizer details");
      }

      onOrganizerChange(await organizerResponse.json());
      await onUserReload?.();
      setIsOpen(false);
      resetModalState();
    } catch (error) {
      console.error("Server error:", error);
      setStatusMessage(error.message || "Could not create organizer");
      setIsErrorMessage(true);
    }
  };

  const handleDemote = async () => {
    try {
      const response = await fetch("http://localhost:8080/v1/demote-organizer", {
        method: "POST",
        credentials: "include",
      });

      const data = await response.json().catch(() => null);

      if (!response.ok) {
        throw new Error(data?.message || "Could not switch back to user");
      }

      onOrganizerChange(null);
      await onUserReload?.();
      setIsOpen(false);
      resetModalState();
    } catch (error) {
      console.error("Server error:", error);
      setStatusMessage(error.message || "Could not switch back to user");
      setIsErrorMessage(true);
    }
  };

  return (
    <div className="user-popup" ref={popupRef}>
      <button
        className="user-popup-trigger"
        onClick={() => setIsOpen((prev) => !prev)}
        aria-expanded={isOpen}
      >
        <span className="user-avatar">{(name || "U").charAt(0).toUpperCase()}</span>
        <span className="user-name"> {displayName}</span>
      </button>

      {isOpen && (
        <div className="user-popup-menu" role="menu">
          <div className="user-popup-info">
            <p className="user-popup-name">{name || "User"}</p>
            <p className="user-popup-email">{email || "No email provided"}</p>
            {isCurrentOrganizerMode ? (
              <button className="user-popup-trigger" onClick={handleDemote}>
                Switch to User
              </button>
            ) : hasOrganizerMembership ? (
              <button className="user-popup-trigger" onClick={handleRoleChange}>
                Switch to Organizer
              </button>
            ) : (
              <button className="user-popup-trigger" onClick={handleRoleChange}>
                Create Organizer
              </button>
            )}
          </div>

          <button className="user-popup-logout" onClick={handleLogout}>
            Log Out
          </button>
        </div>
      )}

      {(showCreatePrompt || showCreateForm) && (
        <div className="organizer-modal-overlay">
          <div className="organizer-modal">
            {statusMessage && (
              <p className={`organizer-status ${isErrorMessage ? "error" : "success"}`}>
                {statusMessage}
              </p>
            )}

            {showCreatePrompt && !showCreateForm && (
              <>
                <h3>You are currently not part of an organizer.</h3>
                <p>Would you like to create your own organizer?</p>
                <div className="organizer-modal-actions">
                  <button className="organizer-submit-btn" onClick={() => setShowCreateForm(true)}>
                    Yes
                  </button>
                  <button className="organizer-cancel-btn" onClick={resetModalState}>
                    No
                  </button>
                </div>
              </>
            )}

            {showCreateForm && (
              <form className="organizer-form" onSubmit={handleCreateOrganizer}>
                <h3>Create Organizer</h3>
                <input
                  type="text"
                  placeholder="Organizer name"
                  value={formData.name}
                  onChange={(event) => setFormData((prev) => ({ ...prev, name: event.target.value }))}
                  required
                />
                <input
                  type="number"
                  placeholder="Organization number"
                  value={formData.orgNumber}
                  onChange={(event) => setFormData((prev) => ({ ...prev, orgNumber: event.target.value }))}
                  required
                />
                <select
                  value={formData.type}
                  onChange={(event) => setFormData((prev) => ({ ...prev, type: event.target.value }))}
                >
                  <option value="Personal">Personal</option>
                  <option value="Official">Official</option>
                </select>
                <div className="organizer-modal-actions">
                  <button type="submit" className="organizer-submit-btn">
                    Create organizer
                  </button>
                  <button type="button" className="organizer-cancel-btn" onClick={resetModalState}>
                    Cancel
                  </button>
                </div>
              </form>
            )}
          </div>
        </div>
      )}
    </div>
  );
}

export default UserPopUp;