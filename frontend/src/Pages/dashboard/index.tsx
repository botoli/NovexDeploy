import { useEffect, useState } from "react";
import styles from "./index.module.scss";
import { api } from "../../shared/api";
import { useNavigate } from "react-router-dom";
import { Rocket, Server, Bot, Activity, Plus } from "lucide-react";

export function DashboardPage() {
  const navigate = useNavigate();
  const [projects, setProjects] = useState<any[]>([]);
  const [deployments, setDeployments] = useState<any[]>([]);
  const [runtimeMap, setRuntimeMap] = useState<Record<string, string>>({});
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const load = async () => {
      setLoading(true);
      try {
        const p = await api.projects();
        setProjects(p);
        const deploymentsList = await Promise.all(
          p.slice(0, 8).map(async (project) => {
            try {
              const d = await api.deployments(project.id);
              return d.slice(0, 3).map((item) => ({ ...item, projectName: project.name, projectId: project.id }));
            } catch {
              return [];
            }
          }),
        );
        setDeployments(
          deploymentsList
            .flat()
            .sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
            .slice(0, 12),
        );

        const runtimePairs = await Promise.all(
          p.slice(0, 8).map(async (project) => {
            try {
              const runtimeData = await api.runtimeStatus(project.id);
              return [project.id, runtimeData?.state || "configured"] as const;
            } catch {
              return [project.id, project.runtime_state || "configured"] as const;
            }
          }),
        );
        setRuntimeMap(Object.fromEntries(runtimePairs));
      } finally {
        setLoading(false);
      }
    };
    load();
  }, []);

  const activeDeployments = deployments.filter((d) => d.status === "deploying" || d.status === "building");
  const successDeployments = deployments.filter((d) => d.status === "ready").length;
  const failedDeployments = deployments.filter((d) => d.status === "failed").length;
  const runningProjects = projects.filter((p) => (runtimeMap[p.id] || p.runtime_state) === "running").length;
  const backendProjects = projects.filter((p) => (p.project_type || "backend") === "backend");
  const telegramProjects = projects.filter((p) => (p.project_type || "backend") === "telegram");

  const formatStatus = (status: string) => {
    if (status === "ready") return "success";
    if (status === "failed") return "failed";
    if (status === "building" || status === "deploying") return "building";
    return status || "unknown";
  };

  const runtimeLabel = (value: string) => value || "configured";

  return (
    <div className={styles.Page}>
      <div className={styles.topBar}>
        <div>
          <h1 className={styles.title}>Control Center</h1>
          <p className={styles.subtitle}>Backend and Telegram deployment platform</p>
        </div>
        <div className={styles.topActions}>
          <button className={styles.actionGhost} onClick={() => navigate("/deploy")}>
            <Plus size={16} /> Import repo
          </button>
          <button className={styles.actionPrimary} onClick={() => navigate("/projects")}>
            Open projects
          </button>
        </div>
      </div>

      <div className={styles.metrics}>
        <div className={styles.metricCard}>
          <div className={styles.metricHeader}><Server size={16} /> Projects</div>
          <div className={styles.metricValue}>{projects.length}</div>
          <div className={styles.metricMeta}>{backendProjects.length} backend • {telegramProjects.length} telegram</div>
        </div>
        <div className={styles.metricCard}>
          <div className={styles.metricHeader}><Activity size={16} /> Runtime</div>
          <div className={styles.metricValue}>{runningProjects}</div>
          <div className={styles.metricMeta}>running services</div>
        </div>
        <div className={styles.metricCard}>
          <div className={styles.metricHeader}><Rocket size={16} /> Deployments</div>
          <div className={styles.metricValue}>{activeDeployments.length}</div>
          <div className={styles.metricMeta}>active now</div>
        </div>
        <div className={styles.metricCard}>
          <div className={styles.metricHeader}><Bot size={16} /> Health</div>
          <div className={styles.metricMeta}>Success: {successDeployments}</div>
          <div className={styles.metricMeta}>Failed: {failedDeployments}</div>
        </div>
      </div>

      <div className={styles.grid}>
        <div className={styles.card}>
          <div className={styles.cardHeader}>
            <h3>Active deployments</h3>
            <button className={styles.linkBtn} onClick={() => navigate("/deployments")}>View all</button>
          </div>
          <div className={styles.list}>
            {activeDeployments.slice(0, 6).map((d) => (
              <div key={d.id} className={styles.row}>
                <div>
                  <b>{d.projectName}</b>
                  <div className={styles.meta}>{d.branch}</div>
                </div>
                <span className={`${styles.status} ${styles[formatStatus(d.status)]}`}>{formatStatus(d.status)}</span>
              </div>
            ))}
            {activeDeployments.length === 0 && <div className={styles.empty}>No active deployments</div>}
          </div>
        </div>

        <div className={styles.card}>
          <div className={styles.cardHeader}>
            <h3>Quick actions</h3>
          </div>
          <div className={styles.quickGrid}>
            <button className={styles.quickBtn} onClick={() => navigate("/projects")}>New backend project</button>
            <button className={styles.quickBtn} onClick={() => navigate("/projects")}>New telegram bot</button>
            <button className={styles.quickBtn} onClick={() => navigate("/deploy")}>Import from GitHub</button>
            <button className={styles.quickBtn} onClick={() => navigate("/settings")}>Platform settings</button>
          </div>
        </div>

        <div className={styles.card}>
          <div className={styles.cardHeader}>
            <h3>Projects</h3>
            <button className={styles.linkBtn} onClick={() => navigate("/projects")}>Manage</button>
          </div>
          <div className={styles.projectGrid}>
            {projects.slice(0, 8).map((p) => (
              <div key={p.id} className={styles.projectCard} onClick={() => navigate(`/projects/${p.id}`)}>
                <div className={styles.projectName}>{p.name}</div>
                <div className={styles.meta}>{p.project_type || "backend"}</div>
                <div className={styles.meta}>runtime: {runtimeLabel(runtimeMap[p.id] || p.runtime_state)}</div>
              </div>
            ))}
            {projects.length === 0 && <div className={styles.empty}>No projects yet</div>}
          </div>
        </div>

        <div className={styles.card}>
          <div className={styles.cardHeader}>
            <h3>Recent deployments</h3>
            <button className={styles.linkBtn} onClick={() => navigate("/logs")}>Open logs</button>
          </div>
          <div className={styles.list}>
            {deployments.slice(0, 10).map((d) => (
              <div key={d.id} className={styles.row}>
                <div>
                  <b>{d.projectName}</b>
                  <div className={styles.meta}>{d.branch} • {new Date(d.created_at).toLocaleString()}</div>
                </div>
                <span className={`${styles.status} ${styles[formatStatus(d.status)]}`}>{formatStatus(d.status)}</span>
              </div>
            ))}
            {deployments.length === 0 && !loading && <div className={styles.empty}>No deployments yet</div>}
            {loading && <div className={styles.empty}>Loading dashboard...</div>}
          </div>
        </div>
      </div>
    </div>
  );
}
