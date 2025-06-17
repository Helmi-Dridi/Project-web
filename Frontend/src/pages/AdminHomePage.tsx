// AdminHomePage.tsx - Admin dashboard showing system stats and recent activity.
// Uses hooks for live data and presents a responsive overview layout.

import { CalendarDays, Users, Shield, Bell, MessageSquare } from "lucide-react";
import Calendar from "react-calendar";
import "react-calendar/dist/Calendar.css";
import { format } from "date-fns";
import { useAdminDetails, useStudents } from "../hooks/useAdmin";
import { useUpcomingAppointments, useAppointmentStats } from "../hooks/useBooking";
import { useNotificationCount, useNotifications } from "../hooks/useNotifications";
import { useUnreadCount } from "../hooks/useMessages";

export default function AdminHomePage() {
  const { data: admins = [], isLoading: loadingAdmins } = useAdminDetails();
  const { data: students = [], isLoading: loadingStudents } = useStudents();
  const { data: upcoming = [], isLoading: loadingUpcoming } = useUpcomingAppointments();
  const { data: stats } = useAppointmentStats();
  const { data: notificationCount } = useNotificationCount();
  const { data: rawNotifications } = useNotifications();
  // Ensure notifications is always an array to avoid runtime errors
  const notifications = Array.isArray(rawNotifications) ? rawNotifications : [];
  const { count: unreadMessages } = useUnreadCount();

  const metrics = [
    { label: "Admins", value: admins.length, icon: <Shield className="w-6 h-6 text-blue-600" /> },
    { label: "Students", value: students.length, icon: <Users className="w-6 h-6 text-green-600" /> },
    { label: "Total Appointments", value: stats?.total ?? 0, icon: <CalendarDays className="w-6 h-6 text-purple-600" /> },
    { label: "Upcoming", value: upcoming.length, icon: <CalendarDays className="w-6 h-6 text-orange-500" /> },
    { label: "Unread Messages", value: unreadMessages, icon: <MessageSquare className="w-6 h-6 text-indigo-600" /> },
    { label: "Notifications", value: notificationCount?.count ?? 0, icon: <Bell className="w-6 h-6 text-red-600" /> },
  ];

  const upcomingPreview = upcoming.slice(0, 5);
  const recentActivities = notifications.slice(0, 5);

  return (
    <div className="p-6 space-y-8">
      <h1 className="text-3xl font-bold text-gray-800">Admin Overview</h1>

      {/* Metrics */}
      <section className="grid gap-6 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3">
        {metrics.map((m) => (
          <div
            key={m.label}
            className="flex items-center bg-white shadow rounded-xl p-4 gap-4"
          >
            {m.icon}
            <div>
              <p className="text-sm text-gray-500">{m.label}</p>
              <p className="text-xl font-semibold">
                {loadingAdmins || loadingStudents || loadingUpcoming ? "..." : m.value}
              </p>
            </div>
          </div>
        ))}
      </section>

      {/* Upcoming appointments */}
      <section className="bg-white rounded-xl shadow p-6">
        <h2 className="text-lg font-semibold mb-4">Upcoming Appointments</h2>
        {loadingUpcoming ? (
          <p className="text-gray-500">Loading...</p>
        ) : upcomingPreview.length === 0 ? (
          <p className="text-gray-500">No upcoming appointments.</p>
        ) : (
          <ul className="space-y-2 text-sm">
            {upcomingPreview.map((appt) => (
              <li key={appt.id} className="flex justify-between items-center border-b pb-2 last:border-b-0">
                <span>{format(new Date(appt.date), "PPP")}</span>
                <span className="text-gray-500">{appt.timeSlot}</span>
                <span className="text-gray-500">{appt.status}</span>
              </li>
            ))}
          </ul>
        )}
      </section>

      {/* Recent activity */}
      <section className="bg-white rounded-xl shadow p-6">
        <h2 className="text-lg font-semibold mb-4">Recent Activity</h2>
        {recentActivities.length === 0 ? (
          <p className="text-gray-500">No recent notifications.</p>
        ) : (
          <ul className="space-y-2 text-sm">
            {recentActivities.map((n) => (
              <li key={n.id} className="flex justify-between items-start border-b pb-2 last:border-b-0">
                <div>
                  <p className="font-medium text-gray-700">{n.type}</p>
                  <p className="text-gray-500">{n.content}</p>
                </div>
                <span className="text-xs text-gray-400">{new Date(n.createdAt).toLocaleDateString()}</span>
              </li>
            ))}
          </ul>
        )}
      </section>

      {/* Calendar snapshot */}
      <section className="bg-white rounded-xl shadow p-6">
        <h2 className="text-lg font-semibold mb-4">Calendar</h2>
        <Calendar className="w-full" tileClassName={() => "!rounded-md"} />
      </section>
    </div>
  );
}