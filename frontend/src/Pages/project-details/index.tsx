import { useEffect, useMemo, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { api } from "../../shared/api";
import styles from "./index.module.scss";

type RuntimeAction = "start" | "stop" | "restart";

export function ProjectDetailsPage() {
  const { projectId = "" } = useParams();
  const [project, setProject] = useState<any>(null);
  const [runtime, setRuntime] = useState<any>(null);
  const [deployments, setDeployments] = useState<any[]>([]);
  const [envVars, setEnvVars] = useState<any[]>([]);
  const [telegram, setTelegram] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [notice, setNotice] = useState("");

  const [envKey, setEnvKey] = useState("");
  const [envValue, setEnvValue] = useState("");
  const [botMode, setBotMode] = useState("polling");
  const [botToken, setBotToken] = useState("");
  const [webhookURL, setWebhookURL] = useState("");

  const [form, setForm] = useState({
    branch: "main",
    root_dir: ".",
    build_command: "",
    start_command: "",
    project_type: "backend",
  });

  const isTelegramProject = useMemo(
    () => form.project_type === "telegram" || project?.project_type === "telegram",
    [form.project_type, project?.project_type],
  );

  const loadProject = async () => {
    if (!projectId) return;
    setLoading(true);
    setError("");
    try {
      const [projectData, runtimeData, envData, deploymentData] = await Promise.all([
        api.project(projectId),
        api.runtimeStatus(projectId),
        api.envList(projectId),
        api.deployments(projectId),
      ]);
      setProject(projectData);
      setRuntime(runtimeData);
      setEnvVars(envData || []);
      setDeployments(deploymentData || []);
      setForm({
        branch: projectData.branch || "main",
        root_dir: projectData.root_dir || ".",
        build_command: projectData.build_command || "",
        start_command: projectData.start_command || "",
        project_type: projectData.project_type || "backend",
      });
      try {
        const tg = await api.telegramStatus(projectId);
        setTelegram(tg);
        setBotMode(tg.mode || "polling");
        setWebhookURL(tg.webhook || "");
      } catch {
        setTelegram(null);
      }
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to load project");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadProject();
  }, [projectId]);

  const saveProject = async () => {
    if (!projectId) return;
    try {
      await api.updateProject(projectId, form);
      setNotice("Project settings updated");
      await loadProject();
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to save settings");
    }
  };

  const runDeploy = async () => {
    if (!projectId) return;
    await api.createDeployment(projectId);
    await loadProject();
  };

  const runRuntime = async (action: RuntimeAction) => {
    if (!projectId) return;
    await api.runtimeAction(projectId, action);
    await loadProject();
  };

  const saveEnv = async () => {
    if (!projectId || !envKey) return;
    await api.envUpsert(projectId, envKey, envValue);
    setEnvKey("");
    setEnvValue("");
    setNotice("Env updated. Restart runtime to apply changes.");
    await loadProject();
  };

  const deleteEnv = async (key: string) => {
    if (!projectId) return;
    if (!window.confirm(`Delete ${key}?`)) return;
    await api.envDelete(projectId, key);
    setNotice("Env deleted. Restart runtime to apply changes.");
    await loadProject();
  };

  const saveTelegram = async () => {
    if (!projectId || !botToken) return;
    await api.telegramConfig(projectId, {
      mode: botMode,
      bot_token: botToken,
      webhook_url: webhookURL,
    });
    setBotToken("");
    setNotice("Telegram config saved");
    await loadProject();
  };

  if (loading) return <div className={styles.root}>Loading...</div>;
  if (error) return <div className={styles.root}>Error: {error}</div>;
  if (!project) return <div className={styles.root}>Project not found</div>;

  return (
    <div className={styles.root}>
      <div className={styles.header}>
        <div>
          <h2 className={styles.title}>{project.name}</h2>
          <div className={styles.muted}>
            Type: {project.project_type} | Runtime: {runtime?.state || "configured"}
          </div>
        </div>
        <Link to="/projects" className={styles.backLink}>Back to projects</Link>
      </div>

      {notice && <div className={styles.notice}>{notice}</div>}

      <div className={styles.card}>
        <h3>Repository and Build</h3>
        <div className={styles.row}>
          <input className={styles.input} value={project.repository || ""} disabled />
          <input className={styles.input} value={form.branch} onChange={(e) => setForm((v) => ({ ...v, branch: e.target.value }))} placeholder="Branch" />
          <input className={styles.input} value={form.root_dir} onChange={(e) => setForm((v) => ({ ...v, root_dir: e.target.value }))} placeholder="Root dir" />
        </div>
        <div className={styles.row}>
          <input className={styles.input} value={form.build_command} onChange={(e) => setForm((v) => ({ ...v, build_command: e.target.value }))} placeholder="Build command" />
          <input className={styles.input} value={form.start_command} onChange={(e) => setForm((v) => ({ ...v, start_command: e.target.value }))} placeholder="Start command" />
          <select className={styles.select} value={form.project_type} onChange={(e) => setForm((v) => ({ ...v, project_type: e.target.value }))}>
            <option value="backend">backend</option>
            <option value="telegram">telegram</option>
          </select>
          <button className={styles.button} onClick={saveProject}>Save</button>
        </div>
      </div>

      <div className={styles.card}>
        <h3>Operations</h3>
        <div className={styles.row}>
          <button className={styles.button} onClick={runDeploy}>Deploy now</button>
          <button className={styles.button} onClick={() => runRuntime("start")}>Start</button>
          <button className={styles.button} onClick={() => runRuntime("stop")}>Stop</button>
          <button className={styles.button} onClick={() => runRuntime("restart")}>Restart</button>
        </div>
      </div>

      <div className={styles.grid}>
        <div className={styles.card}>
          <h3>Environment Variables</h3>
          <div className={styles.row}>
            <input className={styles.input} placeholder="KEY" value={envKey} onChange={(e) => setEnvKey(e.target.value)} />
            <input className={styles.input} placeholder="VALUE" value={envValue} onChange={(e) => setEnvValue(e.target.value)} />
            <button className={styles.button} onClick={saveEnv}>Save env</button>
          </div>
          <div className={styles.list}>
            {envVars.map((item) => (
              <div key={item.id} className={styles.itemRow}>
                <span>{item.key} = {item.masked_value || "****"}</span>
                <button className={styles.button} onClick={() => deleteEnv(item.key)}>Delete</button>
              </div>
            ))}
            {envVars.length === 0 && <div className={styles.muted}>No env vars</div>}
          </div>
        </div>

        <div className={styles.card}>
          <h3>Recent Deployments</h3>
          <div className={styles.list}>
            {deployments.slice(0, 10).map((d) => (
              <div key={d.id} className={styles.itemRow}>
                <span>{d.status} | {d.branch}</span>
                <span className={styles.muted}>{new Date(d.created_at).toLocaleString()}</span>
              </div>
            ))}
            {deployments.length === 0 && <div className={styles.muted}>No deployments yet</div>}
          </div>
        </div>
      </div>

      {isTelegramProject && (
        <div className={styles.card}>
          <h3>Telegram</h3>
          <div className={styles.muted}>
            Current: {telegram?.mode || "not configured"} | {telegram?.is_active ? "active" : "inactive"}
          </div>
          <div className={styles.row}>
            <select className={styles.select} value={botMode} onChange={(e) => setBotMode(e.target.value)}>
              <option value="polling">polling</option>
              <option value="webhook">webhook</option>
            </select>
            <input className={styles.input} placeholder="BOT_TOKEN" value={botToken} onChange={(e) => setBotToken(e.target.value)} />
            {botMode === "webhook" && (
              <input className={styles.input} placeholder="Webhook URL" value={webhookURL} onChange={(e) => setWebhookURL(e.target.value)} />
            )}
            <button className={styles.button} onClick={saveTelegram}>Save Telegram config</button>
          </div>
        </div>
      )}
    </div>
  );
}

