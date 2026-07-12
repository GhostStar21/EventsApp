import "./index.css"
import { useState } from "react";
import PopUp from "./components/Popup.jsx";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom"

function App() {
  const [showModal, setShowModal] = useState(true);

  return (
    <>
      {/* Website */}
      <div className="website">
        <PopUp />
        
      </div>

      {/* Login popup */}
    </>
  );
}

export default App;