import styles from "./ProjectsPage.module.scss";
import { GithubUser } from "../../Store/GithubInfo.store";
import { useEffect, useMemo, useState } from "react";
import { api } from "../../shared/api";
import { useNavigate } from "react-router-dom";
import { Search, Rocket, Play, Square, RefreshCw } from "lucide-react";

export function ProjectsPage() {
  const navigate = useNavigate();
  const [projects, setProjects] = useState<any[]>([]);
  const [name, setName] = useState("");
  const [projectType, setProjectType] = useState("backend");
  const [rootDir, setRootDir] = useState(".");
  const [search, setSearch] = useState("");
  const [filterType, setFilterType] = useState("all");
  const [filterRuntime, setFilterRuntime] = useState("all");
  const [loading, setLoading] = useState(false);
  const [busyProjectId, setBusyProjectId] = useState("");
  const [error, setError] = useState("");

  const loadProjects = async () => {
    try {
      const data = await api.projects();
      setProjects(data);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to load projects");
      setProjects([]);
    }
  };

  useEffect(() => {
    loadProjects();
  }, []);

  const createProject = async () => {
    if (!name.trim()) return;
    setLoading(true);
    setError("");
    try {
      await api.createProject({
        name,
        description: "",
        framework: projectType === "telegram" ? "telegram" : "node",
        project_type: projectType,
        build_command: "npm ci",
        root_dir: rootDir,
        output_dir: ".",
        start_command: "npm start",
      });
      setName("");
      setRootDir(".");
      await loadProjects();
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to create project");
    } finally {
      setLoading(false);
    }
  };

  const filteredProjects = useMemo(() => {
    return projects.filter((p) => {
      const text = `${p.name} ${p.repository || ""} ${p.branch || ""}`.toLowerCase();
      const searchOk = text.includes(search.toLowerCase());
      const type = p.project_type || "backend";
      const runtime = p.runtime_state || "configured";
      const typeOk = filterType === "all" || type === filterType;
      const runtimeOk = filterRuntime === "all" || runtime === filterRuntime;
      return searchOk && typeOk && runtimeOk;
    });
  }, [projects, search, filterType, filterRuntime]);

  const counts = useMemo(() => {
    return {
      total: projects.length,
      backend: projects.filter((p) => (p.project_type || "backend") === "backend").length,
      telegram: projects.filter((p) => (p.project_type || "backend") === "telegram").length,
      running: projects.filter((p) => (p.runtime_state || "configured") === "running").length,
    };
  }, [projects]);

  const runProjectAction = async (projectId: string, action: "deploy" | "start" | "stop" | "restart") => {
    setBusyProjectId(projectId + action);
    setError("");
    try {
      if (action === "deploy") {
        await api.createDeployment(projectId);
      } else {
        await api.runtimeAction(projectId, action);
      }
      await loadProjects();
    } catch (e) {
      setError(e instanceof Error ? e.message : "Action failed");
    } finally {
      setBusyProjectId("");
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
        <div className={styles.topBar}>
          <div>
            <h2 className={styles.title}>Projects</h2>
            <p className={styles.subtitle}>Manage backend services and Telegram bots from one place</p>
          </div>
        </div>

        <div className={styles.metrics}>
          <div className={styles.metricCard}>
            <div className={styles.metricValue}>{counts.total}</div>
            <div className={styles.metricLabel}>Total projects</div>
          </div>
          <div className={styles.metricCard}>
            <div className={styles.metricValue}>{counts.backend}</div>
            <div className={styles.metricLabel}>Backend</div>
          </div>
          <div className={styles.metricCard}>
            <div className={styles.metricValue}>{counts.telegram}</div>
            <div className={styles.metricLabel}>Telegram</div>
          </div>
          <div className={styles.metricCard}>
            <div className={styles.metricValue}>{counts.running}</div>
            <div className={styles.metricLabel}>Running</div>
          </div>
        </div>

        <div className={styles.creatorCard}>
          <div className={styles.creatorTitle}>Create project</div>
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
        </div>

        <div className={styles.filters}>
          <div className={styles.searchWrap}>
            <Search size={16} />
            <input
              className={styles.searchInput}
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              placeholder="Search by project, repo or branch"
            />
          </div>
          <select className={styles.select} value={filterType} onChange={(e) => setFilterType(e.target.value)}>
            <option value="all">All types</option>
            <option value="backend">Backend</option>
            <option value="telegram">Telegram</option>
          </select>
          <select className={styles.select} value={filterRuntime} onChange={(e) => setFilterRuntime(e.target.value)}>
            <option value="all">All runtime states</option>
            <option value="running">running</option>
            <option value="stopped">stopped</option>
            <option value="configured">configured</option>
            <option value="failed">failed</option>
          </select>
        </div>

        {error && <div className={styles.error}>{error}</div>}

        <div className={styles.grid}>
          {filteredProjects.map((p) => (
            <div key={p.id} className={styles.card}>
              <div className={styles.cardHead}>
                <b>{p.name}</b>
                <span className={`${styles.badge} ${styles[p.runtime_state || "configured"]}`}>{p.runtime_state || "configured"}</span>
              </div>
              <div className={styles.meta}>Type: {p.project_type || "backend"}</div>
              <div className={styles.meta}>Repo: {p.repository || "not connected"}</div>
              <div className={styles.meta}>Branch: {p.branch || "main"}</div>
              <div className={styles.meta}>Root: {p.root_dir || "."}</div>

              <div className={styles.actions}>
                <button className={styles.buttonSecondary} onClick={() => navigate(`/projects/${p.id}`)}>
                  Open
                </button>
                <button className={styles.iconBtn} onClick={() => runProjectAction(p.id, "deploy")} disabled={busyProjectId === p.id + "deploy"}>
                  <Rocket size={15} />
                </button>
                <button className={styles.iconBtn} onClick={() => runProjectAction(p.id, "start")} disabled={busyProjectId === p.id + "start"}>
                  <Play size={15} />
                </button>
                <button className={styles.iconBtn} onClick={() => runProjectAction(p.id, "stop")} disabled={busyProjectId === p.id + "stop"}>
                  <Square size={15} />
                </button>
                <button className={styles.iconBtn} onClick={() => runProjectAction(p.id, "restart")} disabled={busyProjectId === p.id + "restart"}>
                  <RefreshCw size={15} />
                </button>
              </div>
            </div>
          ))}
          {filteredProjects.length === 0 && <div className={styles.empty}>No projects found</div>}
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
