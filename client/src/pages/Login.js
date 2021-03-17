import React from "react";
import "./Login.scss";

function Login() {
  return (
    <form className="login-form">
      <input type="password" name="password" autoFocus />
    </form>
  );
}

export default Login;
