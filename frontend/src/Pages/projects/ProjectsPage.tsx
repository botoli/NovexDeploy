import { useEffect, useState } from "react";
import styles from "./ProjectsPage.module.scss";

export function ProjectsPage() {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const fetchGitHubProfile = async () => {
    setLoading(true);
    try {
      // Используем относительный путь - прокси перенаправит
      const response = await fetch("/auth/me", {
        credentials: "include",
      });

      const data = await response.json();

      if (data.ok) {
        setUser(data.data);
      } else {
        setError(data.message);
      }
    } catch (err) {
      console.error("Fetch error:", err);
      setError("Failed to fetch profile");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchGitHubProfile();
  }, []);

  if (loading) return <div>Loading...</div>;
  if (error)
    return (
      <div>
        <p style={{ color: "red" }}>Error: {error}</p>
        <button onClick={fetchGitHubProfile}>Retry</button>
      </div>
    );
  if (!user)
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
        <pre>{JSON.stringify(user, null, 2)}</pre>

        {/* Аватар */}
        {user.avatar_url && (
          <img
            src={user.avatar_url}
            alt={user.name}
            style={{ width: "100px", borderRadius: "50%" }}
          />
        )}

        {/* Информация */}
        <div>
          <p>
            <strong>Name:</strong> {user.name || "N/A"}
          </p>
          <p>
            <strong>Username:</strong> {user.github_login || "N/A"}
          </p>
          <p>
            <strong>Email:</strong> {user.email || "N/A"}
          </p>
        </div>
      </div>
    </div>
  );
}
