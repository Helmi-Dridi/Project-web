import { Routes, Route, Navigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext"; // ✅ Import auth

import PrivateRoute from "./PrivateRoute";
import AdminRoute from "./AdminRoute";

import DashboardLayout from "../components/DashboardLayout";
import AdminDashboardLayout from "../components/AdminDashboardLayout";

import Login from "../pages/Login";
import Dashboard from "../pages/Dashboard";
import Settings from "../pages/Settings";
import Steps from "../pages/Steps";
import Profile from "../pages/Profile/Profile";
import Documents from "../pages/Documents";
import Appointments from "../pages/Appointments";
import Notifications from "../pages/Profile/Notifications";
import Referrals from "../pages/Profile/Referrals";
import AdminDashboard from "../pages/admindashboard";
import AdminList from "../pages/AdminList";
import AdminStudentList from "../pages/AdminStudentList";
import UniversityListPage from "../pages/AdminUniversity";
import AdminCalendarPage from "../pages/AdminCalendarPage";
import ChatPage from "../pages/ChatPage";
import AdminChatPage from "../pages/AdminChatPage";

export default function AppRoutes() {
  const { user,isAuthenticated } = useAuth(); // ✅ Check auth

  return (
    <Routes>
      {/* Public */}
      <Route
  path="/login"
  element={
    isAuthenticated
      ? user?.roles?.includes("Manager") || user?.roles?.includes("CEO")
        ? <Navigate to="/admin/dashboard" />
        : <Navigate to="/dashboard" />
      : <Login />
  }
/>


      {/* Admin-only routes */}
      <Route element={<AdminRoute />}>
        <Route element={<AdminDashboardLayout />}>
          <Route path="/admin/dashboard" element={<AdminDashboard />} />
          <Route path="/admin/list" element={<AdminList />} />
          <Route path="/Student/list" element={<AdminStudentList />} />
          <Route path="/university/list" element={<UniversityListPage />} />
          <Route path="/admin/calendar" element={<AdminCalendarPage />} />
          <Route path="/admin/chat/:userId?" element={<AdminChatPage />} />

          
        </Route>
      </Route>

      {/* Protected */}
      <Route element={<PrivateRoute />}>
        <Route element={<DashboardLayout />}>
          <Route path="/dashboard" element={<Dashboard />} />
          <Route path="/steps" element={<Steps />} />
          <Route path="/profile" element={<Profile />} />
          <Route path="/documents" element={<Documents />} />
          <Route path="/appointments" element={<Appointments />} />
          <Route path="/notifications" element={<Notifications />} />
          <Route path="/referrals" element={<Referrals />} />
          <Route path="/settings" element={<Settings />} />
          <Route path="/chat/:userId?" element={<ChatPage />} />
        </Route>
      </Route>

      {/* Catch-all */}
      <Route path="*" element={<Navigate to="/login" />} />
    </Routes>
  );
}
