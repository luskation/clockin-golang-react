import { Link } from "react-router-dom";
import styles from "./AdminHome.module.css";

export default function AdminHome() {
  return (
    <div className={styles.page}>
      <span className={styles.eyebrow}>Administração</span>
      <h1 className={styles.title}>Empresas e usuários</h1>
      <p className={styles.subtitle}>
        Gerencie os cadastros que sustentam o registro de ponto da sua
        organização.
      </p>

      <div className={styles.grid}>
        <Link to="/admin/companies" className={styles.card}>
          <span className={styles.cardLabel}>Empresas</span>
          <span className={styles.cardHint}>
            Cadastre e edite as empresas vinculadas à conta.
          </span>
        </Link>

        <Link to="/admin/users" className={styles.card}>
          <span className={styles.cardLabel}>Usuários</span>
          <span className={styles.cardHint}>
            Gerencie contas, papéis e vínculos com empresas.
          </span>
        </Link>
      </div>
    </div>
  );
}
