
/**
 * Modernized top navigation bar displaying search, notifications, and user menu.
 * UI uses Tailwind for soft shadows, hover effects and animated badges while
 * preserving all existing logic.
 */
import { Bell, ChevronDown } from "lucide-react";
import { useState, useEffect, useRef } from "react";
import { useAuth } from "../context/AuthContext";
import { useNotifications } from "../hooks/useNotifications";
import { useNavigate } from "react-router-dom";

type TopbarProps = {
  collapsed: boolean;
  setCollapsed: React.Dispatch<React.SetStateAction<boolean>>;
};

export default function Topbar({ collapsed }: TopbarProps) {
  const [showNotifications, setShowNotifications] = useState(false);
  const [showUserMenu, setShowUserMenu] = useState(false);
  const navigate = useNavigate();

  const { user, logout } = useAuth();
  const profileImage = user?.profilePicture || "/default-avatar.png";
  const username = user?.name || "Unknown User";

  const { data: rawNotifications, isLoading, isError } = useNotifications();
  const notifications = Array.isArray(rawNotifications) ? rawNotifications : [];

  const notifRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (notifRef.current && !notifRef.current.contains(event.target as Node)) {
        setShowNotifications(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const unseenCount = notifications.filter((n) => !n.seen).length;

  return (
    <header className="fixed top-0 left-0 right-0 z-30 h-20 bg-white shadow-md border-b flex items-center justify-between px-6 py-3 transition-all">
      {/* Search */}
      <div
        className={`flex items-center transition-all duration-300 ${
          collapsed ? "ml-20" : "ml-64"
        }`}
      >
        <input
          type="text"
          placeholder="Search ScholarRev..."
          className="w-48 sm:w-60 md:w-72 px-4 py-2 rounded-full border border-slate-300 shadow-sm focus:ring-2 focus:ring-blue-500 focus:outline-none transition"
        />
      </div>

      {/* Right Icons */}
      <div className="flex items-center gap-4">
        {/* Notification */}
        <div className="relative" ref={notifRef}>
          <button
            className="relative text-slate-600 hover:text-blue-600 transition"
            onClick={() => setShowNotifications((prev) => !prev)}
          >
            <Bell className="w-6 h-6" />
            {unseenCount > 0 && (
              <span className="absolute -top-1 -right-1 flex">
                <span className="animate-ping absolute inline-flex h-2 w-2 rounded-full bg-red-500 opacity-75" />
                <span className="relative inline-flex rounded-full h-2 w-2 bg-red-600" />
              </span>
            )}
          </button>

          {showNotifications && (
            <div className="absolute right-0 mt-2 w-80 bg-white border rounded-xl shadow-lg p-4 z-50 max-h-96 overflow-y-auto animate-fade-in">
              <h4 className="font-medium mb-2">🔔 Notifications</h4>
              {isLoading ? (
                <p className="text-sm text-gray-500">Loading...</p>
              ) : isError ? (
                <p className="text-sm text-red-500">Failed to load notifications.</p>
              ) : notifications.length === 0 ? (
                <p className="text-sm text-gray-500">No new notifications.</p>
              ) : (
                <ul className="space-y-2">
                  {notifications.slice(0, 5).map((notif) => (
                    <li
                      key={notif.id}
                      className="text-sm bg-gray-50 rounded-lg p-3 hover:bg-gray-100 transition"
                    >
                      <p className="font-semibold">{notif.type}</p>
                      <p className="text-gray-600">{notif.content}</p>
                      <p className="text-xs text-gray-400 mt-1">
                        {new Date(notif.createdAt).toLocaleString()}
                      </p>
                    </li>
                  ))}
                </ul>
              )}
            </div>
          )}
        </div>

        {/* User Menu */}
        <div className="relative">
          <button
            className="flex items-center gap-2 hover:bg-slate-100 px-3 py-2 rounded-lg transition"
            onClick={() => setShowUserMenu(!showUserMenu)}
          >
            <img
              src={profileImage}
              alt="User"
              className="w-10 h-10 rounded-full object-cover border"
            />
            <span className="hidden md:inline text-slate-700 font-medium">
              {username}
            </span>
            <ChevronDown className="w-4 h-4 text-slate-600 hidden md:inline" />
          </button>
          {showUserMenu && (
            <div className="absolute right-0 mt-2 w-40 bg-white shadow-lg border rounded-xl py-2 z-10 animate-fade-in">
              <ul className="text-sm">
                <li className="px-4 py-2 hover:bg-slate-100 cursor-pointer">Profile</li>
                <li className="px-4 py-2 hover:bg-slate-100 cursor-pointer">Settings</li>
                <li
                  className="px-4 py-2 hover:bg-slate-100 cursor-pointer"
                  onClick={() => logout(() => navigate("/login"))}
                >
                  Logout
                </li>
              </ul>
            </div>
          )}
        </div>
      </div>
    </header>
  );
}
