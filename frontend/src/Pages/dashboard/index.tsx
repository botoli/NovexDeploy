import { useEffect, useState } from "react";
import { NotificationBell, UserCircleIcon, UserIcon } from "../../Icons/Icons";
import styles from "./index.module.scss";
import { Cpu, Disc, MemoryStick, MemoryStickIcon, Rabbit } from "lucide-react";
import { MemoryRouter } from "react-router-dom";
import { Button } from "@mui/material";

export function DashboardPage() {
  interface Deployment {
    id: number;
    projectId: number;
    projectName: string;
    status: "pending" | "building" | "success" | "failed";
    commitHash: string;
    commitMessage: string;
    branch: string;
    url: string;
    createdAt: string;
    finishedAt?: string;
    duration?: number; // в секундах
    logs?: string[];
  }
  const [MockDeplouments, setMockDeplouments] = useState<Deployment[]>([
    {
      id: 1,
      projectId: 1,
      projectName: "novex-task",
      status: "success",
      commitHash: "a1b2c3d",
      commitMessage: "feat: add dark mode support",
      branch: "main",
      url: "https://novex-task.vercel.app",
      createdAt: "2024-03-15T10:30:00Z",
      finishedAt: "2024-03-15T10:32:15Z",
      duration: 135,
      logs: [
        "🔨 Cloning repository...",
        "📦 Installing dependencies...",
        "🏗️ Building",
        "✅ Deployment successful!",
      ],
    },
    {
      id: 2,
      projectId: 1,
      projectName: "novex-task",
      status: "building",
      commitHash: "d4e5f6g",
      commitMessage: "fix: header alignment",
      branch: "feature/ui-fixes",
      url: "",
      createdAt: "2024-03-18T14:15:00Z",
      logs: [
        "🔨 Cloning repository...",
        "📦 Installing dependencies...",
        "🏗️ Building...",
      ],
    },
    {
      id: 3,
      projectId: 1,
      projectName: "novex-task",
      status: "failed",
      commitHash: "h7i8j9k",
      commitMessage: "chore: update dependencies",
      branch: "main",
      url: "",
      createdAt: "2024-03-14T09:45:00Z",
      finishedAt: "2024-03-14T09:47:00Z",
      duration: 120,
      logs: [
        "🔨 Cloning repository...",
        "📦 Installing dependencies...",
        "❌ Build failed",
      ],
    },
    {
      id: 4,
      projectId: 3,
      projectName: "go-backend",
      status: "success",
      commitHash: "l0m1n2o",
      commitMessage: "feat: add user auth",
      branch: "main",
      url: "https://go-backend.vercel.app",
      createdAt: "2024-03-10T12:00:00Z",
      finishedAt: "2024-03-10T12:01:30Z",
      duration: 90,
      logs: [
        "🔨 Cloning repository...",
        "📦 Installing Go modules...",
        "🏗️ Building binary...",
        "✅ Deployment successful!",
      ],
    },
    {
      id: 5,
      projectId: 2,
      projectName: "novex-deploy",
      status: "pending",
      commitHash: "p3q4r5s",
      commitMessage: "docs: update README",
      branch: "main",
      url: "",
      createdAt: "2024-03-18T16:50:00Z",
      logs: ["⏳ Queued for deployment..."],
    },
  ]);
  const [metrics, setMetrics] = useState<any>(null);
  const fetchMetrics = () => {
    fetch("/metrics/system")
      .then((response) => response.json())
      .then((data) => {
        if (data && data.data) {
          setMetrics(data.data);
        }
      })
      .catch((error) => {
        console.error("Error fetching metrics:", error);
      });
  };

  useEffect(() => {
    fetchMetrics();
    const interval = setInterval(fetchMetrics, 10000);
    return () => clearInterval(interval);
  }, []);
  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <div className={styles.burgermenu}>
          <div className={styles.line}></div>
          <div className={styles.line}></div>
          <div className={styles.line}></div>
        </div>
        <h1 className={styles.title}>Novex Deploy</h1>
        <div className={styles.navigation}>
          <button className={styles.DeployBtn} onClick={() => {}}>
            <p>Deploy</p>
          </button>
          <button className={styles.notificationBtn}>
            <NotificationBell />
          </button>
          <button className={styles.userBtn}>
            <UserIcon />
          </button>
        </div>
      </header>
      <div className={styles.content}>
        <div className={styles.DeploySection}>
          <div className={styles.Active}>
            <h2>Active Deployments</h2>
            <div className={styles.DeploymentsList}>
              {MockDeplouments.map((d) => (
                <div className={styles.RowsDiv}>
                  <div className={styles.Rows}>{d.projectName}</div>
                  <div className={styles.Rows}>
                    {d.logs?.findLast((log) => log !== undefined)}
                  </div>
                </div>
              ))}
            </div>
          </div>
          <div className={styles.Last}>
            <h2>Last Deployments</h2>
            <div className={styles.DeploymentsList}>
              {MockDeplouments.map((d) => (
                <div className={styles.RowsDiv}>
                  <div className={styles.Rows}>{d.projectName}</div>
                  <div
                    className={
                      styles.RowsStatus +
                      (d.status === "success"
                        ? ` ${styles.RowsStatussuccess}`
                        : d.status === "failed"
                          ? ` ${styles.RowsStatuserror}`
                          : d.status === "building"
                            ? ` ${styles.RowsStatusbuilding}`
                            : "")
                    }
                  >
                    {d.status}
                  </div>
                </div>
              ))}
            </div>
          </div>
          <div className={styles.Usage}>
            <h2>Usage Statistics</h2>
            <div className={styles.DeploymentsList}>
              {metrics ? (
                <>
                  <div className={styles.RowsDiv}>
                    <div className={styles.Rows}>
                      <Cpu className={styles.iconCpu} />
                      CPU Usage
                    </div>
                    <div className={styles.Rows}>
                      <div
                        style={{
                          width: "120px",
                          height: "8px",
                          backgroundColor: "rgba(255,255,255,0.1)",
                          borderRadius: "4px",
                          overflow: "hidden",
                        }}
                      >
                        <div
                          style={{
                            width: `${Math.min(100, Math.max(0, metrics.cpu_percent))}%`,
                            height: "100%",
                            backgroundColor:
                              metrics.cpu_percent > 80 ? "#fb7185" : "#38bdf8",
                            transition: "width 0.5s ease",
                          }}
                        />
                      </div>
                      <div style={{ minWidth: "3rem", textAlign: "right" }}>
                        {metrics.cpu_percent?.toFixed(1)}%
                      </div>
                    </div>
                  </div>
                  <div className={styles.RowsDiv}>
                    <div className={styles.Rows}>
                      <MemoryStick className={styles.iconMemory} />
                      RAM Usage
                    </div>
                    <div className={styles.Rows}>
                      <div
                        style={{
                          width: "120px",
                          height: "8px",
                          backgroundColor: "rgba(255,255,255,0.1)",
                          borderRadius: "4px",
                          overflow: "hidden",
                        }}
                      >
                        <div
                          style={{
                            width: `${Math.min(
                              100,
                              (metrics.ram_mb / metrics.ram_total_mb) * 100,
                            )}%`,
                            height: "100%",
                            backgroundColor: "#818cf8",
                            transition: "width 0.5s ease",
                          }}
                        />
                      </div>
                      <div style={{ minWidth: "7rem", textAlign: "right" }}>
                        {(metrics.ram_mb / 1024).toFixed(1)} /{" "}
                        {(metrics.ram_total_mb / 1024).toFixed(1)} GB
                      </div>
                    </div>
                  </div>
                  <div className={styles.RowsDiv}>
                    <div className={styles.Rows}>
                      <Disc className={styles.iconDisk} />
                      Disk Usage
                    </div>
                    <div className={styles.Rows}>
                      <div
                        style={{
                          width: "120px",
                          height: "8px",
                          backgroundColor: "rgba(255,255,255,0.1)",
                          borderRadius: "4px",
                          overflow: "hidden",
                        }}
                      >
                        <div
                          style={{
                            width: `${Math.min(
                              100,
                              (metrics.disk_mb / metrics.disk_total_mb) * 100,
                            )}%`,
                            height: "100%",
                            backgroundColor: "#34d399",
                            transition: "width 0.5s ease",
                          }}
                        />
                      </div>
                      <div style={{ minWidth: "7rem", textAlign: "right" }}>
                        {(metrics.disk_mb / 1024).toFixed(1)} /{" "}
                        {(metrics.disk_total_mb / 1024).toFixed(1)} GB
                      </div>
                    </div>
                  </div>
                </>
              ) : (
                <div className={styles.RowsDiv}>
                  <div className={styles.Rows}>Loading metrics...</div>
                </div>
              )}
              <Button color="primary" onClick={() => fetchMetrics()}>
                Secondary
              </Button>
            </div>
          </div>
          <div className={styles.QuickActions}>
            <h2>Quick Actions</h2>
            <div className={styles.QuickActionsButtons}>
              <button className={styles.QuickActionBtn}>
                <p>New Frontend</p>
              </button>
              <button className={styles.QuickActionBtn}>
                <p>New Backend</p>
              </button>
              <button className={styles.QuickActionBtn}>
                <p>New Bot</p>
              </button>
            </div>
          </div>
        </div>
        <div className={styles.DeployedAppsSection}>
          <h2>Deployed Applications</h2>
          <div className={styles.AppsList}>
            {/* List of deployed applications */}
          </div>
        </div>
        <div className={styles.HistorySection}>
          <h2>Deployment History</h2>
          <div className={styles.HistoryList}>
            {/* List of deployment history */}
          </div>
        </div>
      </div>
    </div>
  );
}
