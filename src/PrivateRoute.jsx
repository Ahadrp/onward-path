import { useEffect, useState } from "react";
import { Navigate } from "react-router-dom";

export default function PrivateRoute({ children }) {
  const [loading, setLoading] = useState(true);
  const [authenticated, setAuthenticated] = useState(false);

  useEffect(() => {
    fetch("/api/Auth-check", { // <-- Implement this endpoint in Go
      method: "GET",
      credentials: "include",
    })
      .then(res => setAuthenticated(res.ok))
      .catch(() => setAuthenticated(false))
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <div>Loading...</div>;
  if (!authenticated) return <Navigate to="/login" replace />;
  return children;
}
