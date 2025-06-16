import axios from "axios";

export type UserSession = {
  ID: string;
  email: string;
  name: string;
  profilePicture: string;
  workCompanyId: string;
  roles: string[]; // ✅ Add this

};

type LoginResponse = {
  accessToken: string;
  user: UserSession;
};

export const login = async (email: string, password: string): Promise<LoginResponse> => {
  const res = await axios.post("http://localhost:8080/api/auth/signin", {
    email,
    password,
  });

  const data = res.data.data;
  const user = data.user;

  // ✅ Normalize profile picture URL if needed
  if (user.profilePicture && !user.profilePicture.startsWith("http")) {
    user.profilePicture = `http://localhost:8080/${user.profilePicture.replace(/\\/g, "/")}`;
  }

  // ✅ Ensure roles exist and are array
  if (!Array.isArray(user.roles)) {
    user.roles = [];
  }

  localStorage.setItem("token", data.accessToken);
  localStorage.setItem("user", JSON.stringify(user));

  return {
    accessToken: data.accessToken,
    user,
  };
};




export const getCurrentUser = (): UserSession | null => {
  try {
    const userStr = localStorage.getItem("user");
    if (!userStr || userStr === "undefined") return null;
    return JSON.parse(userStr);
  } catch (error) {
    console.error("Failed to parse user from localStorage:", error);
    return null;
  }
};

export const getAuthHeaders = () => {
  const token = localStorage.getItem("token");
  return {
    Authorization: `Bearer ${token}`,
  };
};
