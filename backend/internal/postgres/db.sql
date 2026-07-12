CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    is_exclusive BOOLEAN NOT NULL DEFAULT false,
    event_date DATE NOT NULL,
    event_time TIME NOT NULL,
    location TEXT NOT NULL,
    description TEXT
);

CREATE TABLE organizers (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    org_number INT UNIQUE NOT NULL
);

CREATE TABLE event_organizers (
    event_id INT REFERENCES events(id) ON DELETE CASCADE,
    organizer_id INT REFERENCES organizers(id) ON DELETE CASCADE,
    PRIMARY KEY (event_id, organizer_id)
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE,
    password_hash TEXT,
    role VARCHAR(20) NOT NULL DEFAULT 'USER'
);

CREATE TABLE user_interested_organizers (
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    organizer_id INT REFERENCES organizers(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, organizer_id)
);


INSERT INTO organizers (name, org_number) VALUES
('Tech Trondheim AS', 1001),
('Nordic Music Org', 1002),
('Startup Hub Norway', 1003),
('Arctic Events', 1004);

INSERT INTO users (name) VALUES
('Alice'),
('Bob'),
('Charlie'),
('Diana'),
('Eirik');

INSERT INTO events (name, is_exclusive, event_date, event_time, location, description) VALUES
('AI & Future Summit', true, '2026-07-10', '09:00', 'Trondheim Spektrum', 'A conference about AI, ML and future tech.'),
('Nordic Music Festival', false, '2026-08-02', '18:00', 'Bergen Arena', 'Annual music festival featuring Nordic artists.'),
('Startup Pitch Night', false, '2026-07-15', '17:30', 'Trondheim Workspaces', 'Pitch your startup to investors.'),
('Arctic Innovation Day', true, '2026-09-01', '10:00', 'Tromsø Conference Center', 'Exclusive innovation showcase event.'),
('Open Tech Meetup', false, '2026-06-25', '16:00', 'NTNU Campus', 'Casual tech meetup for developers.');

INSERT INTO event_organizers (event_id, organizer_id) VALUES
(1, 1), -- AI Summit → Tech Trondheim AS
(1, 3), -- AI Summit → Startup Hub Norway
(2, 2), -- Music Festival → Nordic Music Org
(3, 3), -- Startup Night → Startup Hub Norway
(4, 4), -- Arctic Innovation → Arctic Events
(4, 1), -- Arctic Innovation → Tech Trondheim AS
(5, 1); -- Open Meetup → Tech Trondheim AS

INSERT INTO user_interested_organizers (user_id, organizer_id) VALUES
(1, 1), -- Alice → Tech Trondheim AS
(1, 3), -- Alice → Startup Hub Norway
(2, 2), -- Bob → Nordic Music Org
(3, 1), -- Charlie → Tech Trondheim AS
(3, 4), -- Charlie → Arctic Events
(4, 3), -- Diana → Startup Hub Norway
(5, 2), -- Eirik → Nordic Music Org
(5, 1); -- Eirik → Tech Trondheim AS