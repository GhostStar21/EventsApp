import "./Popup.css";
import { useState } from "react";

/*
* Handles the login and register pop up displayed for new users.
*/
function PopUp() {
  const [showModal, setShowModal] = useState(true);
  const [showLoginForm, setShowLoginForm] = useState(false);
  const [formData, setFormData] = useState({ username: "", password: "" });
  const [showRegisterForm, setRegisterForm] = useState(false);
  const [registerData, setRegisterData] = useState( {username: "", password: "", email: ""})

  const handleInputChange = (event) => {
    const { name, value } = event.target;
    setFormData((prevData) => ({ ...prevData, [name]: value }));
  };

  const handleLoginSubmit = (event) => {
    event.preventDefault();
    setShowModal(false);
  };

  let loginOverlay = 
  <>
    <h2>Log In</h2>
    <p>Enter your username and password.</p>
    <form className="login-form" onSubmit={handleLoginSubmit}>
      <input
        type="text"
        name="username"
        placeholder="Username"
        value={formData.username}
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
    <form className="login-form" onSubmit={handleLoginSubmit}>
      <input 
        type="username" 
        name="username"
        placeholder="Username"
        value={registerData.username}
        onChange={handleInputChange}
      />
      <input 
        type="password" 
        name="password" 
        placeholder="Password" 
        value={registerData.password} 
        onChange={handleInputChange} 
      />
      <input 
        type="email" 
        name="email" 
        placeholder="Email" 
        value={registerData.email} 
        onChange={handleInputChange} 
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
            {content}
          </div>
        </div>
      )}
    </>
  );
}

export default PopUp;