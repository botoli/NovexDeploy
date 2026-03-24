import { useEffect, useState } from "react";
import { api } from "../../shared/api";
import styles from "./index.module.scss";

export function DeploymentsPage({ logsOnly = false }: { logsOnly?: boolean }) {
  const [projects, setProjects] = useState<any[]>([]);
  const [selectedProject, setSelectedProject] = useState("");
  const [deployments, setDeployments] = useState<any[]>([]);
  const [logs, setLogs] = useState("");

  useEffect(() => {
    api.projects().then((p) => {
      setProjects(p);
      if (p.length > 0) setSelectedProject(p[0].id);
    });
  }, []);

  useEffect(() => {
    if (!selectedProject) return;
    api.deployments(selectedProject).then(setDeployments).catch(() => setDeployments([]));
  }, [selectedProject]);

  const runDeploy = async () => {
    if (!selectedProject) return;
    await api.createDeployment(selectedProject);
    const updated = await api.deployments(selectedProject);
    setDeployments(updated);
  };

  const loadLogs = async (deploymentId: string) => {
    const data = await api.deploymentLogs(deploymentId);
    setLogs(data.logs || "");
  };

  return (
    <div className={styles.root}>
      <h2 className={styles.title}>{logsOnly ? "Logs" : "Deployments"}</h2>
      <div className={styles.toolbar}>
        <select className={styles.select} value={selectedProject} onChange={(e) => setSelectedProject(e.target.value)}>
          {projects.map((p) => (
            <option key={p.id} value={p.id}>
              {p.name}
            </option>
          ))}
        </select>
        {!logsOnly && <button className={styles.button} onClick={runDeploy}>Run deploy</button>}
      </div>
      <div className={styles.list}>
        {deployments.map((d) => (
          <div key={d.id} className={styles.item}>
            <div className={styles.itemHeader}>
              <b>{d.status}</b> - {d.branch}
            </div>
            <button className={styles.button} onClick={() => loadLogs(d.id)}>Open logs</button>
          </div>
        ))}
      </div>
      <pre className={styles.logs}>{logs}</pre>
    </div>
  );
}
