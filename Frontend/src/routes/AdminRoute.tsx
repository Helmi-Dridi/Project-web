import { Navigate, Outlet } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

export default function AdminRoute() {
  const {  isAuthenticated, user, loading } = useAuth(); // 👈 include loading

  if (loading) return null; // ⏳ Wait until auth is loaded

  if (!isAuthenticated) return <Navigate to="/login" />;

  const isAdmin = user?.roles?.includes("Manager") || user?.roles?.includes("CEO");

  return isAdmin ? <Outlet /> : <Navigate to="/dashboard" />;
}
