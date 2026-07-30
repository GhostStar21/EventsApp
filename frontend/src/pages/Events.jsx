import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import EventCard from "../components/EventCard";
import UserPopUp from "../components/UserPopUp.jsx";
import "../styles/Events.css";

/** 
* Events renders the main Events page, showcasing all available events
*
* @param {user} - The current user logged in
* @param {onLogout} - a function parameter used for logging out a user
* @return {JSX Element} - All events that are currently registered (
* TODO: eventually display events based on organizers user is subscribed to)
*/
function Events({ user, onLogout, onUserReload }) {

  const [events, setEvents] = useState([]);
  const [organizer, setOrganizer] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    fetchEvents();
  }, []);

  useEffect(() => {
    if (user?.role === "ORGANIZER" && user?.organizerId && !organizer) {
      fetch(`http://localhost:8080/v1/organizers/${user.organizerId}`, {
        credentials: "include",
      })
        .then((res) => (res.ok ? res.json() : null))
        .then((data) => {
          if (data) setOrganizer(data);
        })
        .catch((err) => console.error("Failed to load organizer details", err));
    }
  }, [user?.role, user?.organizerId, organizer]);

  const fetchEvents = async () => {
    try {
      const response = await fetch("http://localhost:8080/v1/events", {
        credentials: "include",
      });

      if (!response.ok) {
        throw new Error("Failed to fetch events.");
      }

      const data = await response.json();
      setEvents(data);
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <div className="events-page">
      <div className="events-header">
        <div className="events-header-left">
          <h1>Hva Skjer!</h1>
          {user?.role === "ORGANIZER" && organizer && (
            <button
            className="create-event-button"
            onClick={() => navigate("/events/new", { state: { organizer } })}
            >
              Create event
            </button>
        )}
        </div>
        <UserPopUp 
          name={user?.name} 
          email={user?.email} 
          role={user?.role}
          isOrganizerMember={user?.isOrganizerMember}
          onLogout={onLogout} 
          organizer={organizer} 
          onOrganizerChange={setOrganizer}
          onUserReload={onUserReload}
          />
      </div>

      <div className="event-grid">
        {events.map(event => (
          <EventCard
              key={event.id}
              event={event}
          />
        ))}
      </div>
  	</div>
  )
}

export default Events;
