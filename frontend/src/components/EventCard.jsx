import { useState } from "react";
import "../styles/EventCard.css";

/**
 * Displays an event card, with a control to show its description.
 * 
 * @param {event} - An event fetched from the backend side
 * @param {user} - Currently authenticated user
 * @param {onEdit} - Function to trigger editing an event
 * @param {onDelete} - Function to trigger deleting an event
 * @returns JSX element - The rendered event card
*/

// XSS Protection: Validate and sanitize image URLs
const getSafeImageUrl = (url) => {
  if (!url) return "https://images.unsplash.com/photo-1501281668745-f7f57925c3b4?auto=format&fit=crop&w=600&q=80";
  try {
    const parsedUrl = new URL(url);
    // Only allow http and https protocols
    if (parsedUrl.protocol === 'http:' || parsedUrl.protocol === 'https:') {
      return url;
    }
  } catch (e) {
    // Invalid URL format
    console.warn('Invalid image URL provided:', url);
  }
  return "https://images.unsplash.com/photo-1501281668745-f7f57925c3b4?auto=format&fit=crop&w=600&q=80";
};

// XSS Protection: Validate and sanitize registration URLs
const getSafeRegistrationLink = (link) => {
  if (!link) return null;
  try {
    const parsedUrl = new URL(link);
    // Only allow http and https protocols
    if (parsedUrl.protocol === 'http:' || parsedUrl.protocol === 'https:') {
      return link;
    }
  } catch (e) {
    // Invalid URL format
    console.warn('Invalid registration URL provided:', link);
  }
  return null;
};

  function EventCard({ event, user, onEdit, onDelete }) {
    const [expanded, setExpanded] = useState(false);
    const [showDeleteModal, setShowDeleteModal] = useState(false);

    const title = event.name || event.title;
      let formattedDate = "";
      if (event.event_date) {
        try {
          formattedDate = new Date(event.event_date).toLocaleDateString("no-NO", {
            year: "numeric", month: "long", day: "numeric"
          });
        } catch {
        formattedDate = event.event_date;
        }
      } else if (event.date) {
        formattedDate = event.date;
      }
      let formattedTime = "";
      if (event.event_time) {
        if (event.event_time.includes("T")) {
          formattedTime = event.event_time.split("T")[1].substring(0, 5);
        } else {
            formattedTime = event.event_time.substring(0, 5);
        }
      }
      const canManage = user?.role === "ADMIN" || ( user?.role === "ORGANIZER" &&  
        user?.organizerId && event.organizer_id && 
        user.organizerId === event.organizer_id
      );

      const handleDelete = (e) => {
        e.stopPropagation();
        // Authorization check: Prevent unauthorized access
        if (!canManage) {
          console.warn("Security Alert: Attempted to delete event without manage permissions.");
          return;
        }
        setShowDeleteModal(true);
      };

      const confirmDelete = () => {
        // Authorization check: Prevent unauthorized access
        if (!canManage) {
          console.warn("Security Alert: Attempted to confirm delete without manage permissions.");
          setShowDeleteModal(false);
          return;
        }
        onDelete(event.id);
        setShowDeleteModal(false);
      };

      const handleEdit = (e) => {
        e.stopPropagation();
        // Authorization check: Prevent unauthorized access
        if (!canManage) {
          console.warn("Security Alert: Attempted to edit event without manage permissions.");
          return;
        }
        onEdit(event);
      };
    return (
      <div className={expanded ? "event-card expanded" : "event-card"}>
      <img
      src={getSafeImageUrl(event.image)}
      alt={title}
      />
      <h2>{title}</h2>
      {event.organizer_name && (
      <span className="event-organizer-tag">
      {event.organizer_name}
      </span>
      )}
      <p className="event-meta">
      📅 {formattedDate} {formattedTime ? `at ${formattedTime}` : ""}
      </p>
      <p className="event-meta">📍 {event.location}</p>
      <div className="event-card-view-action">
        <button
          className="view-btn"
          onClick={() => setExpanded((isExpanded) => !isExpanded)}
          aria-expanded={expanded}
        >
          {expanded ? "Hide details" : "View"}
        </button>
      </div>
      {expanded && (
      <div className="description">
        <p>{event.description}</p>
        <div className="event-card-actions">
          {event.isRegistration && (
            getSafeRegistrationLink(event.registrationLink) ? (
              <a
                className="view-btn"
                href={getSafeRegistrationLink(event.registrationLink)}
                target="_blank"
                rel="noreferrer"
                onClick={(e) => e.stopPropagation()}
              >
                Register
              </a>
            ) : (
              <button className="view-btn" onClick={(e) => e.stopPropagation()}>
                Register
              </button>
            )
          )}
            {canManage && (
              <>
              <button className="edit-btn" onClick={handleEdit}>
              ✏️ Edit
              </button>
              <button className="delete-btn" onClick={handleDelete}>
              🗑️ Delete
              </button>
              </>
            )}
        </div>
      </div>
      )}
      {showDeleteModal && (
        <div className="delete-modal-overlay" onClick={() => setShowDeleteModal(false)}>
          <div className="delete-modal" onClick={(e) => e.stopPropagation()}>
            <h3>Delete event?</h3>
            <p>Are you sure you want to delete "{title}"?</p>
            <div className="delete-modal-actions">
              <button className="cancel-btn" onClick={() => setShowDeleteModal(false)}>
                Cancel
              </button>
              <button className="delete-btn" onClick={confirmDelete}>
                Delete
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
export default EventCard;
