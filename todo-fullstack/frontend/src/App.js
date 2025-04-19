import React, { useState, useEffect } from "react";
import { register, login, logout, getToken } from "./auth";

const API = "http://localhost:8080/todos";

function App() {
  const [todos, setTodos] = useState([]);
  const [task, setTask] = useState("");
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [authMode, setAuthMode] = useState("login");
  const [error, setError] = useState("");
  const [isLoggedIn, setIsLoggedIn] = useState(!!getToken());

  useEffect(() => {
    if (isLoggedIn) {
      fetch(API, {
        headers: { Authorization: `Bearer ${getToken()}` }
      })
        .then(res => res.json())
        .then(setTodos);
    }
  }, [isLoggedIn]);

  const handleRegister = async e => {
    e.preventDefault();
    try {
      await register(username, password);
      setAuthMode("login");
      setError("");
      alert("Registration successful! Please login.");
    } catch {
      setError("Registration failed");
    }
  };

  const handleLogin = async e => {
    e.preventDefault();
    try {
      await login(username, password);
      setIsLoggedIn(true);
      setError("");
    } catch {
      setError("Login failed");
    }
  };

  const handleLogout = () => {
    logout();
    setIsLoggedIn(false);
    setTodos([]);
    setUsername("");
    setPassword("");
  };

  const addTodo = async e => {
    e.preventDefault();
    const res = await fetch(API, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${getToken()}`
      },
      body: JSON.stringify({ task, done: false })
    });
    const newTodo = await res.json();
    setTodos([...todos, newTodo]);
    setTask("");
  };

  const toggleTodo = async id => {
    const todo = todos.find(t => t.id === id);
    const res = await fetch(`${API}/${id}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${getToken()}`
      },
      body: JSON.stringify({ ...todo, done: !todo.done })
    });
    const updated = await res.json();
    setTodos(todos.map(t => (t.id === id ? updated : t)));
  };

  const deleteTodo = async id => {
    await fetch(`${API}/${id}`, {
      method: "DELETE",
      headers: { Authorization: `Bearer ${getToken()}` }
    });
    setTodos(todos.filter(t => t.id !== id));
  };

  if (!isLoggedIn) {
    return (
      <div style={{ maxWidth: 400, margin: "auto" }}>
        <h1>Todo App</h1>
        <form onSubmit={authMode === "login" ? handleLogin : handleRegister}>
          <input
            value={username}
            onChange={e => setUsername(e.target.value)}
            placeholder="Username"
            required
          />
          <input
            type="password"
            value={password}
            onChange={e => setPassword(e.target.value)}
            placeholder="Password"
            required
          />
          <button type="submit">{authMode === "login" ? "Login" : "Register"}</button>
        </form>
        <button onClick={() => setAuthMode(authMode === "login" ? "register" : "login")}
          style={{ marginTop: 8 }}>
          {authMode === "login" ? "Need an account? Register" : "Already have an account? Login"}
        </button>
        {error && <div style={{ color: "red" }}>{error}</div>}
      </div>
    );
  }

  return (
    <div style={{ maxWidth: 400, margin: "auto" }}>
      <h1>Todo App</h1>
      <button onClick={handleLogout} style={{ float: "right" }}>Logout</button>
      <form onSubmit={addTodo}>
        <input
          value={task}
          onChange={e => setTask(e.target.value)}
          placeholder="New todo"
          required
        />
        <button type="submit">Add</button>
      </form>
      <ul>
        {todos.map(todo => (
          <li key={todo.id}>
            <span
              style={{ textDecoration: todo.done ? "line-through" : "none", cursor: "pointer" }}
              onClick={() => toggleTodo(todo.id)}
            >
              {todo.task}
            </span>
            <button onClick={() => deleteTodo(todo.id)} style={{ marginLeft: 8 }}>
              Delete
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default App;
