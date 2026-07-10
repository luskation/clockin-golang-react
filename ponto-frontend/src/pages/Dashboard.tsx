import { useEffect, useState } from "react";
import axios from "axios";
import styles from "./Dashboard.module.css";

type Status = "checking" | "online" | "offline";

const apiRoot = import.meta.env.VITE_API_URL.replace(/\/api\/v1\/?$/, "");

export default function Dashboard() {
  const [status, setStatus] = useState<Status>("checking");

  useEffect(() => {
    axios
      .get(`${apiRoot}/health`)
      .then(() => setStatus("online"))
      .catch(() => setStatus("offline"));
  }, []);

  return (
    <div className={styles.dashboard}>
      <h1 className={styles.title}>Dashboard</h1>
      <p className={styles.subtitle}>
        A tela de bater ponto chega na semana 5. Por enquanto, aqui está o
        status da API:
      </p>
      <p className={styles.status} data-status={status}>
        {status}
      </p>
    </div>
  );
}
