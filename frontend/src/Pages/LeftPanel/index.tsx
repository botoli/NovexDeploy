import React from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { Home, Box, Rocket, Terminal, Settings, Bot } from "lucide-react";
import styles from "./index.module.scss";

const NAV_ITEMS = [
  {
    label: "Dashboard",
    icon: Home,
    path: "/dashboard",
  },
  {
    label: "Projects",
    icon: Box,
    path: "/projects",
  },
  {
    label: "Deployments",
    icon: Rocket,
    path: "/deployments",
  },
  {
    label: "Logs",
    icon: Terminal,
    path: "/logs",
  },
  // если хочешь включить ботов — раскомментируй
  // {
  //   label: "Bots",
  //   icon: Bot,
  //   path: "/bots",
  // },
  {
    label: "Settings",
    icon: Settings,
    path: "/settings",
  },
];

export function LeftPanel() {
  const location = useLocation();
  const navigate = useNavigate();

  const isActive = (path: string) => location.pathname.startsWith(path);

  return (
    <div className={styles.container}>
      {/* Logo */}
      <div className={styles.logo}>⚡</div>

      {/* Navigation */}
      <div className={styles.nav}>
        {NAV_ITEMS.map((item) => {
          const Icon = item.icon;
          const active = isActive(item.path);

          return (
            <div
              key={item.path}
              className={`${styles.item} ${active ? styles.active : ""}`}
              onClick={() => navigate(item.path)}
            >
              <Icon size={20} />
              <span className={styles.label}>{item.label}</span>
            </div>
          );
        })}
      </div>

      {/* Bottom (user) */}
      <div className={styles.bottom}>
        <div className={styles.user}>
          <div className={styles.avatar}>A</div>
        </div>
      </div>
    </div>
  );
}
