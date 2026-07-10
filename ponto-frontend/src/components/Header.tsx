import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { getCurrentUser, logout } from "../services/auth";
import styles from "./Header.module.css";

export default function Header() {
  const navigate = useNavigate();
  const user = getCurrentUser();
  const [menuOpen, setMenuOpen] = useState(false);

  function handleLogout() {
    logout();
    navigate("/login", { replace: true });
  }

  return (
    <header className={styles.header}>
      <Link to="/" className={styles.brand} onClick={() => setMenuOpen(false)}>
        Ponto
      </Link>

      <button
        className={styles.menuToggle}
        type="button"
        aria-expanded={menuOpen}
        aria-controls="main-nav"
        onClick={() => setMenuOpen((open) => !open)}
      >
        <span className={styles.menuIcon} data-open={menuOpen} />
        <span className={styles.srOnly}>Menu</span>
      </button>

      <nav
        id="main-nav"
        className={styles.nav}
        data-open={menuOpen}
      >
        {user?.role === "admin" && (
          <Link to="/admin" className={styles.navLink} onClick={() => setMenuOpen(false)}>
            Empresas e usuários
          </Link>
        )}

        <div className={styles.userInfo}>
          <span className={styles.userName}>{user?.name}</span>
          <span className={styles.userRole}>{user?.role}</span>
        </div>

        <button className={styles.logout} type="button" onClick={handleLogout}>
          Sair
        </button>
      </nav>
    </header>
  );
}
