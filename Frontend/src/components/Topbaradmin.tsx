/**
 * Enhanced admin top bar with animated notification badge, full notification
 * history and delete controls, while keeping existing functionality intact.
 */
import { Bell, Trash2, ChevronDown } from "lucide-react";
import { useState, useRef, useEffect } from "react";
import { useAuth } from "../context/AuthContext";
import {
  useNotifications,
  useDeleteNotification,
} from "../hooks/useNotifications";
import { useNavigate } from "react-router-dom";

type TopbarProps = {
  collapsed: boolean;
  setCollapsed: React.Dispatch<React.SetStateAction<boolean>>;
};

export default function Topbaradmin({ collapsed }: TopbarProps) {
  const [showNotifications, setShowNotifications] = useState(false);
  const [showUserMenu, setShowUserMenu] = useState(false);
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const profileImage = user?.profilePicture || "/default-avatar.png";
  const username = user?.name || "Unknown User";

  const { data: notificationsRaw, isLoading, isError, refetch } =
    useNotifications();
  const deleteMutation = useDeleteNotification();

  const rawArray = notificationsRaw
    ? Array.isArray(notificationsRaw)
      ? notificationsRaw
      : [notificationsRaw]
    : [];

  const notifications = rawArray
    .map((n: any) => ({
      id: n.ID || n.id,
      type: n.Type || n.type,
      content: n.Content || n.content,
      seen: n.Seen ?? n.seen,
      createdAt: n.CreatedAt || n.createdAt,
    }))
    .filter((n) => !isNaN(new Date(n.createdAt).getTime()))
    .sort(
      (a, b) =>
        new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
    );

  const unseenCount = notifications.filter((n) => !n.seen).length;

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

  useEffect(() => {
    if (showNotifications) {
      refetch(); // ✅ Refresh when dropdown opens
    }
  }, [showNotifications, refetch]);

  const handleDelete = async (id: string) => {
    try {
      await deleteMutation.mutateAsync(id);
    } catch (err) {
      console.error("❌ Failed to delete notification:", err);
    }
  };

  return (
    <header className="fixed top-0 left-0 right-0 z-30 h-20 bg-white shadow-md border-b flex items-center justify-between px-6 py-3 transition-all">
      {/* Search */}
      <div className={`flex items-center transition-all duration-300 ${collapsed ? "ml-20" : "ml-64"}`}>
        <input
          type="text"
          placeholder="Search ScholarRev..."
          className="w-48 sm:w-60 md:w-72 px-4 py-2 rounded-full border border-slate-300 shadow-sm focus:ring-2 focus:ring-blue-500 focus:outline-none transition"
        />
      </div>

      {/* Right Icons */}
      <div className="flex items-center gap-4">
        {/* Notifications */}
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
            <div className="absolute right-0 mt-2 w-96 bg-white border rounded-xl shadow-lg p-4 z-50 max-h-96 overflow-y-auto animate-fade-in">
              <h4 className="font-medium mb-2">🔔 Notifications</h4>
              {isLoading ? (
                <p className="text-sm text-gray-500">Loading...</p>
              ) : isError ? (
                <p className="text-sm text-red-500">Error loading notifications.</p>
              ) : notifications.length === 0 ? (
                <p className="text-sm text-gray-500">No notifications.</p>
              ) : (
                <ul className="space-y-3">
                  {notifications.map((notif) => (
                    <li
                      key={notif.id}
                      className="text-sm bg-gray-50 rounded-lg p-3 flex justify-between items-start gap-3 hover:bg-gray-100 transition"
                    >
                      <div className="flex-1">
                        <p className="font-semibold text-blue-600">{notif.type}</p>
                        <p className="text-gray-700 mt-1">{notif.content}</p>
                        <p className="text-xs text-gray-400 mt-1">
                          {notif.createdAt
                            ? new Date(notif.createdAt).toLocaleString()
                            : "Unknown date"}
                        </p>
                      </div>
                      <button
                        onClick={() => handleDelete(notif.id)}
                        className="text-slate-400 hover:text-red-600"
                        title="Delete notification"
                      >
                        <Trash2 className="w-4 h-4" />
                      </button>
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
            <span className="hidden md:inline text-slate-700 font-medium">{username}</span>
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