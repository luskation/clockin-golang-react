import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import Pagination from "./Pagination";

describe("Pagination", () => {
  it("não renderiza nada quando tudo cabe em uma página", () => {
    const { container } = render(<Pagination page={1} limit={10} total={5} onPageChange={() => {}} />);

    expect(container).toBeEmptyDOMElement();
  });

  it("desabilita 'Anterior' na primeira página e habilita 'Próxima'", () => {
    render(<Pagination page={1} limit={10} total={25} onPageChange={() => {}} />);

    expect(screen.getByText("Anterior")).toBeDisabled();
    expect(screen.getByText("Próxima")).toBeEnabled();
  });

  it("desabilita 'Próxima' na última página", () => {
    render(<Pagination page={3} limit={10} total={25} onPageChange={() => {}} />);

    expect(screen.getByText("Próxima")).toBeDisabled();
    expect(screen.getByText("Anterior")).toBeEnabled();
  });

  it("chama onPageChange com a página seguinte ao clicar em 'Próxima'", async () => {
    const onPageChange = vi.fn();
    const user = userEvent.setup();
    render(<Pagination page={2} limit={10} total={25} onPageChange={onPageChange} />);

    await user.click(screen.getByText("Próxima"));

    expect(onPageChange).toHaveBeenCalledWith(3);
  });

  it("mostra a página atual e o total de páginas", () => {
    render(<Pagination page={2} limit={10} total={25} onPageChange={() => {}} />);

    expect(screen.getByText("Página 2 de 3")).toBeInTheDocument();
  });
});
