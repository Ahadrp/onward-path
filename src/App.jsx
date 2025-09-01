import PrivateRoute from "./PrivateRoute";
import { Routes, Route, Link, useNavigate } from "react-router-dom";
import { useState } from "react";

// Home page with navigation
function Home() {
  return (
    <div>
      <h2>Welcome!</h2>
      <nav>
        <Link to="/register">Register</Link> |{" "}
        <Link to="/login">Login</Link> |{" "}
        <Link to="/profile">Profile</Link>
      </nav>
    </div>
  );
}

// Register Page
function Register() {
  const [email, setEmail] = useState("");
  const [passwd, setPasswd] = useState("");
  const [message, setMessage] = useState("");
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setMessage("");

    try {
      const res = await fetch("/api/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Accept": "application/json",
        },
        body: JSON.stringify({ email, passwd }),
      });

      if (res.ok) {
        setMessage("Registration successful! You can now log in.");
        setTimeout(() => navigate("/login"), 1500);
      } else {
        const data = await res.json();
        setMessage(data.message || "Registration failed.");
      }
    } catch (err) {
      setMessage("Error: " + err.message);
    }
  };

  return (
    <div>
      <h2>Register</h2>
      <form onSubmit={handleSubmit}>
        <div>
          <label>Email: </label>
          <input
            type="email"
            value={email}
            onChange={e => setEmail(e.target.value)}
            required
          />
        </div>
        <div>
          <label>Password: </label>
          <input
            type="password"
            value={passwd}
            onChange={e => setPasswd(e.target.value)}
            required
          />
        </div>
        <button type="submit">Register</button>
      </form>
      <div style={{ color: "red", marginTop: 10 }}>{message}</div>
    </div>
  );
}

// Login Page
function Login() {
const [email, setEmail] = useState("");
  const [passwd, setPasswd] = useState("");
  const [message, setMessage] = useState("");
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setMessage("");

    try {
      const res = await fetch("/api/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Accept": "application/json",
        },
        body: JSON.stringify({ email, passwd }),
        credentials: "include", // <-- Important!
      });

      if (res.ok) {
        setMessage("Login successful!");
        setTimeout(() => navigate("/profile"), 1000);
      } else {
        const data = await res.json();
        setMessage(data.message || "Login failed.");
      }
    } catch (err) {
      setMessage("Error: " + err.message);
    }
  };

  return (
    <div>
      <h2>Login</h2>
      <form onSubmit={handleSubmit}>
        <div>
          <label>Email: </label>
          <input
            type="email"
            value={email}
            onChange={e => setEmail(e.target.value)}
            required
          />
        </div>
        <div>
          <label>Password: </label>
          <input
            type="password"
            value={passwd}
            onChange={e => setPasswd(e.target.value)}
            required
          />
        </div>
        <button type="submit">Login</button>
      </form>
      <div style={{ color: "red", marginTop: 10 }}>{message}</div>
    </div>
  );
}

// Profile Page placeholder
function Profile() {
  return <h2>Profile Page (to be implemented)</h2>;
}

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/register" element={<Register />} />
      <Route path="/login" element={<Login />} />
      <Route
        path="/profile"
        element={
          <PrivateRoute>
            <Profile />
          </PrivateRoute>
        }
      />
    </Routes>
  );
}
