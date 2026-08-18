import { useEffect, useState } from "react";
import { Navigate, useNavigate } from "react-router-dom";
import "../styles/Admin.css";
import ConfirmModal from "../components/ConfirmModal.jsx";
import { API_URL } from "../config.js";

const blankUser = { id: "", name: "", email: "", role: "USER", password: "" };
const blankOrganizer = { id: "", name: "", orgNumber: "", type: "Personal" };
const blankEvent = {
  id: "",
  name: "",
  isExclusive: false,
  event_date: "",
  event_time: "",
  location: "",
  description: "",
  isRegistration: false,
  registrationLink: "",
  organizerId: "",
};

function Admin({ user, onUserReload }) {
  const navigate = useNavigate();
  const [section, setSection] = useState("users");
  const [users, setUsers] = useState([]);
  const [organizers, setOrganizers] = useState([]);
  const [organizerFilter, setOrganizerFilter] = useState("all");
  const [showAdvanced, setShowAdvanced] = useState(false);
  const [events, setEvents] = useState([]);
  const [userForm, setUserForm] = useState(blankUser);
  const [organizerForm, setOrganizerForm] = useState(blankOrganizer);
  const [eventForm, setEventForm] = useState(blankEvent);
  const [message, setMessage] = useState("");
  const [isError, setIsError] = useState(false);
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [confirmTitle, setConfirmTitle] = useState("");
  const [confirmMessage, setConfirmMessage] = useState("");
  const [confirmRequireInput, setConfirmRequireInput] = useState(false);
  const [confirmInputLabel, setConfirmInputLabel] = useState("");
  const [confirmInputValue, setConfirmInputValue] = useState("");
  const [pendingAction, setPendingAction] = useState(null);

  useEffect(() => {
    if (user?.role !== "ADMIN") return;
    loadAllData();
  }, [user]);

  const loadAllData = async () => {
    try {
      const [usersRes, organizersRes, eventsRes] = await Promise.all([
        fetch(`${API_URL}/v1/users`, { credentials: "include" }),
        fetch(`${API_URL}/v1/organizers`, { credentials: "include" }),
        fetch(`${API_URL}/v1/events`, { credentials: "include" }),
      ]);

      if (!usersRes.ok || !organizersRes.ok || !eventsRes.ok) {
        throw new Error("Failed to load admin data.");
      }

      const [usersData, organizersData, eventsData] = await Promise.all([
        usersRes.json(),
        organizersRes.json(),
        eventsRes.json(),
      ]);
      setUsers(usersData || []);
      setOrganizers(organizersData || []);
      setEvents(eventsData || []);
    } catch (err) {
      console.error(err);
      setMessage("Could not load admin data. Please refresh.");
      setIsError(true);
    }
  };

  const resetMessages = () => {
    setMessage("");
    setIsError(false);
  };

  const handleUserFormChange = (event) => {
    const { name, value } = event.target;
    setUserForm((prev) => ({ ...prev, [name]: value }));
  };

  const handleOrganizerFormChange = (event) => {
    const { name, value } = event.target;
    setOrganizerForm((prev) => ({ ...prev, [name]: value }));
  };

  const handleEventFormChange = (event) => {
    const { name, value, type, checked } = event.target;
    setEventForm((prev) => ({
      ...prev,
      [name]: type === "checkbox" ? checked : value,
    }));
  };

  const saveUser = async (event) => {
    event.preventDefault();
    resetMessages();

    const payload = {
      name: userForm.name,
      email: userForm.email,
      role: userForm.role,
    };
    if (!userForm.id && userForm.password) {
      payload.password = userForm.password;
    }

    const method = userForm.id ? "PUT" : "POST";
    const url = userForm.id ? `${API_URL}/v1/users/${userForm.id}` : `${API_URL}/v1/users`;

    try {
      const response = await fetch(url, {
        method,
        credentials: "include",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      const data = await response.json().catch(() => null);
      if (!response.ok) {
        throw new Error(data?.message || "Failed to save user.");
      }
      setMessage(userForm.id ? "User updated." : "User created.");
      setIsError(false);
      setUserForm(blankUser);
      await loadAllData();
      await onUserReload?.();
    } catch (err) {
      console.error(err);
      setMessage(err.message);
      setIsError(true);
    }
  };

  const deleteUser = (id) => {
    resetMessages();
    setConfirmTitle("Delete user?");
    setConfirmMessage("This will permanently delete the user.");
    setConfirmRequireInput(false);
    setPendingAction({ type: "deleteUser", id });
    setConfirmOpen(true);
  };

  const editUser = (userData) => {
    setUserForm({
      id: userData.id,
      name: userData.name || "",
      email: userData.email || "",
      role: userData.role || "USER",
      password: "",
    });
    setSection("users");
    resetMessages();
  };

  const saveOrganizer = async (event) => {
    event.preventDefault();
    resetMessages();

    const payload = {
      name: organizerForm.name,
      orgNumber: Number(organizerForm.orgNumber),
      type: organizerForm.type,
    };
    const method = organizerForm.id ? "PUT" : "POST";
    const url = organizerForm.id ? `${API_URL}/v1/organizers/${organizerForm.id}` : `${API_URL}/v1/organizers`;

    try {
      const response = await fetch(url, {
        method,
        credentials: "include",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      const data = await response.json().catch(() => null);
      if (!response.ok) {
        throw new Error(data?.message || "Failed to save organizer.");
      }
      setMessage(organizerForm.id ? "Organizer updated." : "Organizer created.");
      setIsError(false);
      setOrganizerForm(blankOrganizer);
      await loadAllData();
    } catch (err) {
      console.error(err);
      setMessage(err.message);
      setIsError(true);
    }
  };

  const deleteOrganizer = (id) => {
    resetMessages();
    setConfirmTitle("Delete organizer?");
    setConfirmMessage("This will permanently delete the organizer.");
    setConfirmRequireInput(false);
    setPendingAction({ type: "deleteOrganizer", id });
    setConfirmOpen(true);
  };

  const editOrganizer = (org) => {
    setOrganizerForm({
      id: org.id,
      name: org.name || "",
      orgNumber: org.orgnumber || org.orgNumber || "",
      type: org.type || "Personal",
    });
    setSection("organizers");
    resetMessages();
  };

  const setOrganizerStatus = (org, status) => {
    resetMessages();
    setConfirmTitle(`${status} organizer?`);
    setConfirmMessage(`Set organizer ${org.name} to ${status}?`);
    setConfirmRequireInput(false);
    setPendingAction({ type: "setOrganizerStatus", org, status });
    setConfirmOpen(true);
  };

  const deleteAll = (type) => {
    resetMessages();
    setConfirmTitle(`Delete all ${type}?`);
    setConfirmMessage(`This is irreversible. Enter your admin password to confirm.`);
    setConfirmRequireInput(true);
    setConfirmInputLabel("Admin password");
    setConfirmInputValue("");
    setPendingAction({ type: "deleteAll", subtype: type });
    setConfirmOpen(true);
  };

  const saveEvent = async (event) => {
    event.preventDefault();
    resetMessages();

    const payload = {
      name: eventForm.name,
      isExclusive: eventForm.isExclusive,
      event_date: eventForm.event_date ? `${eventForm.event_date}T00:00:00Z` : "",
      event_time: eventForm.event_time ? `0001-01-01T${eventForm.event_time}:00Z` : "",
      location: eventForm.location,
      description: eventForm.description,
      isRegistration: eventForm.isRegistration,
      registrationLink: eventForm.registrationLink,
    };
    if (eventForm.organizerId) {
      payload.organizer_id = Number(eventForm.organizerId);
    }

    const method = eventForm.id ? "PUT" : "POST";
    const url = eventForm.id ? `${API_URL}/v1/events/${eventForm.id}` : `${API_URL}/v1/events`;

    try {
      const response = await fetch(url, {
        method,
        credentials: "include",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      const data = await response.json().catch(() => null);
      if (!response.ok) {
        throw new Error(data?.message || "Failed to save event.");
      }
      setMessage(eventForm.id ? "Event updated." : "Event created.");
      setIsError(false);
      setEventForm(blankEvent);
      await loadAllData();
    } catch (err) {
      console.error(err);
      setMessage(err.message);
      setIsError(true);
    }
  };

  const deleteEvent = async (id) => {
    resetMessages();
    setConfirmTitle("Delete event?");
    setConfirmMessage("This will permanently delete the event.");
    setConfirmRequireInput(false);
    setPendingAction({ type: "deleteEvent", id });
    setConfirmOpen(true);
  };

  const editEvent = (ev) => {
    setEventForm({
      id: ev.id,
      name: ev.name || "",
      isExclusive: Boolean(ev.isExclusive),
      event_date: ev.event_date ? String(ev.event_date).split("T")[0] : "",
      event_time: ev.event_time
        ? ev.event_time.includes("T")
          ? ev.event_time.split("T")[1].substring(0, 5)
          : ev.event_time.substring(0, 5)
        : "",
      location: ev.location || "",
      description: ev.description || "",
      isRegistration: Boolean(ev.isRegistration),
      registrationLink: ev.registrationLink || "",
      organizerId: ev.organizer_id || ev.organizerId || "",
    });
    setSection("events");
    resetMessages();
  };

  const handleConfirmCancel = () => {
    setConfirmOpen(false);
    setPendingAction(null);
    setConfirmInputValue("");
  };

  const handleConfirm = async (inputValue) => {
    setConfirmOpen(false);
    const action = pendingAction;
    setPendingAction(null);
    try {
      if (!action) return;
      if (action.type === "deleteUser") {
        const response = await fetch(`${API_URL}/v1/users/${action.id}`, { method: "DELETE", credentials: "include" });
        if (!response.ok) throw new Error("Failed to delete user");
        setMessage("User deleted.");
        await loadAllData();
        await onUserReload?.();
        return;
      }
      if (action.type === "deleteOrganizer") {
        const response = await fetch(`${API_URL}/v1/organizers/${action.id}`, { method: "DELETE", credentials: "include" });
        if (!response.ok) throw new Error("Failed to delete organizer");
        setMessage("Organizer deleted.");
        await loadAllData();
        return;
      }
      if (action.type === "setOrganizerStatus") {
        const { org, status } = action;
        const payload = {
          id: org.id,
          name: org.name,
          orgnumber: org.orgnumber || org.orgNumber || 0,
          type: org.type || "Personal",
          is_approved: status,
        };
        const response = await fetch(`${API_URL}/v1/organizers/${org.id}`, {
          method: "PUT",
          credentials: "include",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(payload),
        });
        const data = await response.json().catch(() => null);
        if (!response.ok) throw new Error(data?.message || "Failed to update organizer approval");
        setMessage(`Organizer ${status.toLowerCase()}.`);
        await loadAllData();
        return;
      }
      if (action.type === "deleteAll") {
        const url = `${API_URL}/v1/${action.subtype}`;
        // Note: backend expects admin privileges; password check should be server-side.
        const response = await fetch(url, { method: "DELETE", credentials: "include" });
        if (!response.ok) throw new Error("Failed to delete");
        setMessage(`All ${action.subtype} deleted.`);
        await loadAllData();
        return;
      }
      if (action.type === "deleteEvent") {
        const response = await fetch(`${API_URL}/v1/events/${action.id}`, { method: "DELETE", credentials: "include" });
        if (!response.ok) throw new Error("Failed to delete event");
        setMessage("Event deleted.");
        await loadAllData();
        return;
      }
    } catch (err) {
      console.error(err);
      setMessage(err.message || "Action failed");
      setIsError(true);
    }
  };

  if (!user) {
    return <Navigate to="/" replace />;
  }

  if (user.role !== "ADMIN") {
    return <Navigate to="/events" replace />;
  }

  const displayedOrganizers = organizers.filter((o) => {
    if (organizerFilter === "pending") return (o.is_approved || o.isApproved || o.IsApproved || "PENDING") === "PENDING";
    return true;
  });

  return (
    <div className="admin-root">
    <div className="admin-page">
      <div className="admin-header">
        <div>
          <h1>Admin dashboard</h1>
          <p>Manage users, organizers, and events.</p>
        </div>
        <div style={{display: 'flex', gap: '0.5rem', alignItems: 'center'}}>
          <button className="admin-back-btn" onClick={() => navigate("/events")}>Back to events</button>
          <button className="admin-back-btn" onClick={() => setShowAdvanced((s) => !s)}>{showAdvanced ? 'Hide Advanced' : 'Advanced'}</button>
        </div>
      </div>

      <div className="admin-tabs">
        <button className={section === "users" ? "active" : ""} onClick={() => setSection("users")}>Users</button>
        <button className={section === "organizers" ? "active" : ""} onClick={() => setSection("organizers")}>Organizers</button>
        <button className={section === "events" ? "active" : ""} onClick={() => setSection("events")}>Events</button>
      </div>

      {message && <div className={`admin-message ${isError ? "error" : "success"}`}>{message}</div>}

      {showAdvanced && (
        <div className="admin-section" style={{marginBottom: '1rem'}}>
          <div className="admin-grid">
            <div className="admin-list-card">
              <h3>Advanced: Dangerous actions</h3>
              <p>Deleting is permanent. You will be asked to confirm with your admin password.</p>
              <div className="advanced-actions" style={{marginTop: '0.5rem'}}>
                <button className="danger" onClick={() => deleteAll('events')}>Delete all events</button>
                <button className="danger" onClick={() => deleteAll('users')}>Delete all users</button>
                <button className="danger" onClick={() => deleteAll('organizers')}>Delete all organizers</button>
              </div>
            </div>
          </div>
        </div>
      )}

      {section === "users" && (
        <div className="admin-section">
          <div className="admin-grid">
            <div className="admin-list-card">
              <h2>Users</h2>
              <div className="table-scroll">
              <table>
                <thead>
                  <tr>
                    <th>ID</th>
                    <th>Name</th>
                    <th>Email</th>
                    <th>Role</th>
                    <th>Organizer</th>
                    <th>Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {users.map((u) => (
                    <tr key={u.id}>
                      <td>{u.id}</td>
                      <td>{u.name}</td>
                      <td>{u.email}</td>
                      <td>{u.role}</td>
                      <td>{u.organizerId || "—"}</td>
                      <td>
                        <div className="action-buttons">
                          <button className="small-btn btn-edit" onClick={() => editUser(u)}>Edit</button>
                          <button className="small-btn btn-delete" onClick={() => deleteUser(u.id)}>Delete</button>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
              </div>
            </div>
            <div className="admin-form-card">
              <h2>{userForm.id ? "Edit user" : "Create user"}</h2>
              <form onSubmit={saveUser}>
                <label>Name</label>
                <input name="name" value={userForm.name} onChange={handleUserFormChange} required />
                <label>Email</label>
                <input name="email" type="email" value={userForm.email} onChange={handleUserFormChange} required />
                <label>Role</label>
                <select name="role" value={userForm.role} onChange={handleUserFormChange}>
                  <option value="USER">USER</option>
                  <option value="ORGANIZER">ORGANIZER</option>
                  <option value="ADMIN">ADMIN</option>
                </select>
                {!userForm.id && (
                  <>
                    <label>Password</label>
                    <input name="password" type="password" value={userForm.password} onChange={handleUserFormChange} required />
                  </>
                )}
                <button type="submit">{userForm.id ? "Save user" : "Create user"}</button>
                {userForm.id && <button type="button" onClick={() => setUserForm(blankUser)}>Clear</button>}
              </form>
            </div>
          </div>
        </div>
      )}

      {section === "organizers" && (
        <div className="admin-section">
          <div className="admin-grid">
            <div className="admin-list-card">
              <h2>Organizers</h2>
              <div className="organizer-filters" style={{marginBottom: '0.5rem'}}>
                <button className={organizerFilter === 'all' ? 'active' : ''} onClick={() => setOrganizerFilter('all')}>All</button>
                <button className={organizerFilter === 'pending' ? 'active' : ''} onClick={() => setOrganizerFilter('pending')}>Requires approval</button>
              </div>
              <div className="table-scroll">
              <table>
                <thead>
                  <tr>
                    <th>ID</th>
                    <th>Name</th>
                    <th>Org number</th>
                    <th>Type</th>
                    <th>Status</th>
                    <th>Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {displayedOrganizers.map((org) => (
                    <tr key={org.id}>
                      <td>{org.id}</td>
                      <td>{org.name}</td>
                      <td>{org.orgnumber || org.orgNumber}</td>
                      <td>{org.type}</td>
                      <td>{org.is_approved || org.isApproved || org.IsApproved || 'PENDING'}</td>
                      <td>
                        <div className="action-buttons">
                          <button className="small-btn btn-edit" onClick={() => editOrganizer(org)}>Edit</button>
                          <button className="small-btn btn-delete" onClick={() => deleteOrganizer(org.id)}>Delete</button>
                          {( (org.is_approved || org.isApproved || org.IsApproved || 'PENDING') === 'PENDING') && (
                            <>
                              <button className="small-btn btn-approve" onClick={() => setOrganizerStatus(org, 'APPROVED')}>Approve</button>
                              <button className="small-btn btn-reject" onClick={() => setOrganizerStatus(org, 'REJECTED')}>Reject</button>
                              <button className="small-btn btn-suspend" onClick={() => setOrganizerStatus(org, 'SUSPENDED')}>Suspend</button>
                            </>
                          )}
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
              </div>
            </div>
            <div className="admin-form-card">
              <h2>{organizerForm.id ? "Edit organizer" : "Create organizer"}</h2>
              <form onSubmit={saveOrganizer}>
                <label>Name</label>
                <input name="name" value={organizerForm.name} onChange={handleOrganizerFormChange} required />
                <label>Org Number</label>
                <input name="orgNumber" type="number" value={organizerForm.orgNumber} onChange={handleOrganizerFormChange} required />
                <label>Type</label>
                <select name="type" value={organizerForm.type} onChange={handleOrganizerFormChange}>
                  <option value="Personal">Personal</option>
                  <option value="Official">Official</option>
                </select>
                <button type="submit">{organizerForm.id ? "Save organizer" : "Create organizer"}</button>
                {organizerForm.id && <button type="button" onClick={() => setOrganizerForm(blankOrganizer)}>Clear</button>}
              </form>
            </div>
          </div>
        </div>
      )}

      {section === "events" && (
        <div className="admin-section">
          <div className="admin-grid">
            <div className="admin-list-card">
              <h2>Events</h2>
              <div className="table-scroll">
              <table>
                <thead>
                  <tr>
                    <th>ID</th>
                    <th>Name</th>
                    <th>Date</th>
                    <th>Organizer</th>
                    <th>Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {events.map((ev) => (
                    <tr key={ev.id}>
                      <td>{ev.id}</td>
                      <td>{ev.name}</td>
                      <td>{ev.event_date ? ev.event_date.split("T")[0] : "—"}</td>
                      <td>{ev.organizer_name || ev.organizer_name || "—"}</td>
                      <td>
                        <div className="action-buttons">
                          <button className="small-btn btn-edit" onClick={() => editEvent(ev)}>Edit</button>
                          <button className="small-btn btn-delete" onClick={() => deleteEvent(ev.id)}>Delete</button>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
              </div>
            </div>
            <div className="admin-form-card">
              <h2>{eventForm.id ? "Edit event" : "Create event"}</h2>
              <form onSubmit={saveEvent}>
                <label>Name</label>
                <input name="name" value={eventForm.name} onChange={handleEventFormChange} required />
                <div className="form-row-two">
                  <div>
                    <label>Date</label>
                    <input name="event_date" type="date" value={eventForm.event_date} onChange={handleEventFormChange} required />
                  </div>
                  <div>
                    <label>Time</label>
                    <input name="event_time" type="time" value={eventForm.event_time} onChange={handleEventFormChange} required />
                  </div>
                </div>
                <label>Location</label>
                <input name="location" value={eventForm.location} onChange={handleEventFormChange} required />
                <label>Description</label>
                <textarea name="description" value={eventForm.description} onChange={handleEventFormChange} rows="4" />
                <label>Organizer</label>
                <select name="organizerId" value={eventForm.organizerId} onChange={handleEventFormChange}>
                  <option value="">No organizer</option>
                  {organizers.map((org) => (
                    <option key={org.id} value={org.id}>{org.name}</option>
                  ))}
                </select>
                <label className="checkbox-label">
                  <input name="isExclusive" type="checkbox" checked={eventForm.isExclusive} onChange={handleEventFormChange} />
                  Exclusive event
                </label>
                <label className="checkbox-label">
                  <input name="isRegistration" type="checkbox" checked={eventForm.isRegistration} onChange={handleEventFormChange} />
                  Registration required
                </label>
                {eventForm.isRegistration && (
                  <>
                    <label>Registration link</label>
                    <input name="registrationLink" type="url" value={eventForm.registrationLink} onChange={handleEventFormChange} />
                  </>
                )}
                <button type="submit">{eventForm.id ? "Save event" : "Create event"}</button>
                {eventForm.id && <button type="button" onClick={() => setEventForm(blankEvent)}>Clear</button>}
              </form>
            </div>
          </div>
        </div>
      )}
      
    </div>
    <ConfirmModal
      open={confirmOpen}
      title={confirmTitle}
      message={confirmMessage}
      requireInput={confirmRequireInput}
      inputLabel={confirmInputLabel}
      inputValue={confirmInputValue}
      onInputChange={setConfirmInputValue}
      onCancel={handleConfirmCancel}
      onConfirm={handleConfirm}
    />
    </div>
  );
}

export default Admin;
