import { makeAutoObservable } from "mobx";

export const GithubUser = {
  user: null,
  loading: false,
  error: "",

  setUser(data: any) {
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
