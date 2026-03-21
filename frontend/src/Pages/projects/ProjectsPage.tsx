import { useEffect, useState } from "react";
import styles from "./ProjectsPage.module.scss";
import { GithubUser } from "../../Store/GithubInfo.store";

export function ProjectsPage() {
  if (GithubUser.loading) return <div>Loading...</div>;
  if (GithubUser.error)
    return (
      <div>
        <p style={{ color: "red" }}>Error: {GithubUser.error}</p>
      </div>
    );
  if (!GithubUser.user)
    return (
      <div>
        <p>Please login</p>
        <a href="http://localhost:8888/auth/github/login">
          <button>Login with GitHub</button>
        </a>
      </div>
    );

  return (
    <div className={styles.root}>
      <div style={{ padding: "20px", border: "1px solid #ccc" }}>
        <h2>GitHub Profile</h2>

        {/* Проверяем что пришло */}
        <pre>{JSON.stringify(GithubUser.user, null, 2)}</pre>

        {/* Аватар */}
        {GithubUser.user.avatar_url && (
          <img
            src={GithubUser.user.avatar_url}
            alt={GithubUser.user.name}
            style={{ width: "100px", borderRadius: "50%" }}
          />
        )}

        {/* Информация */}
        <div>
          <p>
            <strong>Name:</strong> {GithubUser.user.name || "N/A"}
          </p>
          <p>
            <strong>Username:</strong> {GithubUser.user.github_login || "N/A"}
          </p>
          <p>
            <strong>Email:</strong> {GithubUser.user.email || "N/A"}
          </p>
        </div>
      </div>
    </div>
  );
}
