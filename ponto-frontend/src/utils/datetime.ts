export function formatTime(date: Date): string {
  return date.toLocaleTimeString("pt-BR", { hour12: false });
}

export function formatDate(date: Date): string {
  const formatted = date.toLocaleDateString("pt-BR", {
    weekday: "long",
    day: "2-digit",
    month: "long",
  });
  return formatted.charAt(0).toUpperCase() + formatted.slice(1);
}
