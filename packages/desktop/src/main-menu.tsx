import React from "react";
import ReactDOM from "react-dom/client";
import AppMenu from "./App-menu";
import "./styles.css";

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <React.StrictMode>
    <AppMenu />
  </React.StrictMode>,
);
