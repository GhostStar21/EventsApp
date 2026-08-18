import { useLocation, useNavigate, useParams } from "react-router-dom";
import { useEffect, useState } from "react";
import "../styles/CreateEvent.css";
import { API_URL } from "../config.js";

/** 
 * Displays a form for creating or updating an event.
 * 
 * @returns JSX element - The rendered form for creating/updating an event.
*/
function CreateEvent() {
  const navigate = useNavigate();
  const { id } = useParams();
  const { state } = useLocation();
  const organizer = state?.organizer;
  const isEdit = Boolean(id);

  const [eventData, setEventData] = useState({
    id: "", 
    name: "", 
    isExclusive: false,
    event_date: "",
    event_time: "",
    location: "",
    description: "",
    isRegistration: false,
    registrationLink: "",
  });

  const [statusMessage, setStatusMessage] = useState("");
  const [isErrorMessage, setIsErrorMessage] = useState(false);
  const [eventSubmitted, setEventSubmitted] = useState(false);

  useEffect(() => {
    if (isEdit) {
      if (state?.event) {
        populateForm(state.event);
      } else {
        fetch(`${API_URL}/v1/events/${id}`, {
          credentials: "include",
        })
          .then((res) => (res.ok ? res.json() : Promise.reject("Event not found")))
          .then((data) => populateForm(data))
          .catch((err) => {
            console.error("Error loading event", err);
            setStatusMessage("Could not load event for editing.");
            setIsErrorMessage(true);
          });
      }
    }
  }, [id, isEdit, state?.event]);

  const populateForm = (ev) => {
    let dateStr = "";
    if (ev.event_date) {
      dateStr = ev.event_date.split("T")[0];
    } else if (ev.date) {
      dateStr = ev.date.split("T")[0];
    }

    let timeStr = "";
    if (ev.event_time) {
      timeStr = ev.event_time.includes("T") ? ev.event_time.split("T")[1].substring(0, 5) : ev.event_time.substring(0, 5);
    } else if (ev.time) {
      timeStr = ev.time.includes("T") ? ev.time.split("T")[1].substring(0, 5) : ev.time.substring(0, 5);
    }

    setEventData({
      id: ev.id,
      name: ev.name || ev.title || "",
      isExclusive: Boolean(ev.isExclusive),
      event_date: dateStr,
      event_time: timeStr,
      location: ev.location || "",
      description: ev.description || "",
      isRegistration: Boolean(ev.isRegistration),
      registrationLink: ev.registrationLink || "",
    });
  };

  const handleInputChange = (event) => {
    const { name, type, value, checked } = event.target;

    setEventData((previous) => ({
      ...previous,
      [name]: type === "checkbox" ? checked : value,
    }));
  };

  const handleEvent = async (event) => {
    setStatusMessage("");
    setIsErrorMessage(false);
    event.preventDefault();

    const url = isEdit ? `${API_URL}/v1/events/${id}` : `${API_URL}/v1/events`;
    const method = isEdit ? "PUT" : "POST";

    const formattedDate = eventData.event_date
      ? (eventData.event_date.includes("T") ? eventData.event_date : `${eventData.event_date}T00:00:00Z`)
      : "";
      
    const formattedTime = eventData.event_time
      ? (eventData.event_time.includes("T") ? eventData.event_time : `0001-01-01T${eventData.event_time}:00Z`)
      : "";

    try {
      const response = await fetch(url, {
        method,
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          id: isEdit ? Number(id) : undefined,
          name: eventData.name,
          isExclusive: eventData.isExclusive,
          event_date: formattedDate,
          event_time: formattedTime,
          location: eventData.location,
          description: eventData.description,
          isRegistration: eventData.isRegistration,
          registrationLink: eventData.registrationLink,
        }),
      });

      const data = await response.json().catch(() => null);

      if (response.ok) {
        setStatusMessage(isEdit ? "Event updated successfully!" : "Event created successfully!");
        setIsErrorMessage(false);
        setEventSubmitted(true);
      } else {
        setStatusMessage(data?.message || (isEdit ? "Could not update event." : "Could not create event."));
        setIsErrorMessage(true);
      }

    } catch (error) {
      console.error("Server error:", error);
      setStatusMessage("Could not reach the server.");
      setIsErrorMessage(true);
    }
  };

  let eventOverlay = (
    <section className="create-event-card">
      <h1>{isEdit ? "Edit event" : "Create event"}</h1>
      {organizer && <p>{isEdit ? "Editing" : "Creating"} event for {organizer.name}.</p>}
      {eventSubmitted ? (
        <div className="event-created-message">
          <p>{statusMessage}</p>
          <button className="submit-event-button" onClick={() => navigate("/events")}>
            Back to events
          </button>
        </div>
      ) : (
        <form className="create-event-form" onSubmit={handleEvent}>
          <div className="form-field">
            <label htmlFor="name">Event name</label>
            <input
              id="name"
              type="text"
              name="name"
              placeholder="e.g. Autumn welcome party"
              value={eventData.name}
              onChange={handleInputChange}
              required
            />
          </div>

          <div className="form-row">
            <div className="form-field">
              <label htmlFor="event_date">Date</label>
              <input
                id="event_date"
                type="date"
                name="event_date"
                value={eventData.event_date}
                onChange={handleInputChange}
                required
              />
            </div>

            <div className="form-field">
              <label htmlFor="event_time">Time</label>
              <input
                id="event_time"
                type="time"
                name="event_time"
                value={eventData.event_time}
                onChange={handleInputChange}
                required
              />
            </div>
          </div>

          <div className="form-field">
            <label htmlFor="location">Location</label>
            <input
              id="location"
              type="text"
              name="location"
              placeholder="e.g. Trondheim Spektrum"
              value={eventData.location}
              onChange={handleInputChange}
              required
            />
          </div>

          <div className="form-field">
            <label htmlFor="description">Description</label>
            <textarea
              id="description"
              name="description"
              placeholder="Tell attendees what to expect"
              value={eventData.description}
              onChange={handleInputChange}
              rows="5"
            />
          </div>

          <label className="checkbox-field">
            <input
              type="checkbox"
              name="isExclusive"
              checked={eventData.isExclusive}
              onChange={handleInputChange}
            />
            Exclusive event (Only Linjeforening members?)
          </label>

          <label className="checkbox-field">
            <input
              type="checkbox"
              name="isRegistration"
              checked={eventData.isRegistration}
              onChange={handleInputChange}
            />
            Registration required
          </label>

          {eventData.isRegistration && (
            <div className="form-field">
              <label htmlFor="registrationLink">Registration link</label>
              <input
                id="registrationLink"
                type="url"
                name="registrationLink"
                placeholder="https://example.com/register"
                value={eventData.registrationLink}
                onChange={handleInputChange}
              />
            </div>
          )}

          {statusMessage && (
            <p className={isErrorMessage ? "form-message error" : "form-message success"}>
              {statusMessage}
            </p>
          )}

          <button className="submit-event-button" type="submit">
            {isEdit ? "Update event" : "Create event"}
          </button>
        </form>
      )}
    </section>
  );

  return (
    <main className="create-event-page">
      {!eventSubmitted && (
        <button className="back-button" onClick={() => navigate("/events")}>
          ← Back to events
        </button>
      )}
      {eventOverlay}
    </main>
  );
}

export default CreateEvent;
