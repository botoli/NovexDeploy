import { useEffect, useState } from "react";
import styles from "./index.module.scss";

export function OverviewPage() {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const fetchGitHubProfile = async () => {
    setLoading(true);
    try {
      const response = await fetch("http://localhost:8888/auth/me", {
        credentials: "include", // важно для cookie
      });
      const data = await response.json();

      if (data.ok) {
        setUser(data.data);
      } else {
        setError(data.message);
      }
    } catch (err) {
      setError("Failed to fetch profile");
    } finally {
      setLoading(false);
    }
  };

  // Загружаем при монтировании
  useEffect(() => {
    fetchGitHubProfile();
  }, []);

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;
  if (!user) return <div>Please login</div>;

  return (
    <div className={styles.root}>
      <div
        style={{
          padding: "20px",
          border: "1px solid #ccc",
          borderRadius: "8px",
        }}
      >
        <h2>GitHub Profile</h2>

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
            <strong>Username:</strong> {user.github_login}
          </p>
          <p>
            <strong>Email:</strong> {user.email || "N/A"}
          </p>
          <p>
            <strong>GitHub ID:</strong> {user.github_id}
          </p>
          <p>
            <strong>Last Login:</strong>{" "}
            {new Date(user.last_login_at).toLocaleString()}
          </p>
        </div>
      </div>
    </div>
  );
}
