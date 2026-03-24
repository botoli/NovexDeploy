import { makeAutoObservable } from "mobx";

export interface GitHubRepo {
  id: number;
  name: string;
  full_name: string;
  description: string;
  private: boolean;
  html_url: string;
  clone_url: string;
  language: string;
  updated_at: string;
  default_branch: string;
}

class GithubRepoStore {
  repos: GitHubRepo[] = [];
  loading = false;
  error = "";

  constructor() {
    makeAutoObservable(this);
  }

  setRepos(data: GitHubRepo[]) {
    if (Array.isArray(data)) {
      this.repos = data;
    } else {
      console.warn("Received non-array repos data:", data);
      this.repos = [];
    }
    this.error = "";
  }

  setError(error: string) {
    this.error = error;
  }

  setLoading(loading: boolean) {
    this.loading = loading;
  }
}

export const GithubRepo = new GithubRepoStore();
