import { useEffect, useState } from "react";
import { api } from "../../shared/api";
import styles from "./index.module.scss";

export function SettingsPage() {
  const [projects, setProjects] = useState<any[]>([]);
  const [selectedProject, setSelectedProject] = useState("");
  const [selectedProjectData, setSelectedProjectData] = useState<any>(null);
  const [deployments, setDeployments] = useState<any[]>([]);
  const [runtime, setRuntime] = useState<any>(null);
  const [telegram, setTelegram] = useState<any>(null);
  const [envVars, setEnvVars] = useState<any[]>([]);
  const [envKey, setEnvKey] = useState("");
  const [envValue, setEnvValue] = useState("");
  const [editingEnvKey, setEditingEnvKey] = useState("");
  const [editingEnvValue, setEditingEnvValue] = useState("");
  const [envDirty, setEnvDirty] = useState(false);
  const [saveMessage, setSaveMessage] = useState("");
  const [botToken, setBotToken] = useState("");
  const [botMode, setBotMode] = useState("polling");
  const [projectForm, setProjectForm] = useState({
    branch: "main",
    root_dir: ".",
    build_command: "",
    start_command: "",
    project_type: "backend",
  });

  const load = async (projectId: string) => {
    try {
      const [projectData, runtimeData, envData, deploymentData] = await Promise.all([
        api.projects().then((list) => list.find((p) => p.id === projectId)),
        api.runtimeStatus(projectId),
        api.envList(projectId),
        api.deployments(projectId),
      ]);
      setSelectedProjectData(projectData || null);
      setProjectForm({
        branch: projectData?.branch || "main",
        root_dir: projectData?.root_dir || ".",
        build_command: projectData?.build_command || "",
        start_command: projectData?.start_command || "",
        project_type: projectData?.project_type || "backend",
      });
      setRuntime(runtimeData);
      setEnvVars(envData);
      setDeployments(deploymentData || []);
      try {
        setTelegram(await api.telegramStatus(projectId));
      } catch {
        setTelegram(null);
      }
    } catch {}
  };

  useEffect(() => {
    api.projects().then((p) => {
      setProjects(p);
      if (p.length > 0) setSelectedProject(p[0].id);
    });
  }, []);

  useEffect(() => {
    if (!selectedProject) return;
    load(selectedProject);
  }, [selectedProject]);

  const runRuntime = async (action: "start" | "stop" | "restart") => {
    if (!selectedProject) return;
    await api.runtimeAction(selectedProject, action);
    await load(selectedProject);
  };

  const saveEnv = async () => {
    if (!selectedProject || !envKey) return;
    await api.envUpsert(selectedProject, envKey, envValue);
    setEnvKey("");
    setEnvValue("");
    setEnvDirty(true);
    await load(selectedProject);
  };

  const startEditEnv = (key: string, value: string) => {
    setEditingEnvKey(key);
    setEditingEnvValue(value);
  };

  const saveEditEnv = async () => {
    if (!selectedProject || !editingEnvKey) return;
    await api.envUpsert(selectedProject, editingEnvKey, editingEnvValue);
    setEditingEnvKey("");
    setEditingEnvValue("");
    setEnvDirty(true);
    await load(selectedProject);
  };

  const removeEnv = async (key: string) => {
    if (!selectedProject) return;
    if (!window.confirm(`Delete env var ${key}?`)) return;
    await api.envDelete(selectedProject, key);
    setEnvDirty(true);
    await load(selectedProject);
  };

  const saveProjectConfig = async () => {
    if (!selectedProject) return;
    await api.updateProject(selectedProject, projectForm);
    setSaveMessage("Project configuration saved");
    await load(selectedProject);
  };

  const saveTelegram = async () => {
    if (!selectedProject || !botToken) return;
    await api.telegramConfig(selectedProject, { mode: botMode, bot_token: botToken });
    setBotToken("");
    await load(selectedProject);
  };

  return (
    <div className={styles.root}>
      <h2 className={styles.title}>Settings</h2>
      <select className={styles.select} value={selectedProject} onChange={(e) => setSelectedProject(e.target.value)}>
        {projects.map((p) => (
          <option key={p.id} value={p.id}>
            {p.name}
          </option>
        ))}
      </select>

      <div className={styles.card}>
        <h3>Project Configuration</h3>
        <div className={styles.row}>
          <input className={styles.input} value={selectedProjectData?.repository || ""} disabled />
          <input className={styles.input} value={projectForm.branch} onChange={(e) => setProjectForm((v) => ({ ...v, branch: e.target.value }))} placeholder="Branch" />
        </div>
        <div className={styles.row}>
          <input className={styles.input} value={projectForm.root_dir} onChange={(e) => setProjectForm((v) => ({ ...v, root_dir: e.target.value }))} placeholder="Root dir" />
          <select className={styles.select} value={projectForm.project_type} onChange={(e) => setProjectForm((v) => ({ ...v, project_type: e.target.value }))}>
            <option value="backend">backend</option>
            <option value="telegram">telegram</option>
          </select>
        </div>
        <div className={styles.row}>
          <input className={styles.input} value={projectForm.build_command} onChange={(e) => setProjectForm((v) => ({ ...v, build_command: e.target.value }))} placeholder="Build command" />
          <input className={styles.input} value={projectForm.start_command} onChange={(e) => setProjectForm((v) => ({ ...v, start_command: e.target.value }))} placeholder="Start command" />
          <button className={styles.btn} onClick={saveProjectConfig}>Save config</button>
        </div>
        {saveMessage && <div>{saveMessage}</div>}
      </div>

      <div className={styles.card}>
        <h3>Runtime</h3>
        <div>State: {runtime?.state || "unknown"}</div>
        <div className={styles.row}>
          <button className={styles.btn} onClick={() => runRuntime("start")}>Start</button>
          <button className={styles.btn} onClick={() => runRuntime("stop")}>Stop</button>
          <button className={styles.btn} onClick={() => runRuntime("restart")}>Restart</button>
        </div>
      </div>

      <div className={styles.card}>
        <h3>Environment Variables</h3>
        {envDirty && <div>Restart runtime or redeploy to apply env changes.</div>}
        <div className={styles.row}>
          <input className={styles.input} placeholder="KEY" value={envKey} onChange={(e) => setEnvKey(e.target.value)} />
          <input className={styles.input} placeholder="VALUE" value={envValue} onChange={(e) => setEnvValue(e.target.value)} />
          <button className={styles.btn} onClick={saveEnv}>Save</button>
        </div>
        <div className={styles.envList}>
          {envVars.map((item) => (
            <div key={item.id} className={styles.envItem}>
              {editingEnvKey === item.key ? (
                <div className={styles.row}>
                  <input className={styles.input} value={editingEnvValue} onChange={(e) => setEditingEnvValue(e.target.value)} />
                  <button className={styles.btn} onClick={saveEditEnv}>Update</button>
                </div>
              ) : (
                <div className={styles.row}>
                  <span>{item.key} = {item.masked_value || "****"}</span>
                  <button className={styles.btn} onClick={() => startEditEnv(item.key, item.value)}>Edit</button>
                  <button className={styles.btn} onClick={() => removeEnv(item.key)}>Delete</button>
                </div>
              )}
            </div>
          ))}
        </div>
      </div>

      <div className={styles.card}>
        <h3>Telegram Bot</h3>
        <div>Current mode: {telegram?.mode || "not configured"}</div>
        <div>Status: {telegram?.is_active ? "active" : "inactive"}</div>
        <div className={styles.row}>
          <select className={styles.select} value={botMode} onChange={(e) => setBotMode(e.target.value)}>
            <option value="polling">polling</option>
            <option value="webhook">webhook</option>
          </select>
          <input
            className={styles.input}
            placeholder="BOT_TOKEN"
            value={botToken}
            onChange={(e) => setBotToken(e.target.value)}
          />
          <button className={styles.btn} onClick={saveTelegram}>Save config</button>
        </div>
      </div>

      <div className={styles.card}>
        <h3>Recent Deployments</h3>
        <div className={styles.envList}>
          {deployments.slice(0, 5).map((d) => (
            <div className={styles.envItem} key={d.id}>
              {d.status} | {d.branch} | {new Date(d.created_at).toLocaleString()}
            </div>
          ))}
          {deployments.length === 0 && <div>No deployments yet</div>}
        </div>
      </div>
    </div>
  );
}
