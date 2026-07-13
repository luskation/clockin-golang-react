import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter } from "react-router-dom";
import Login from "./Login";
import * as auth from "../services/auth";

test("renderiza campos de email e senha", () => {
  render(
    <MemoryRouter>
      <Login />
    </MemoryRouter>,
  );

  expect(screen.getByLabelText(/email/i)).toBeInTheDocument();
  expect(screen.getByLabelText(/senha/i)).toBeInTheDocument();
});

test("mostra mensagem de erro quando o login falha", async () => {
  vi.spyOn(auth, "login").mockRejectedValueOnce(new Error("credenciais inválidas"));
  const user = userEvent.setup();

  render(
    <MemoryRouter>
      <Login />
    </MemoryRouter>,
  );

  await user.type(screen.getByLabelText(/email/i), "ana@example.com");
  await user.type(screen.getByLabelText(/senha/i), "senha-errada");
  await user.click(screen.getByRole("button", { name: /entrar/i }));

  expect(await screen.findByRole("alert")).toHaveTextContent("Email ou senha incorretos.");
});
