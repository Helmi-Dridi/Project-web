import { useState } from "react";
import { Mail, Lock, Eye, EyeOff } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { login as loginService } from "../services/authService";

export default function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [show, setShow] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const navigate = useNavigate();
  const { login: saveToken } = useAuth();

  const handleSubmit = async (e: React.FormEvent) => {
  e.preventDefault();
  setLoading(true);
  setError("");

  try {
    const response = await loginService(email, password);

    // Save token and user
    saveToken(response.accessToken, response.user);

    // Redirect based on role
    const roles = response.user.roles || [];

    if (roles.includes("Manager") || roles.includes("CEO")) {
navigate("/admin/dashboard", { replace: true });
    } else {
navigate("/dashboard", { replace: true });
    }
  } catch (err) {
    console.error("Login failed:", err);
    setError("Invalid email or password");
  } finally {
    setLoading(false);
  }
};



  return (
    <div className="min-h-screen flex flex-col md:flex-row bg-gradient-to-br from-blue-100 via-slate-100 to-white">
  {/* LEFT SIDE: Logo */}
  <div className="md:w-1/2 flex items-center justify-center bg-black p-6">
    <img
      src="/lovable-uploads/Scholarevakber.jpg"
      alt="ScholarRev full logo"
      className="max-w-[90%] h-auto"
    />
  </div>

  {/* RIGHT SIDE: Form */}
  <div className="md:w-1/2 flex items-center justify-center p-8">
    <div className="w-full max-w-md bg-white/90 backdrop-blur-md shadow-xl border border-slate-200 rounded-xl p-10 animate-fade-in">
      
      {/* Updated Heading */}
      <div className="mb-10 text-center">
        <h2 className="text-4xl font-extrabold text-slate-800 mb-2 tracking-tight">
          Welcome Back
        </h2>
        <p className="text-xl text-slate-500 font-medium">
          Please enter your credentials
        </p>
      </div>

      {/* Optional error */}
      {error && <p className="text-red-500 text-center mb-4">{error}</p>}

      {/* Form */}
      <form onSubmit={handleSubmit} className="space-y-7 text-xl">
        <div>
          <label htmlFor="email" className="block mb-2 text-xl text-slate-700 font-semibold">
            Email Address
          </label>
          <div className="relative">
            <Mail className="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400 h-6 w-6" />
            <input
              type="email"
              id="email"
              className="w-full pl-14 pr-4 py-4 text-xl rounded-xl border border-slate-300 focus:ring-2 focus:ring-blue-500 focus:outline-none transition-all"
              placeholder="you@example.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </div>
        </div>

        <div>
          <label htmlFor="password" className="block mb-2 text-xl text-slate-700 font-semibold">
            Password
          </label>
          <div className="relative">
            <Lock className="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400 h-6 w-6" />
            <input
              type={show ? "text" : "password"}
              id="password"
              className="w-full pl-14 pr-12 py-4 text-xl rounded-xl border border-slate-300 focus:ring-2 focus:ring-blue-500 focus:outline-none transition-all"
              placeholder="••••••••"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
            <button
              type="button"
              onClick={() => setShow(!show)}
              className="absolute right-4 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600"
            >
              {show ? <EyeOff size={24} /> : <Eye size={24} />}
            </button>
          </div>
        </div>

        <button
          type="submit"
          disabled={loading}
          className="w-full py-4 rounded-xl bg-blue-600 text-white text-xl font-bold hover:bg-blue-700 transition-all duration-200 flex justify-center items-center hover:scale-[1.02] disabled:opacity-50"
        >
          {loading && (
            <div className="w-6 h-6 border-2 border-white border-t-transparent rounded-full animate-spin mr-3"></div>
          )}
          {loading ? "Signing in..." : "Sign In"}
        </button>
      </form>
    </div>
  </div>
</div>

  );
}
