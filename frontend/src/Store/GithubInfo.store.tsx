import { makeAutoObservable } from "mobx";

type GitHubUserType = {
  id: string;
  email: string;
  name: string;
  avatar_url?: string;
  github_login?: string;
  github_id?: number;
};

export const GithubUser = {
  user: null as GitHubUserType | null,
  loading: false,
  error: "",

  setUser(data: GitHubUserType) {
    this.user = data;
    this.error = "";
  },

  setError(error: string) {
    this.error = error;
  },

  setLoading(loading: boolean) {
    this.loading = loading;
  },
};

makeAutoObservable(GithubUser);
