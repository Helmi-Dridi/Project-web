// src/components/SidebarAdmin.tsx
import {
  FaUserShield,
  FaUserGraduate,
  FaUniversity,
  FaCalendarAlt,
  FaSignOutAlt,
  FaCog,
  FaHome,
} from "react-icons/fa";
import type { FC } from "react";
import { useAuth } from "../context/AuthContext";
import { useNavigate, useLocation } from "react-router-dom";
import classNames from "classnames";

type SidebarProps = {
  collapsed: boolean;
  setCollapsed: React.Dispatch<React.SetStateAction<boolean>>;
};

const SidebarAdmin: FC<SidebarProps> = ({ collapsed, setCollapsed }) => {
  const { logout } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  const isActive = (path: string) => location.pathname === path;

  const menuItems = [
    { icon: <FaHome/>, label: "Home", path: "/admin/home" },
    { icon: <FaUserShield />, label: "Admins", path: "/admin/list" },
    { icon: <FaUserGraduate />, label: "Students", path: "/Student/list" },
    { icon: <FaUniversity />, label: "Universities", path: "/university/list" },
    { icon: <FaCalendarAlt />, label: "Calendar", path: "/admin/calendar" },
    { icon: <FaCalendarAlt />, label: "Chat", path: "/admin/chat" },
    
  ];

  return (
    <aside
      className={classNames(
        "fixed top-0 left-0 h-screen bg-black text-white z-40 shadow-md flex flex-col transition-all duration-300 ease-in-out",
        {
          "w-64": !collapsed,
          "w-16": collapsed,
        }
      )}
    >
      {/* Toggle */}
      <button
        onClick={() => setCollapsed(!collapsed)}
        className="absolute -right-4 top-6 w-8 h-8 bg-white text-black border rounded-full shadow-lg flex items-center justify-center z-50"
        title="Toggle Sidebar"
      >
        <svg
          className={classNames("w-5 h-5 transition-transform", {
            "rotate-180": collapsed,
          })}
          fill="none"
          stroke="currentColor"
          strokeWidth={2}
          viewBox="0 0 24 24"
        >
          <path strokeLinecap="round" strokeLinejoin="round" d="M9 5l7 7-7 7" />
        </svg>
      </button>

      {/* Logo */}
      <div className="flex justify-center py-6 border-b border-slate-800 transition-all duration-300">
        {!collapsed && (
          <img
            src="/lovable-uploads/Scholarevakber.jpg"
            alt="Logo"
            className="h-16 object-contain"
          />
        )}
      </div>

      {/* Menu */}
      <nav className="flex flex-col gap-2 px-2 mt-4">
        {menuItems.map((item, idx) => (
          <SidebarItem
            key={idx}
            icon={item.icon}
            label={item.label}
            onClick={() => navigate(item.path)}
            active={isActive(item.path)}
            collapsed={collapsed}
          />
        ))}
      </nav>

      {/* Bottom */}
      <div className="mt-auto px-2 mb-4 flex flex-col gap-2">
        <SidebarItem
          icon={<FaCog />}
          label="Settings"
          onClick={() => navigate("/admin/settings")}
          collapsed={collapsed}
          active={isActive("/admin/settings")}
        />
        <SidebarItem
          icon={<FaSignOutAlt />}
          label="Logout"
          onClick={() => logout(() => navigate("/login"))} // ✅ now navigates to login
          collapsed={collapsed}
        />
      </div>
    </aside>
  );
};

type SidebarItemProps = {
  icon: React.ReactNode;
  label: string;
  onClick: () => void;
  active?: boolean;
  collapsed: boolean;
};

const SidebarItem: FC<SidebarItemProps> = ({
  icon,
  label,
  onClick,
  active,
  collapsed,
}) => (
  <button
    onClick={onClick}
    className={classNames(
      "flex items-center w-full px-4 py-3 text-left rounded-md hover:bg-slate-800 transition-colors",
      {
        "bg-slate-800": active,
        "justify-center": collapsed,
      }
    )}
  >
    <span className="text-xl">{icon}</span>
    {!collapsed && <span className="ml-4 text-sm font-medium">{label}</span>}
  </button>
);

export default SidebarAdmin;
