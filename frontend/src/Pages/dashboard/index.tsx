import { useState } from "react";
import { NotificationBell, UserCircleIcon, UserIcon } from "../../Icons/Icons";
import styles from "./index.module.scss";

export function DashboardPage() {
  interface Depl {
    id: number;
    name: string;
    status: string;
    category: string;
  }
  const [MockDeplouments, setMockDeplouments] = useState<Depl[]>([
    { id: 1, name: "NovexTask", status: "building", category: "Frontend" },
    { id: 1, name: "NovexDeploy", status: "Active", category: "Fullstack" },
  ]);
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
                  <div className={styles.Rows}>{d.name}</div>
                  <div className={styles.Rows}>{d.status}</div>
                </div>
              ))}
            </div>
          </div>
          <div className={styles.Last}>
            <h2>Last Deployments</h2>
            <div className={styles.DeploymentsList}>
              {/* List of last deployments */}
            </div>
          </div>
          <div className={styles.Usage}>
            <h2>Usage Statistics</h2>
            <div className={styles.DeploymentsList}>
              {/* List of usage statistics */}
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
