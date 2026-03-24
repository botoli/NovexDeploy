import { useState, useMemo } from "react";
import { observer } from "mobx-react-lite";
import { ArrowBack } from "@mui/icons-material";
import {
  Link as LinkIcon,
  Search,
  FilePlus,
  Lock,
  Globe
} from "lucide-react";
import styles from "./index.module.scss";
import { useNavigate } from "react-router-dom";
import { GithubRepo } from "../../Store/Repos.store";
import { api } from "../../shared/api";

const DeployPage = observer(() => {
  const navigate = useNavigate();
  const [importingRepoId, setImportingRepoId] = useState<number | null>(null);
  const [searchQuery, setSearchQuery] = useState("");
  const [rootDir, setRootDir] = useState(".");
  const [buildCmd, setBuildCmd] = useState("npm ci");
  const [outputDir, setOutputDir] = useState(".");
  const [startCmd, setStartCmd] = useState("npm start");

  const filteredRepos = useMemo(() => {
    if (!GithubRepo.repos) return [];
    return GithubRepo.repos.filter((repo) =>
      repo.name.toLowerCase().includes(searchQuery.toLowerCase()),
    );
  }, [GithubRepo.repos, searchQuery]);

  const handleImport = async (repo: any) => {
    setImportingRepoId(repo.id);
    try {
      // 1. Create Project
      const created = await api.createProject({
        name: repo.name,
        description: repo.description || "",
        framework: "node",
        project_type: "backend",
        build_command: buildCmd,
        root_dir: rootDir,
        output_dir: outputDir,
        start_command: startCmd,
      });

      const projectId = created.id;

      // 2. Connect Repo & Trigger Build
      await api.connectRepo(projectId, {
        repo_full_name: repo.full_name,
        branch: repo.default_branch,
        build_command: buildCmd,
        root_dir: rootDir,
        output_dir: outputDir,
      });

      navigate("/deployments");
    } catch (err: any) {
      alert("Error importing repository: " + err.message);
    } finally {
      setImportingRepoId(null);
    }
  };

  return (
    <div className={styles.page}>
      <header className={styles.header}>
        <a href="/dashboard" className={styles.backLink}>
          <ArrowBack fontSize="small" />
          <span>Back</span>
        </a>

        {/* Center title or logo could go here */}

        {/* Right side user or nothing, assuming layout handles specific user icon */}
      </header>

      <h1>Deploy backend service or Telegram bot</h1>

      <div className={styles.inputSection}>
        <div className={styles.inputWrapper}>
          <LinkIcon size={16} className={styles.linkIcon} />
          <input
            type="text"
            placeholder="Select a repository below (GitHub OAuth required)"
          />
        </div>
        <div style={{ display: "grid", gridTemplateColumns: "180px 1fr 180px 1fr", gap: 8, marginTop: 10 }}>
          <input value={rootDir} onChange={(e) => setRootDir(e.target.value)} placeholder="Root dir (.)" />
          <input value={buildCmd} onChange={(e) => setBuildCmd(e.target.value)} placeholder="Build command" />
          <input value={outputDir} onChange={(e) => setOutputDir(e.target.value)} placeholder="Output dir (use . for backend runtime root)" />
          <input value={startCmd} onChange={(e) => setStartCmd(e.target.value)} placeholder="Start command" />
        </div>
      </div>

      <div className={styles.mainGrid}>
        {/* Left Column: Import Git Repo */}
        <section className={styles.importSection}>
          <div className={styles.sectionTitle}>Import Git Repository</div>

          <div className={styles.repoSelector}>
            <div className={styles.search}>
              <Search size={16} />
              <input
                placeholder="Search..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
            </div>
          </div>

          <div className={styles.repoList}>
            {GithubRepo.loading && (
              <div
                style={{
                  padding: 20,
                  textAlign: "center",
                  color: "#666",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                  gap: 10,
                }}
              >
                <div className={styles.spinner} /> Loading...
              </div>
            )}
            {GithubRepo.error && (
              <div style={{ color: "red", padding: 10 }}>
                Error: {GithubRepo.error}
              </div>
            )}

            {!GithubRepo.loading &&
              !GithubRepo.error &&
              filteredRepos.length === 0 && (
                <div
                  style={{ padding: 20, textAlign: "center", color: "#666" }}
                >
                  {searchQuery
                    ? "No repositories match query."
                    : "No repositories found."}
                </div>
              )}

            {filteredRepos.map((repo) => (
              <div key={repo.id} className={styles.repoItem}>
                <div className={styles.repoInfo}>
                  <div className={styles.repoIcon}>
                    {repo.private ? <Lock size={16} /> : <Globe size={16} />}
                  </div>
                  <div className={styles.repoDetail}>
                    <div
                      style={{
                        display: "flex",
                        alignItems: "center",
                        gap: "8px",
                      }}
                    >
                      <span className={styles.name}>{repo.name}</span>
                      {repo.private && (
                        <span className={styles.privateBadge}>Private</span>
                      )}
                    </div>
                    <span className={styles.time}>
                      {repo.language && `${repo.language} • `} Updated{" "}
                      {new Date(repo.updated_at).toLocaleDateString()}
                    </span>
                  </div>
                </div>
                <button
                  className={styles.importBtn}
                  onClick={() => handleImport(repo)}
                  disabled={importingRepoId === repo.id}
                >
                  {importingRepoId === repo.id ? "Importing..." : "Import"}
                </button>
              </div>
            ))}
          </div>
        </section>

        {/* Right Column: Clone Template */}
      </div>

      <footer className={styles.footer}>
        <div className={styles.footerBtnWrapper}>
          <div className={styles.footerText}>
            <FilePlus size={20} />
            <div>
              <h3>Create Empty Project</h3>
              <p>
                Create a backend/telegram project without repository import.
              </p>
            </div>
          </div>

          <button className={styles.emptyProjectBtn}>
            Create Empty Project
          </button>
        </div>
      </footer>
    </div>
  );
});

export default DeployPage;
