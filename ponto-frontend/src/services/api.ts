import axios from "axios";

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  if (token) config.headers.Authorization = `Bearer ${token}`;
  return config;
});

api.interceptors.response.use(
  (res) => res,
  (err) => {
    // Um 401 em /auth/* é login/senha errados, tratado pela própria tela —
    // não deve disparar o logout global, senão a mensagem de erro do
    // formulário nunca chegaria a aparecer.
    const isAuthEndpoint = err.config?.url?.startsWith("/auth/");
    if (err.response?.status === 401 && !isAuthEndpoint) {
      localStorage.clear();
      window.location.href = "/login";
    }
    return Promise.reject(err);
  },
);

export default api;
