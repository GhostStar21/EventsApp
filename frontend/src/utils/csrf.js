import { API_URL } from "../config.js"

// CSRF Protection utility
export async function getCsrfToken() {
  // Try to get CSRF token from meta tag first
  const metaToken = document.querySelector('meta[name="csrf-token"]')?.getAttribute('content');
  if (metaToken) return metaToken;
  
  // Try to get from cookie
  const cookie = document.cookie.split('; ').find(row => row.startsWith('XSRF-TOKEN='));
  if (cookie) return cookie.split('=')[1];
  
  const response = await fetch(`${API_URL}/v1/csrf-token`, {
    method: "GET",
    credentials: "include",
  });

  if (!response.ok) {
    throw new Error("Could not get CSRF token");
  }

  const data = await response.json();
  return data.token;
}

export default getCsrfToken;
