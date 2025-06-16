import { Navigate, Outlet } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

export default function PrivateRoute() {
  const { isAuthenticated, loading } = useAuth();

  console.log("🔐 PrivateRoute:", {
    loading,
    isAuthenticated,
  });

  if (loading) return null;

  return isAuthenticated ? <Outlet /> : <Navigate to="/login" />;
}
