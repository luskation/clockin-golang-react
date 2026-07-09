import { useEffect, useState } from "react";
import axios from "axios";

type Status = "checking" | "online" | "offline";

const apiRoot = import.meta.env.VITE_API_URL.replace(/\/api\/v1\/?$/, "");

function Health() {
  const [status, setStatus] = useState<Status>("checking");

  useEffect(() => {
    axios
      .get(`${apiRoot}/health`)
      .then(() => setStatus("online"))
      .catch(() => setStatus("offline"));
  }, []);

  return (
    <div className="health-check">
      <h1>Ponto</h1>
      <p>
        API: <strong>{status}</strong>
      </p>
    </div>
  );
}

export default Health;
