import { Route, Routes } from "react-router";
import { DashboardPage } from "./Pages/dashboard";
import { ProjectsPage } from "./Pages/projects/ProjectsPage";
import { LeftPanel } from "./Pages/LeftPanel";
import "./GlobalStyles/App.scss";
function App() {
  return (
    <div className="AllApp">
      <LeftPanel />
      <Routes>
        <Route path="/dashboard" element={<DashboardPage />} />
        <Route path="/projects" element={<ProjectsPage />} />
      </Routes>
    </div>
  );
}

export default App;
