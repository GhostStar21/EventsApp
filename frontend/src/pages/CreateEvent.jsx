import { useLocation, useNavigate } from "react-router-dom";
import { useState } from "react";
import "../styles/CreateEvent.css";

function CreateEvent() {
  const navigate = useNavigate();
  const { state } = useLocation();
  const organizer = state?.organizer;
  const [eventData, setEventData] = useState(
    {
      id: "", 
      name: "", 
      isExclusive: false,
      event_date: "",
      event_time: "",
      location: "",
      description: "",
      isRegistration: false,
    }
  );
  const [statusMessage, setStatusMessage] = useState("");
  const [isErrorMessage, setIsErrorMessage] = useState(false);
  const [eventCreated, setEventCreated] = useState(false);

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

    try {
      const response = await fetch("http://localhost:8080/v1/events", {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          name: eventData.name,
          isExclusive: eventData.isExclusive,
          event_date: eventData.event_date,
          event_time: eventData.event_time,
          location: eventData.location,
          description: eventData.description,
          isRegistration: eventData.isRegistration,
        }),
      });

      const data = await response.json().catch(() => null);
      console.log("Response status:", response.status);

      if (response.ok) {
        setStatusMessage("Event created successfully!");
        setIsErrorMessage(false);
        setEventCreated(true);
      } else {
        setStatusMessage(data?.message || "Could not create event.");
        setIsErrorMessage(true);
      }

    }
    catch (error) {
      console.error("Server error:", error);
      setStatusMessage("Could not reach the server.");
      setIsErrorMessage(true);
    }
  } 

  let eventOverlay = 
  <>
    <section className="create-event-card">
      <h1>Create event</h1>
      {organizer && <p>Creating an event for {organizer.name}.</p>}
      {eventCreated ? (
        <div className="event-created-message">
          <p>{statusMessage}</p>
          <button onClick={() => navigate("/events")}>
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

          {statusMessage && (
            <p className={isErrorMessage ? "form-message error" : "form-message success"}>
              {statusMessage}
            </p>
          )}

          <button className="submit-event-button" type="submit">Create event</button>
        </form>
      )}
    </section>
  
  </>


  return (
    <main className="create-event-page">
      {!eventCreated && (
        <button className="back-button" onClick={() => navigate("/events")}>
          Back to events
        </button>
      )}
      {eventOverlay}
    </main>
  );
}

export default CreateEvent;
