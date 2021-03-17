import React, { memo, useContext } from "react";
import { Redirect } from "@reach/router";
import AuthContext from "@shared/contexts/Auth";

function ProtectedComponent({ component: Component, ...rest }) {
  const loggedIn = useContext(AuthContext);

  // Not using `noThrow` causing UI to break. Don't know why
  if (!loggedIn) return <Redirect from="" to="login" noThrow />;

  return <Component {...rest} />;
}

export default memo(ProtectedComponent);
