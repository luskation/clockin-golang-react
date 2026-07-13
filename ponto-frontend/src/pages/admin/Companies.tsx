import { useEffect, useState, type FormEvent } from "react";
import api from "../../services/api";
import Modal from "../../components/Modal";
import Toast from "../../components/Toast";
import Pagination from "../../components/Pagination";
import { formatCNPJ, isValidCNPJ } from "../../utils/validation";
import styles from "./Companies.module.css";

interface Company {
  id: string;
  name: string;
  cnpj: string;
  created_at: string;
}

interface FormState {
  name: string;
  cnpj: string;
}

const emptyForm: FormState = { name: "", cnpj: "" };
const LIMIT = 10;

export default function Companies() {
  const [companies, setCompanies] = useState<Company[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(true);
  const [toast, setToast] = useState<{ message: string; variant: "success" | "error" } | null>(null);

  const [modalOpen, setModalOpen] = useState(false);
  const [editing, setEditing] = useState<Company | null>(null);
  const [form, setForm] = useState<FormState>(emptyForm);
  const [fieldErrors, setFieldErrors] = useState<Partial<FormState>>({});
  const [formError, setFormError] = useState("");
  const [saving, setSaving] = useState(false);

  const [confirmingId, setConfirmingId] = useState<string | null>(null);
  const [deletingId, setDeletingId] = useState<string | null>(null);

  function loadCompanies() {
    setLoading(true);
    api
      .get("/companies", { params: { page, limit: LIMIT } })
      .then((res) => {
        setCompanies(res.data.data ?? []);
        setTotal(res.data.total ?? 0);
      })
      .catch(() => setToast({ message: "Não foi possível carregar as empresas.", variant: "error" }))
      .finally(() => setLoading(false));
  }

  useEffect(loadCompanies, [page]);

  function openCreate() {
    setEditing(null);
    setForm(emptyForm);
    setFieldErrors({});
    setFormError("");
    setModalOpen(true);
  }

  function openEdit(company: Company) {
    setEditing(company);
    setForm({ name: company.name, cnpj: company.cnpj });
    setFieldErrors({});
    setFormError("");
    setModalOpen(true);
  }

  function validate(): boolean {
    const errors: Partial<FormState> = {};
    if (!form.name.trim()) errors.name = "Informe o nome da empresa.";
    if (!isValidCNPJ(form.cnpj)) errors.cnpj = "Use o formato 00.000.000/0000-00.";
    setFieldErrors(errors);
    return Object.keys(errors).length === 0;
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setFormError("");
    if (!validate()) return;

    setSaving(true);
    try {
      if (editing) {
        await api.put(`/companies/${editing.id}`, form);
        setToast({ message: "Empresa atualizada.", variant: "success" });
      } else {
        await api.post("/companies", form);
        setToast({ message: "Empresa criada.", variant: "success" });
      }
      setModalOpen(false);
      loadCompanies();
    } catch (err) {
      setFormError(extractMessage(err, "Não foi possível salvar a empresa."));
    } finally {
      setSaving(false);
    }
  }

  async function handleDelete(id: string) {
    setDeletingId(id);
    try {
      await api.delete(`/companies/${id}`);
      setToast({ message: "Empresa excluída.", variant: "success" });
      setConfirmingId(null);
      loadCompanies();
    } catch (err) {
      setToast({ message: extractMessage(err, "Não foi possível excluir a empresa."), variant: "error" });
    } finally {
      setDeletingId(null);
    }
  }

  return (
    <div className={styles.page}>
      <div className={styles.headerRow}>
        <div>
          <span className={styles.eyebrow}>Administração</span>
          <h1 className={styles.title}>Empresas</h1>
        </div>
        <button className={styles.primaryButton} type="button" onClick={openCreate}>
          Nova empresa
        </button>
      </div>

      {loading ? (
        <p className={styles.hint}>Carregando empresas…</p>
      ) : companies.length === 0 ? (
        <p className={styles.hint}>Nenhuma empresa cadastrada ainda. Crie a primeira acima.</p>
      ) : (
        <div className={styles.tableWrap}>
          <table className={styles.table}>
            <thead>
              <tr>
                <th>Nome</th>
                <th>CNPJ</th>
                <th aria-label="Ações" />
              </tr>
            </thead>
            <tbody>
              {companies.map((company) => (
                <tr key={company.id}>
                  <td>{company.name}</td>
                  <td className={styles.mono}>{company.cnpj}</td>
                  <td className={styles.actions}>
                    {confirmingId === company.id ? (
                      <span className={styles.confirmGroup}>
                        <span className={styles.confirmLabel}>Excluir?</span>
                        <button
                          className={styles.dangerLink}
                          type="button"
                          disabled={deletingId === company.id}
                          onClick={() => handleDelete(company.id)}
                        >
                          {deletingId === company.id ? "Excluindo…" : "Confirmar"}
                        </button>
                        <button
                          className={styles.link}
                          type="button"
                          onClick={() => setConfirmingId(null)}
                        >
                          Cancelar
                        </button>
                      </span>
                    ) : (
                      <span className={styles.confirmGroup}>
                        <button className={styles.link} type="button" onClick={() => openEdit(company)}>
                          Editar
                        </button>
                        <button
                          className={styles.dangerLink}
                          type="button"
                          onClick={() => setConfirmingId(company.id)}
                        >
                          Excluir
                        </button>
                      </span>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      <Pagination page={page} limit={LIMIT} total={total} onPageChange={setPage} />

      {modalOpen && (
        <Modal title={editing ? "Editar empresa" : "Nova empresa"} onClose={() => setModalOpen(false)}>
          <form className={styles.form} onSubmit={handleSubmit} noValidate>
            <label className={styles.label} htmlFor="company-name">
              Nome
            </label>
            <input
              id="company-name"
              className={styles.input}
              value={form.name}
              onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))}
              required
            />
            {fieldErrors.name && <p className={styles.fieldError}>{fieldErrors.name}</p>}

            <label className={styles.label} htmlFor="company-cnpj">
              CNPJ
            </label>
            <input
              id="company-cnpj"
              className={styles.input}
              value={form.cnpj}
              onChange={(e) => setForm((f) => ({ ...f, cnpj: formatCNPJ(e.target.value) }))}
              placeholder="00.000.000/0000-00"
              inputMode="numeric"
              required
            />
            {fieldErrors.cnpj && <p className={styles.fieldError}>{fieldErrors.cnpj}</p>}

            {formError && (
              <p className={styles.formError} role="alert">
                {formError}
              </p>
            )}

            <button className={styles.submit} type="submit" disabled={saving}>
              {saving ? "Salvando…" : "Salvar"}
            </button>
          </form>
        </Modal>
      )}

      {toast && <Toast message={toast.message} variant={toast.variant} onDismiss={() => setToast(null)} />}
    </div>
  );
}

function extractMessage(err: unknown, fallback: string): string {
  if (typeof err === "object" && err !== null && "response" in err) {
    const response = (err as { response?: { data?: { message?: string } } }).response;
    if (response?.data?.message) return response.data.message;
  }
  return fallback;
}
