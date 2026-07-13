import { useEffect } from "react";
import styles from "./Toast.module.css";

type ToastVariant = "success" | "error";

interface ToastProps {
  message: string;
  variant: ToastVariant;
  onDismiss: () => void;
}

export default function Toast({ message, variant, onDismiss }: ToastProps) {
  useEffect(() => {
    const id = setTimeout(onDismiss, 4000);
    return () => clearTimeout(id);
  }, [message, onDismiss]);

  return (
    <div className={styles.toast} data-variant={variant} role="status">
      {message}
      <button
        className={styles.dismiss}
        type="button"
        onClick={onDismiss}
        aria-label="Fechar aviso"
      >
        ×
      </button>
    </div>
  );
}
