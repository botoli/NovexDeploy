import { makeAutoObservable } from "mobx";

interface User {
  name: string;
  github_login: string;
  email: string;
  avatar_url: string;
}

export const GithubUser = {
  user: null as User | null,
  loading: false,
  error: "",

  setUser(data: User) {
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
