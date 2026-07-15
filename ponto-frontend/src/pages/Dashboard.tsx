import { useEffect, useState } from "react";
import api from "../services/api";
import Toast from "../components/Toast";
import Pagination from "../components/Pagination";
import { useLiveClock } from "../hooks/useLiveClock";
import { formatDate, formatTime } from "../utils/datetime";
import { getCurrentUser } from "../services/auth";
import styles from "./Dashboard.module.css";

type EntryType = "clock_in" | "clock_out";

interface TimeEntry {
  id: string;
  type: EntryType;
  recorded_at: string;
}

const LIMIT = 5;

const typeLabel: Record<EntryType, string> = { clock_in: "Entrada", clock_out: "Saída" };
const actionLabel: Record<EntryType, string> = { clock_in: "Registrar entrada", clock_out: "Registrar saída" };

export default function Dashboard() {
  const now = useLiveClock();
  const user = getCurrentUser();

  const [entries, setEntries] = useState<TimeEntry[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [loadingEntries, setLoadingEntries] = useState(true);
  const [clocking, setClocking] = useState(false);
  const [toast, setToast] = useState<{ message: string; variant: "success" | "error" } | null>(null);

  function loadEntries() {
    setLoadingEntries(true);
    api
      .get("/time-entries/me", { params: { page, limit: LIMIT } })
      .then((res) => {
        setEntries(res.data.data ?? []);
        setTotal(res.data.total ?? 0);
      })
      .catch(() => setToast({ message: "Não foi possível carregar seu histórico.", variant: "error" }))
      .finally(() => setLoadingEntries(false));
  }

  useEffect(loadEntries, [page]);

  // Só um preview do rótulo do botão — quem decide o tipo de verdade é o
  // backend (RegisterEntry), a partir do último registro real no banco.
  const lastEntry = page === 1 ? entries[0] : undefined;
  const nextType: EntryType = lastEntry?.type === "clock_in" ? "clock_out" : "clock_in";

  async function handleClock() {
    setClocking(true);
    try {
      const { data } = await api.post("/time-entries");
      const registeredType = data.type as EntryType;
      setToast({
        message: `${typeLabel[registeredType]} registrada às ${formatTime(new Date(data.recorded_at))}.`,
        variant: "success",
      });
      if (page === 1) {
        loadEntries();
      } else {
        setPage(1);
      }
    } catch {
      setToast({ message: "Erro ao registrar ponto. Tente novamente.", variant: "error" });
    } finally {
      setClocking(false);
    }
  }

  return (
    <div className={styles.dashboard}>
      <span className={styles.eyebrow}>{formatDate(now)}</span>
      <h1 className={styles.title}>Olá, {user?.name?.split(" ")[0] ?? "colaborador"}</h1>
      <p className={styles.clock}>{formatTime(now)}</p>

      <button className={styles.clockButton} type="button" onClick={handleClock} disabled={clocking}>
        {clocking ? "Registrando…" : actionLabel[nextType]}
      </button>

      <section className={styles.history}>
        <h2 className={styles.historyTitle}>Últimos registros</h2>

        {loadingEntries ? (
          <p className={styles.hint}>Carregando…</p>
        ) : entries.length === 0 ? (
          <p className={styles.hint}>Nenhum registro ainda. Bata seu primeiro ponto acima.</p>
        ) : (
          <ul className={styles.list}>
            {entries.map((entry) => (
              <li key={entry.id} className={styles.listItem}>
                <span className={styles.entryType} data-type={entry.type}>
                  {typeLabel[entry.type]}
                </span>
                <span className={styles.entryTime}>{new Date(entry.recorded_at).toLocaleString("pt-BR")}</span>
              </li>
            ))}
          </ul>
        )}

        <Pagination page={page} limit={LIMIT} total={total} onPageChange={setPage} />
      </section>

      {toast && <Toast message={toast.message} variant={toast.variant} onDismiss={() => setToast(null)} />}
    </div>
  );
}
