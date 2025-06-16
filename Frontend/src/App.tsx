import AppRoutes from "./routes";
import { Toaster } from "react-hot-toast";
import { useAuth } from "./context/AuthContext";
import { useEffect } from "react";

function App() {
  const { user, token, isAuthenticated, loading } = useAuth();

  useEffect(() => {
    console.log("🧩 App mounted");
    console.log("🔐 Token:", token);
    console.log("👤 User:", user);
    console.log("✅ isAuthenticated:", isAuthenticated);
    console.log("⏳ Auth loading:", loading);
  }, [token, user, isAuthenticated, loading]);

  return (
    <>
      <AppRoutes />
      <Toaster position="top-right" />
    </>
  );
}

export default App;
