import { Navigate, Route, Routes } from "react-router-dom";
import { DashboardPage } from "./Pages/dashboard";
import { ProjectsPage } from "./Pages/projects/ProjectsPage";
import "./GlobalStyles/App.scss";
import { useEffect } from "react";
import { GithubUser } from "./Store/GithubInfo.store";
import { observer } from "mobx-react-lite";
import DeployPage from "./Pages/deploy";
import { GithubRepo } from "./Store/Repos.store";
import { LeftPanel } from "./Pages/LeftPanel";
import { DeploymentsPage } from "./Pages/Deployments";
import { api } from "./shared/api";
import { SettingsPage } from "./Pages/Overview";

const App = observer(() => {
  const fetchRepos = async () => {
    GithubRepo.setLoading(true);
    try {
      const repos = await api.repos();
      GithubRepo.setRepos(repos);
    } catch (err) {
      GithubRepo.setError(
        err instanceof Error ? err.message : "Failed to fetch repos",
      );
    } finally {
      GithubRepo.setLoading(false);
    }
  };

  const fetchGitHubProfile = async () => {
    GithubUser.setLoading(true);
    try {
      const user = await api.me();
      GithubUser.setUser(user);
    } catch (err) {
      GithubUser.setError(
        err instanceof Error ? err.message : "Failed to fetch profile",
      );
    } finally {
      GithubUser.setLoading(false);
    }
  };

  useEffect(() => {
    const init = async () => {
      await fetchGitHubProfile();
      await fetchRepos();
    };

    init();
  }, []);

  return (
    <div className="AllApp">
      <LeftPanel />
      <div className="MainContent">
        <Routes>
          <Route path="/" element={<Navigate to="/dashboard" replace />} />
          <Route path="/dashboard" element={<DashboardPage />} />
          <Route path="/projects" element={<ProjectsPage />} />
          <Route path="/deploy" element={<DeployPage />} />
          <Route path="/deployments" element={<DeploymentsPage />} />
          <Route path="/logs" element={<DeploymentsPage logsOnly />} />
          <Route path="/settings" element={<SettingsPage />} />
          <Route path="*" element={<Navigate to="/dashboard" replace />} />
        </Routes>
      </div>
    </div>
  );
});

export default App;
