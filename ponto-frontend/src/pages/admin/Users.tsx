import { useEffect, useState, type FormEvent } from "react";
import api from "../../services/api";
import Modal from "../../components/Modal";
import Toast from "../../components/Toast";
import Pagination from "../../components/Pagination";
import { getCurrentUser } from "../../services/auth";
import { isValidEmail, isValidPassword } from "../../utils/validation";
import styles from "./Users.module.css";

type Role = "admin" | "employee";

interface Company {
  id: string;
  name: string;
}

interface User {
  id: string;
  company_id: string;
  name: string;
  email: string;
  role: Role;
}

interface FormState {
  name: string;
  email: string;
  password: string;
  role: Role;
  company_id: string;
}

const emptyForm: FormState = { name: "", email: "", password: "", role: "employee", company_id: "" };
const LIMIT = 10;

export default function Users() {
  const currentUser = getCurrentUser();

  const [users, setUsers] = useState<User[]>([]);
  const [companies, setCompanies] = useState<Company[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(true);
  const [toast, setToast] = useState<{ message: string; variant: "success" | "error" } | null>(null);

  const [modalOpen, setModalOpen] = useState(false);
  const [editing, setEditing] = useState<User | null>(null);
  const [form, setForm] = useState<FormState>(emptyForm);
  const [fieldErrors, setFieldErrors] = useState<Partial<Record<keyof FormState, string>>>({});
  const [formError, setFormError] = useState("");
  const [saving, setSaving] = useState(false);

  const [confirmingId, setConfirmingId] = useState<string | null>(null);
  const [deletingId, setDeletingId] = useState<string | null>(null);

  function loadUsers() {
    setLoading(true);
    api
      .get("/users", { params: { page, limit: LIMIT } })
      .then((res) => {
        setUsers(res.data.data ?? []);
        setTotal(res.data.total ?? 0);
      })
      .catch(() => setToast({ message: "Não foi possível carregar os usuários.", variant: "error" }))
      .finally(() => setLoading(false));
  }

  useEffect(loadUsers, [page]);

  useEffect(() => {
    api
      .get("/companies", { params: { page: 1, limit: 100 } })
      .then((res) => setCompanies(res.data.data ?? []))
      .catch(() => setToast({ message: "Não foi possível carregar as empresas.", variant: "error" }));
  }, []);

  function companyName(companyId: string): string {
    return companies.find((c) => c.id === companyId)?.name ?? "—";
  }

  function openCreate() {
    setEditing(null);
    setForm({ ...emptyForm, company_id: companies[0]?.id ?? "" });
    setFieldErrors({});
    setFormError("");
    setModalOpen(true);
  }

  function openEdit(user: User) {
    setEditing(user);
    setForm({ name: user.name, email: user.email, password: "", role: user.role, company_id: user.company_id });
    setFieldErrors({});
    setFormError("");
    setModalOpen(true);
  }

  function validate(): boolean {
    const errors: Partial<Record<keyof FormState, string>> = {};
    if (!form.name.trim()) errors.name = "Informe o nome.";
    if (!isValidEmail(form.email)) errors.email = "Informe um e-mail válido.";
    if (!editing) {
      if (!form.company_id) errors.company_id = "Selecione uma empresa.";
      if (!isValidPassword(form.password)) errors.password = "A senha deve ter ao menos 6 caracteres.";
    }
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
        await api.put(`/users/${editing.id}`, { name: form.name, email: form.email, role: form.role });
        setToast({ message: "Usuário atualizado.", variant: "success" });
      } else {
        await api.post("/users", form);
        setToast({ message: "Usuário criado.", variant: "success" });
      }
      setModalOpen(false);
      loadUsers();
    } catch (err) {
      setFormError(extractMessage(err, "Não foi possível salvar o usuário."));
    } finally {
      setSaving(false);
    }
  }

  async function handleDelete(id: string) {
    setDeletingId(id);
    try {
      await api.delete(`/users/${id}`);
      setToast({ message: "Usuário excluído.", variant: "success" });
      setConfirmingId(null);
      loadUsers();
    } catch (err) {
      setToast({ message: extractMessage(err, "Não foi possível excluir o usuário."), variant: "error" });
    } finally {
      setDeletingId(null);
    }
  }

  return (
    <div className={styles.page}>
      <div className={styles.headerRow}>
        <div>
          <span className={styles.eyebrow}>Administração</span>
          <h1 className={styles.title}>Usuários</h1>
        </div>
        <button className={styles.primaryButton} type="button" onClick={openCreate}>
          Novo usuário
        </button>
      </div>

      {loading ? (
        <p className={styles.hint}>Carregando usuários…</p>
      ) : users.length === 0 ? (
        <p className={styles.hint}>Nenhum usuário cadastrado ainda. Crie o primeiro acima.</p>
      ) : (
        <div className={styles.tableWrap}>
          <table className={styles.table}>
            <thead>
              <tr>
                <th>Nome</th>
                <th>E-mail</th>
                <th>Empresa</th>
                <th>Papel</th>
                <th aria-label="Ações" />
              </tr>
            </thead>
            <tbody>
              {users.map((user) => (
                <tr key={user.id}>
                  <td>{user.name}</td>
                  <td>{user.email}</td>
                  <td>{companyName(user.company_id)}</td>
                  <td>
                    <span className={styles.roleBadge} data-role={user.role}>
                      {user.role === "admin" ? "Administrador" : "Colaborador"}
                    </span>
                  </td>
                  <td className={styles.actions}>
                    {confirmingId === user.id ? (
                      <span className={styles.confirmGroup}>
                        <span className={styles.confirmLabel}>Excluir?</span>
                        <button
                          className={styles.dangerLink}
                          type="button"
                          disabled={deletingId === user.id}
                          onClick={() => handleDelete(user.id)}
                        >
                          {deletingId === user.id ? "Excluindo…" : "Confirmar"}
                        </button>
                        <button className={styles.link} type="button" onClick={() => setConfirmingId(null)}>
                          Cancelar
                        </button>
                      </span>
                    ) : (
                      <span className={styles.confirmGroup}>
                        <button className={styles.link} type="button" onClick={() => openEdit(user)}>
                          Editar
                        </button>
                        {currentUser?.id !== user.id && (
                          <button
                            className={styles.dangerLink}
                            type="button"
                            onClick={() => setConfirmingId(user.id)}
                          >
                            Excluir
                          </button>
                        )}
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
        <Modal title={editing ? "Editar usuário" : "Novo usuário"} onClose={() => setModalOpen(false)}>
          <form className={styles.form} onSubmit={handleSubmit} noValidate>
            <label className={styles.label} htmlFor="user-name">
              Nome
            </label>
            <input
              id="user-name"
              className={styles.input}
              value={form.name}
              onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))}
              required
            />
            {fieldErrors.name && <p className={styles.fieldError}>{fieldErrors.name}</p>}

            <label className={styles.label} htmlFor="user-email">
              E-mail
            </label>
            <input
              id="user-email"
              className={styles.input}
              type="email"
              value={form.email}
              onChange={(e) => setForm((f) => ({ ...f, email: e.target.value }))}
              required
            />
            {fieldErrors.email && <p className={styles.fieldError}>{fieldErrors.email}</p>}

            {/* Senha e empresa só aparecem na criação: o endpoint de update
                no backend não altera esses dois campos, só nome/e-mail/papel. */}
            {!editing && (
              <>
                <label className={styles.label} htmlFor="user-password">
                  Senha
                </label>
                <input
                  id="user-password"
                  className={styles.input}
                  type="password"
                  value={form.password}
                  onChange={(e) => setForm((f) => ({ ...f, password: e.target.value }))}
                  required
                />
                {fieldErrors.password && <p className={styles.fieldError}>{fieldErrors.password}</p>}

                <label className={styles.label} htmlFor="user-company">
                  Empresa
                </label>
                <select
                  id="user-company"
                  className={styles.select}
                  value={form.company_id}
                  onChange={(e) => setForm((f) => ({ ...f, company_id: e.target.value }))}
                  required
                >
                  <option value="" disabled>
                    Selecione…
                  </option>
                  {companies.map((company) => (
                    <option key={company.id} value={company.id}>
                      {company.name}
                    </option>
                  ))}
                </select>
                {fieldErrors.company_id && <p className={styles.fieldError}>{fieldErrors.company_id}</p>}
              </>
            )}

            <label className={styles.label} htmlFor="user-role">
              Papel
            </label>
            <select
              id="user-role"
              className={styles.select}
              value={form.role}
              onChange={(e) => setForm((f) => ({ ...f, role: e.target.value as Role }))}
            >
              <option value="employee">Colaborador</option>
              <option value="admin">Administrador</option>
            </select>

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
