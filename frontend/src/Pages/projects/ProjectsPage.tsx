import styles from "./ProjectsPage.module.scss";
import { GithubUser } from "../../Store/GithubInfo.store";
import { useEffect, useState } from "react";
import { api } from "../../shared/api";

export function ProjectsPage() {
  const [projects, setProjects] = useState<any[]>([]);
  const [name, setName] = useState("");
  const [projectType, setProjectType] = useState("backend");
  const [rootDir, setRootDir] = useState(".");
  const [loading, setLoading] = useState(false);

  const loadProjects = async () => {
    try {
      const data = await api.projects();
      setProjects(data);
    } catch {
      setProjects([]);
    }
  };

  useEffect(() => {
    loadProjects();
  }, []);

  const createProject = async () => {
    if (!name.trim()) return;
    setLoading(true);
    try {
      await api.createProject({
        name,
        description: "",
        framework: projectType === "telegram" ? "telegram" : "node",
        project_type: projectType,
        build_command: "npm ci && npm run build",
        root_dir: rootDir,
        output_dir: "dist",
        start_command: "npm start",
      });
      setName("");
      setRootDir(".");
      await loadProjects();
    } finally {
      setLoading(false);
    }
  };

  if (GithubUser.loading) return <div>Loading...</div>;
  if (GithubUser.error)
    return (
      <div className={styles.authBox}>
        <p style={{ color: "red" }}>Error: {GithubUser.error}</p>
      </div>
    );
  if (!GithubUser.user)
    return (
      <div className={styles.authBox}>
        <p>Please login</p>
        <a href="/v1/auth/github/login">
          <button className={styles.button}>Login with GitHub</button>
        </a>
      </div>
    );

  return (
    <div className={styles.Page}>
      <div className={styles.root}>
        <h2 className={styles.title}>Projects</h2>
        <div className={styles.creatorRow}>
          <input
            className={styles.input}
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="Project name"
          />
          <input
            className={styles.input}
            value={rootDir}
            onChange={(e) => setRootDir(e.target.value)}
            placeholder="Root dir (.)"
          />
          <select className={styles.select} value={projectType} onChange={(e) => setProjectType(e.target.value)}>
            <option value="backend">Backend Service</option>
            <option value="telegram">Telegram Bot</option>
          </select>
          <button className={styles.button} disabled={loading} onClick={createProject}>
            {loading ? "Creating..." : "Create project"}
          </button>
        </div>
        <div className={styles.grid}>
          {projects.map((p) => (
            <div key={p.id} className={styles.card}>
              <b>{p.name}</b>
              <div className={styles.meta}>Type: {p.project_type || "backend"}</div>
              <div className={styles.meta}>Repo: {p.repository || "not connected"}</div>
              <div className={styles.meta}>Branch: {p.branch || "main"}</div>
              <div className={styles.meta}>Root: {p.root_dir || "."}</div>
              <div className={styles.meta}>Runtime: {p.runtime_state || "stopped"}</div>
            </div>
          ))}
          {projects.length === 0 && <div className={styles.empty}>No projects yet</div>}
        </div>
        <div className={styles.userInfo}>
          <b>Connected GitHub user:</b>
          <div>
            {GithubUser.user.github_login} ({GithubUser.user.email || "no email"})
          </div>
        </div>
      </div>
    </div>
  );
}
