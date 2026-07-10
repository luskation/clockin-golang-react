import { Navigate, Outlet } from "react-router-dom";
import { getCurrentUser } from "../services/auth";

export default function AdminRoute() {
  const user = getCurrentUser();
  return user?.role === "admin" ? <Outlet /> : <Navigate to="/" replace />;
}
