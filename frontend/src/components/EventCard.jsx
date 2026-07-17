import { useState } from "react";
import "../styles/EventCard.css";

/**
 * Displays an event card, and its description based on clicks.
 * 
 * @param {event} - An event fetched from the backend side
 * @returns JSX element - The rendered event card
 */
function EventCard({ event }) {
    const [expanded, setExpanded] = useState(false);
    return (
        <div
            className={expanded ? "event-card expanded" : "event-card"}
            onClick={() => setExpanded(!expanded)}
        >
            <img
                src={event.image}
                alt={event.title}
            />
            <h2>{event.title}</h2>
            <p>{event.date}</p>
            <p>{event.location}</p>
            {expanded && (
                <div className="description">
                    <p>
                        {event.description}
                    </p>
                    <button>
                        View event
                    </button>
                </div>
            )}
        </div>
    );
}

export default EventCard;