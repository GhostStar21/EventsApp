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
* @param {onUserReload} - a function parameter to reload user details
* @return {JSX Element} - All events that are currently registered
*/
function Events({ user, onLogout, onUserReload }) {

  const [events, setEvents] = useState([]);
  const [organizer, setOrganizer] = useState(null);
  const [showEarlierEvents, setShowEarlierEvents] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    fetchEvents();
  }, []);

  useEffect(() => {
    if (user?.role === "ORGANIZER" && user?.organizerId) {
      if (!organizer || organizer.id !== user.organizerId) {
        fetch(`http://localhost:8080/v1/organizers/${user.organizerId}`, {
          credentials: "include",
        })
          .then((res) => (res.ok ? res.json() : null))
          .then((data) => {
            if (data) setOrganizer(data);
          })
          .catch((err) => console.error("Failed to load organizer details", err));
      }
      return;
    }

    setOrganizer(null);
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
      const today = new Date();
      today.setHours(0, 0, 0, 0);

      const sortedEvents = (data || []).slice().sort((a, b) => {
        const aDate = new Date(a.event_date || a.date || a.eventDate || "");
        const bDate = new Date(b.event_date || b.date || b.eventDate || "");

        const aValid = !Number.isNaN(aDate.getTime());
        const bValid = !Number.isNaN(bDate.getTime());

        if (!aValid && !bValid) return 0;
        if (!aValid) return 1;
        if (!bValid) return -1;

        const aTime = aDate.getTime();
        const bTime = bDate.getTime();

        const aIsPast = aTime < today.getTime();
        const bIsPast = bTime < today.getTime();

        if (aIsPast && !bIsPast) return 1;
        if (!aIsPast && bIsPast) return -1;

        return aTime - bTime;
      });

      setEvents(sortedEvents);
    } catch (error) {
      console.error(error);
    }
  };

  const handleDeleteEvent = async (eventId) => {
    try {
      const response = await fetch(`http://localhost:8080/v1/events/${eventId}`, {
        method: "DELETE",
        credentials: "include",
      });

      if (!response.ok) {
        const data = await response.json().catch(() => null);
        throw new Error(data?.message || "Failed to delete event.");
      }

      setEvents((prev) => prev.filter((e) => e.id !== eventId));
    } catch (error) {
      console.error("Delete error:", error);
      alert(error.message || "Could not delete event.");
    }
  };

  const handleEditEvent = (eventToEdit) => {
    navigate(`/events/edit/${eventToEdit.id}`, { state: { event: eventToEdit, organizer } });
  };

  const today = new Date();
  today.setHours(0, 0, 0, 0);

  const upcomingEvents = events.filter((event) => {
    const eventDate = new Date(event.event_date || event.date || event.eventDate || "");
    return !Number.isNaN(eventDate.getTime()) && eventDate.getTime() >= today.getTime();
  });

  const earlierEvents = events.filter((event) => {
    const eventDate = new Date(event.event_date || event.date || event.eventDate || "");
    return Number.isNaN(eventDate.getTime()) || eventDate.getTime() < today.getTime();
  });

  return (
    <div className="events-page">
      <div className="events-header">
        <div className="events-header-left">
          <h1>Hva Skjer!</h1>
          {user?.role === "ORGANIZER" && (
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
        {upcomingEvents.map((event) => (
          <EventCard
            key={event.id}
            event={event}
            user={user}
            onEdit={handleEditEvent}
            onDelete={handleDeleteEvent}
          />
        ))}
      </div>

      {earlierEvents.length > 0 && (
        <div className="earlier-events-section">
          <button
            className="earlier-events-toggle"
            onClick={() => setShowEarlierEvents((prev) => !prev)}
          >
            {showEarlierEvents ? "Hide earlier events" : "Show earlier events"}
          </button>

          {showEarlierEvents && (
            <div className="event-grid earlier-events-grid">
              {earlierEvents.map((event) => (
                <EventCard
                  key={event.id}
                  event={event}
                  user={user}
                  onEdit={handleEditEvent}
                  onDelete={handleDeleteEvent}
                />
              ))}
            </div>
          )}
        </div>
      )}
  	</div>
  )
}

export default Events;
