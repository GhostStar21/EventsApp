import React from "react";
import "../styles/ConfirmModal.css";

export default function ConfirmModal({ open, title, message, requireInput, inputLabel, inputValue, onInputChange, onCancel, onConfirm }) {
  if (!open) return null;

  return (
    <div className="confirm-backdrop">
      <div className="confirm-modal">
        <h3>{title}</h3>
        <p>{message}</p>
        {requireInput && (
          <div style={{marginTop: '0.5rem'}}>
            <label>{inputLabel}</label>
            <input type={inputLabel && inputLabel.toLowerCase().includes('password') ? 'password' : 'text'} value={inputValue} onChange={(e) => onInputChange(e.target.value)} />
          </div>
        )}
        <div className="confirm-actions">
          <button className="cancel" onClick={onCancel}>Cancel</button>
          <button className="confirm" onClick={() => onConfirm(inputValue)}>Confirm</button>
        </div>
      </div>
    </div>
  );
}
