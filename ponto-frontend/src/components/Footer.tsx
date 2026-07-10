import styles from "./Footer.module.css";

export default function Footer() {
  return (
    <footer className={styles.footer}>
      <span>Ponto</span>
      <span>&copy; {new Date().getFullYear()}</span>
    </footer>
  );
}
