import {
  createContext,
  useContext,
  useState,
  useEffect,
  type ReactNode,
} from 'react';

export interface UserSession {
  ID: string;
  email: string;
  name: string;
  profilePicture: string;
  workCompanyId: string;
  roles: string[];
}

interface AuthContextType {
  token: string | null;
  user: UserSession | null;
  login: (token: string, user: UserSession) => void;
  logout: (callback?: () => void) => void;
  isAuthenticated: boolean;
  setUser: React.Dispatch<React.SetStateAction<UserSession | null>>;
  loading: boolean; // ✅ New
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [token, setToken] = useState<string | null>(null);
  const [user, setUser] = useState<UserSession | null>(null);
  const [loading, setLoading] = useState(true); // ✅ New

  useEffect(() => {
  try {
    const storedToken = localStorage.getItem('token');
    const storedUser = localStorage.getItem('user');
    console.log("🔐 Loading auth from localStorage...");
    console.log("📦 Token:", storedToken);
    console.log("👤 User:", storedUser);

    if (storedToken) setToken(storedToken);
    if (storedUser && storedUser !== 'undefined') {
      const parsedUser = JSON.parse(storedUser);
      setUser(parsedUser);
      console.log("✅ User set:", parsedUser);
    }
  } catch (error) {
    console.error('❌ Error parsing auth data:', error);
  } finally {
    setLoading(false);
    console.log("⏱️ Auth loading complete");
  }
}, []);


  const login = (jwt: string, userData: UserSession) => {
  console.log("🔓 Login:", { jwt, userData });
  localStorage.setItem('token', jwt);
  localStorage.setItem('user', JSON.stringify(userData));
  setToken(jwt);
  setUser(userData);
};
const logout = (callback?: () => void) => {
  console.log("🔒 Logging out...");
  localStorage.removeItem('token');
  localStorage.removeItem('user');
  setToken(null);
  setUser(null);
  if (callback) callback(); // ✅ invoke navigate here
};

  const isAuthenticated = !!token;

  return (
    <AuthContext.Provider
      value={{ token, user, login, logout, isAuthenticated, setUser, loading }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) throw new Error('useAuth must be used within AuthProvider');
  return context;
};
