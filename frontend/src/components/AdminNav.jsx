import { useNavigate } from "react-router-dom";

function AdminNav({ user }) {
  const navigate = useNavigate();

  if (!user || user.role !== "ADMIN") return null;

  return (
    <button className="admin-link-button" onClick={() => navigate('/admin')}>
      Admin dashboard
    </button>
  );
}

export default AdminNav;
