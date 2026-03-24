import { useEffect, useState } from "react";
import { api } from "../../shared/api";
import styles from "./index.module.scss";

export function SettingsPage() {
  const [projects, setProjects] = useState<any[]>([]);
  const [selectedProject, setSelectedProject] = useState("");
  const [runtime, setRuntime] = useState<any>(null);
  const [telegram, setTelegram] = useState<any>(null);
  const [envVars, setEnvVars] = useState<any[]>([]);
  const [envKey, setEnvKey] = useState("");
  const [envValue, setEnvValue] = useState("");
  const [botToken, setBotToken] = useState("");
  const [botMode, setBotMode] = useState("polling");

  const load = async (projectId: string) => {
    try {
      const [runtimeData, envData] = await Promise.all([
        api.runtimeStatus(projectId),
        api.envList(projectId),
      ]);
      setRuntime(runtimeData);
      setEnvVars(envData);
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
        <div className={styles.row}>
          <input className={styles.input} placeholder="KEY" value={envKey} onChange={(e) => setEnvKey(e.target.value)} />
          <input className={styles.input} placeholder="VALUE" value={envValue} onChange={(e) => setEnvValue(e.target.value)} />
          <button className={styles.btn} onClick={saveEnv}>Save</button>
        </div>
        <div className={styles.envList}>
          {envVars.map((item) => (
            <div key={item.id} className={styles.envItem}>
              {item.key} = {item.value}
            </div>
          ))}
        </div>
      </div>

      <div className={styles.card}>
        <h3>Telegram Bot</h3>
        <div>Current mode: {telegram?.mode || "not configured"}</div>
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
    </div>
  );
}
