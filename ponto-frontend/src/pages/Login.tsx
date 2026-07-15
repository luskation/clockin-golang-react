import { useState, type FormEvent } from "react";
import { useNavigate } from "react-router-dom";
import { login } from "../services/auth";
import { useLiveClock } from "../hooks/useLiveClock";
import { formatDate, formatTime } from "../utils/datetime";
import mark from "../assets/brand/mark-white.svg";
import styles from "./Login.module.css";

export default function Login() {
  const now = useLiveClock();
  const navigate = useNavigate();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      await login(email, password);
      navigate("/", { replace: true });
    } catch {
      setError("Email ou senha incorretos.");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className={styles.screen}>
      <aside className={styles.clockPanel}>
        <span className={styles.date}>{formatDate(now)}</span>
        <span className={styles.time}>{formatTime(now)}</span>
        <div className={styles.divider} />
        <img className={styles.brandMark} src={mark} alt="" aria-hidden="true" />
        <h1 className={styles.brand}>Cronos</h1>
        <p className={styles.tagline}>
          Entrada e saída registradas na hora certa, para o time da Comp Júnior.
        </p>
      </aside>

      <main className={styles.formPanel}>
        <form className={styles.form} onSubmit={handleSubmit} noValidate>
          <h2 className={styles.formTitle}>Entrar</h2>

          <label className={styles.label} htmlFor="email">
            Email
          </label>
          <input
            id="email"
            className={styles.input}
            type="email"
            autoComplete="username"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />

          <label className={styles.label} htmlFor="password">
            Senha
          </label>
          <input
            id="password"
            className={styles.input}
            type="password"
            autoComplete="current-password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />

          {error && (
            <p className={styles.error} role="alert">
              {error}
            </p>
          )}

          <button className={styles.submit} type="submit" disabled={loading}>
            {loading ? "Entrando…" : "Entrar"}
          </button>
        </form>
      </main>
    </div>
  );
}
