import { Route, Routes } from "react-router";
import { DashboardPage } from "./Pages/dashboard";
import { ProjectsPage } from "./Pages/projects/ProjectsPage";
import { LeftPanel } from "./Pages/LeftPanel";
import "./GlobalStyles/App.scss";
import { useEffect } from "react";
import { GithubUser } from "./Store/GithubInfo.store";
import { observer } from "mobx-react-lite";
const App = observer(() => {
  useEffect(() => {
    const fetchGitHubProfile = async () => {
      GithubUser.setLoading(true);
      try {
        // Используем относительный путь - прокси перенаправит
        const response = await fetch("/auth/me", {
          credentials: "include",
        });

        const data = await response.json();

        if (data.ok) {
          GithubUser.setUser(data.data);
        } else {
          GithubUser.setError(data.message);
        }
      } catch (err) {
        console.error("Fetch error:", err);
        GithubUser.setError("Failed to fetch profile");
      } finally {
        GithubUser.setLoading(false);
      }
    };
    fetchGitHubProfile();
  }, []);
  return (
    <div className="AllApp">
      <LeftPanel />
      <Routes>
        <Route path="/dashboard" element={<DashboardPage />} />
        <Route path="/projects" element={<ProjectsPage />} />
      </Routes>
    </div>
  );
});
export default App;
