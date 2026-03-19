import React from "react";
import styles from "./index.module.scss";

export function LeftPanel() {
  return (
    <div className={styles.container}>
      <h2 className={styles.title}>Projects</h2>
      <ul className={styles.projectList}>{/* List of projects */}</ul>
    </div>
  );
}
