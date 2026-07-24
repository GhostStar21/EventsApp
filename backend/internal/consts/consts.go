package consts

const EventsPath = "/v1/events/"
const OrganizersPath = "/v1/organizers/"
const UsersPath = "/v1/users/"
const MePath = "/v1/me"
const RegisterPath = "/v1/register"
const LoginPath = "/v1/login"
const LogoutPath = "/v1/logout"
const RegisterOrganizerPath = "/v1/register-organizer"
const DemoteOrganizerPath = "/v1/demote-organizer"

type Role string

const (
	RoleUser      Role = "USER"
	RoleOrganizer Role = "ORGANIZER"
	RoleAdmin     Role = "ADMIN"
)
