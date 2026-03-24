import { Route, Routes } from "react-router";
import { DashboardPage } from "./Pages/dashboard";
import { ProjectsPage } from "./Pages/projects/ProjectsPage";
import { LeftPanel } from "./Pages/LeftPanel";
import "./GlobalStyles/App.scss";
import { useEffect } from "react";
import { GithubUser } from "./Store/GithubInfo.store";
import { observer } from "mobx-react-lite";
import DeployPage from "./Pages/deploy";
import { GithubRepo } from "./Store/Repos.store";

const App = observer(() => {
  const fetchRepos = async () => {
    GithubRepo.setLoading(true);
    try {
      const response = await fetch("/git/repos", {
        credentials: "include",
        headers: {
          Accept: "application/json", // Явно указываем, что ждем JSON
        },
      });

      // Проверяем статус ответа
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      // Проверяем тип контента
      const contentType = response.headers.get("content-type");
      if (!contentType || !contentType.includes("application/json")) {
        throw new Error("Server didn't return JSON");
      }

      const data = await response.json();
      console.log("Full Server Response:", data); // Debug logging

      if (data.ok) {
        if (data.data) {
          console.log(
            "Array info:",
            Array.isArray(data.data),
            "Length:",
            data.data.length,
          );
          GithubRepo.setRepos(data.data);
        } else {
          console.warn("Data field is missing or null", data);
          GithubRepo.setRepos([]);
        }
      } else {
        GithubRepo.setError(data.message || "Unknown error");
      }
    } catch (err) {
      console.error("Fetch repos error:", err);
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
      const response = await fetch("/auth/me", {
        credentials: "include",
        headers: {
          Accept: "application/json",
        },
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const contentType = response.headers.get("content-type");
      if (!contentType || !contentType.includes("application/json")) {
        throw new Error("Server didn't return JSON");
      }

      const data = await response.json();

      if (data.ok) {
        GithubUser.setUser(data.data);
      } else {
        GithubUser.setError(data.message || "Unknown error");
      }
    } catch (err) {
      console.error("Fetch error:", err);
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
      <Routes>
        <Route path="/dashboard" element={<DashboardPage />} />
        <Route path="/projects" element={<ProjectsPage />} />
        <Route path="/deploy" element={<DeployPage />} />
      </Routes>
    </div>
  );
});

export default App;
