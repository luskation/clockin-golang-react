import styles from "./Footer.module.css";

export default function Footer() {
  return (
    <footer className={styles.footer}>
      <span>Cronos · Comp Júnior</span>
      <span>&copy; {new Date().getFullYear()}</span>
    </footer>
  );
}
