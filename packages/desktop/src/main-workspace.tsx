import React from "react";
import ReactDOM from "react-dom/client";
import AppWorkspace from "./App-workspace";
import "./styles.css";

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <React.StrictMode>
    <AppWorkspace />
  </React.StrictMode>,
);
