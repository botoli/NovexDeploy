import { useEffect, useState } from "react";
import styles from "./index.module.scss";
import { api } from "../../shared/api";

export function DashboardPage() {
  const [projects, setProjects] = useState<any[]>([]);
  const [deployments, setDeployments] = useState<any[]>([]);

  useEffect(() => {
    const load = async () => {
      const p = await api.projects();
      setProjects(p);
      const deploymentsList = await Promise.all(
        p.slice(0, 5).map(async (project) => {
          try {
            const d = await api.deployments(project.id);
            return d.slice(0, 2).map((item) => ({ ...item, projectName: project.name }));
          } catch {
            return [];
          }
        }),
      );
      setDeployments(deploymentsList.flat());
    };
    load();
  }, []);
  return (
    <div className={styles.Page}>
      <h1 className={styles.title}>Platform Overview</h1>
      <div className={styles.grid}>
        <div className={styles.card}>
          <h3>Projects</h3>
          <div className={styles.list}>
            {projects.map((p) => (
              <div key={p.id} className={styles.row}>
                <b>{p.name}</b> - {p.project_type || "backend"} - {p.runtime_state || "configured"}
              </div>
            ))}
            {projects.length === 0 && <div>No projects</div>}
          </div>
        </div>
        <div className={styles.card}>
          <h3>Recent Deployments</h3>
          <div className={styles.list}>
            {deployments.map((d) => (
              <div key={d.id} className={styles.row}>
                <b>{d.projectName}</b> - {d.status} - {d.branch}
              </div>
            ))}
            {deployments.length === 0 && <div>No deployments yet</div>}
          </div>
        </div>
      </div>
    </div>
  );
}
