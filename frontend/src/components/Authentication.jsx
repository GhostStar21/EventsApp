import "../styles/Authentication.css";
import { useState } from "react";
import { useNavigate } from "react-router-dom";

/**
* Handles the login and register pop up displayed for new users.
*
* @param {setUser} - A function parameter for the user being authenticated
* @return {JSX Element} - An authentication UI (Login/Register) 
*/
function Authentication({ setUser }) {
  const navigate = useNavigate();
  const [showModal, setShowModal] = useState(true);
  const [showLoginForm, setShowLoginForm] = useState(false);
  const [formData, setFormData] = useState({ email: "", password: "" });
  const [showRegisterForm, setRegisterForm] = useState(false);
  const [registerData, setRegisterData] = useState({ username: "", password: "", email: "" });
  const [statusMessage, setStatusMessage] = useState("");
  const [isErrorMessage, setIsErrorMessage] = useState(false);

  const handleInputChange = (event) => {
    const { name, value } = event.target;
    setFormData((prevData) => ({ ...prevData, [name]: value }));
  };

  const handleRegisterInputChange = (event) => {
    const { name, value } = event.target;
    setRegisterData((prevData) => ({ ...prevData, [name]: value }));
  };

  const handleLogin = async (event) => {
    event.preventDefault();
    setStatusMessage("");
    setIsErrorMessage(false);

    console.log("Attempting login with:", formData.email, formData.password);

    try {
      const response = await fetch("http://localhost:8080/v1/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          email: formData.email,
          password: formData.password,
        }),
      });

      const data = await response.json().catch(() => null);
      console.log("Response status:", response.status);

      if (response.ok) {
        console.log("Login successful!");
        console.log("Token:", data?.token);
        const token = data?.token || "";
        if (!token) {
          throw new Error("Missing token from login response");
        }

        localStorage.setItem("token", token);

        const meResponse = await fetch("http://localhost:8080/v1/me", {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (!meResponse.ok) {
          throw new Error("Failed to load user profile");
        }

        const user = await meResponse.json();
        setUser(user);
        setStatusMessage("Login successful!");
        setIsErrorMessage(false);
        navigate("/events");
        setShowModal(false);
      } else {
        const message = data?.message || "Incorrect email or password. Please try again.";
        console.log("Login failed:", message);
        setStatusMessage(message);
        setIsErrorMessage(true);
      }
    } catch (error) {
      console.error("Server error:", error);
      setStatusMessage("Could not reach the server. Make sure the backend is running.");
      setIsErrorMessage(true);
    }
  };

  const handleRegisterSubmit = async (event) => {
    event.preventDefault();
    setStatusMessage("");
    setIsErrorMessage(false);

    try {
      const response = await fetch("http://localhost:8080/v1/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          name: registerData.username,
          email: registerData.email,
          password: registerData.password,
        }),
      });

      if (response.ok) {
        setStatusMessage("Account created. You can now log in.");
        setRegisterForm(false);
        setShowLoginForm(true);
      } else {
        const data = await response.json().catch(() => null);
        setStatusMessage(data?.message || "Registration failed.");
        setIsErrorMessage(true);
      }
    } catch (error) {
      console.error("Registration error:", error);
      setStatusMessage("Could not reach the server. Make sure the backend is running.");
      setIsErrorMessage(true);
    }
  };

  let loginOverlay = 
  <>
    <h2>Log In</h2>
    <p>Enter your username and password.</p>
    <form className="login-form" onSubmit={handleLogin}>
      <input
        type="text"
        name="email"
        placeholder="Email"
        value={formData.email}
        onChange={handleInputChange}
      />
      <input
        type="password"
        name="password"
        placeholder="Password"
        value={formData.password}
        onChange={handleInputChange}
      />
      <button type="submit" className="primary-btn">Log In</button>
      <button
        type="button"
        className="secondary-btn"
        onClick={() => setShowLoginForm(false)}
      >
        Back
      </button>
    </form>
  </>

  let registerOverlay = 
  <>
    <h2>Register</h2>
    <p>Enter your details.</p>
    <form className="login-form" onSubmit={handleRegisterSubmit}>
      <input 
        type="username" 
        name="username"
        placeholder="Username"
        value={registerData.username}
        onChange={handleRegisterInputChange}
      />
      <input 
        type="password" 
        name="password" 
        placeholder="Password" 
        value={registerData.password} 
        onChange={handleRegisterInputChange} 
      />
      <input 
        type="email" 
        name="email" 
        placeholder="Email" 
        value={registerData.email} 
        onChange={handleRegisterInputChange} 
      />
      <button type="submit" className="primary-btn">Register</button>
      <button
        type="button"
        className="secondary-btn"
        onClick={() => setRegisterForm(false)}
      >
        Back
      </button>
    </form>
    </>

  let frontOverlay =  
  <>
    <h2>Welcome!</h2>
    <p>Please log in or register to continue.</p>
    <button type="button" onClick={() => setShowLoginForm(true)}>
      Log In
    </button>
    <button type="button" onClick={() => setRegisterForm(true)}>Register</button>
  </> 

  let content = frontOverlay;

  if (showLoginForm) {
    content = loginOverlay;
  }
  else if (showRegisterForm) {
    content = registerOverlay;
  }
  else {
    content = frontOverlay;
  }

  return (
    <>
      {showModal && (
        <div className="overlay">
          <div className="modal">
            {statusMessage && (
              <p className={`status-message ${isErrorMessage ? "error" : "success"}`}>
                {statusMessage}
              </p>
            )}
            {content}
          </div>
        </div>
      )}
    </>
  );
}

export default Authentication;