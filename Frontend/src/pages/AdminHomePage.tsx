// AdminHomePage.tsx
// Dashboard using data from useAdminDetails, useStudents, useUpcomingAppointments,
// useAppointmentStats, useNotificationCount, useNotifications, useUnreadCount,
// useUniversities, useUniversityPrograms and useAllStudentEnrollments.
// Displays hero stats, a university snapshot, metrics grid, upcoming appointments,
// recent activity feed and progress bars for appointment statistics.

import {
  CalendarDays,
  Users,
  Shield,
  Bell,
  MessageSquare,
  CheckCircle2,
  XCircle,
  GraduationCap,
  Building2,
} from "lucide-react";
import { format } from "date-fns";
import { useAdminDetails, useStudents } from "../hooks/useAdmin";
import { useUpcomingAppointments, useAppointmentStats } from "../hooks/useBooking";
import { useNotificationCount, useNotifications } from "../hooks/useNotifications";
import { useUnreadCount } from "../hooks/useMessages";
import { useAllStudentEnrollments } from "../hooks/useEnrollment";
import { useUniversities, useUniversityPrograms } from "../hooks/useUniversity";

export default function AdminHomePage() {
  const { data: admins = [], isLoading: loadingAdmins } = useAdminDetails();
  const { data: students = [], isLoading: loadingStudents } = useStudents();
  const { data: upcoming = [], isLoading: loadingUpcoming } = useUpcomingAppointments();
  const { data: stats, isLoading: loadingStats } = useAppointmentStats();
  const { data: notificationCount } = useNotificationCount();
  const { data: rawNotifications } = useNotifications();
  const { count: unreadMessages } = useUnreadCount();
  const { data: enrollments = [], isLoading: loadingEnrollments } = useAllStudentEnrollments();
  const { data: universityData } = useUniversities(1, 1);
  const firstUniversity = universityData?.items?.[0];
  const { data: programs = [], isLoading: loadingPrograms } = useUniversityPrograms(firstUniversity?.id ?? "");
  // Ensure notifications is always an array to avoid runtime errors
  const notifications = Array.isArray(rawNotifications) ? rawNotifications : [];

  const metricsLoading =
    loadingAdmins ||
    loadingStudents ||
    loadingUpcoming ||
    loadingStats ||
    loadingEnrollments ||
    loadingPrograms;

  const completedRate = stats ? Math.round((stats.completed / stats.total) * 100) : 0;
  const canceledRate = stats ? Math.round((stats.canceled / stats.total) * 100) : 0;

  const heroStats = [
    {
      label: "Total Appointments",
      value: stats?.total ?? 0,
      icon: <CalendarDays className="w-8 h-8" />,
    },
    {
      label: "Completed",
      value: stats?.completed ?? 0,
      icon: <CheckCircle2 className="w-8 h-8" />,
    },
    {
      label: "Canceled",
      value: stats?.canceled ?? 0,
      icon: <XCircle className="w-8 h-8" />,
    },
    {
      label: "Total Enrollments",
      value: enrollments.length,
      icon: <GraduationCap className="w-8 h-8" />,
    },
  ];

  const metrics = [
    {
      label: "Admins",
      value: admins.length,
      icon: <Shield className="w-6 h-6 text-blue-600" />,
    },
    {
      label: "Students",
      value: students.length,
      icon: <Users className="w-6 h-6 text-green-600" />,
    },
    {
      label: "Upcoming",
      value: upcoming.length,
      icon: <CalendarDays className="w-6 h-6 text-orange-500" />,
    },
    {
      label: "Notifications",
      value: notificationCount?.count ?? 0,
      icon: <Bell className="w-6 h-6 text-red-600" />,
    },
    {
      label: "Unread Messages",
      value: unreadMessages,
      icon: <MessageSquare className="w-6 h-6 text-indigo-600" />,
    },
  ];

  const upcomingPreview = upcoming.slice(0, 5);
  const recentActivities = notifications.slice(0, 5);

  return (
    <div className="p-6 space-y-8">
      <h1 className="text-3xl font-bold text-gray-800">Admin Dashboard</h1>

      {/* Hero Stats */}
      <section className="grid gap-4 grid-cols-1 sm:grid-cols-3 animate-fade-in">
        {heroStats.map((s) => (
          <div
            key={s.label}
            className="flex items-center gap-4 p-6 rounded-2xl shadow-lg bg-gradient-to-r from-indigo-500 to-purple-600 text-white transition-transform duration-300 hover:scale-[1.02]"
          >
            {s.icon}
            <div>
              <p className="text-sm opacity-80">{s.label}</p>
              <p className="text-3xl font-bold">
                {metricsLoading ? "..." : s.value}
              </p>
            </div>
          </div>
        ))}
      </section>

      {/* Website Snapshot */}
      <section className="bg-white rounded-xl shadow p-6 animate-fade-in">
        <h2 className="text-lg font-semibold mb-4 flex items-center gap-2">
          <Building2 className="text-indigo-500" /> Website Snapshot
        </h2>
        {firstUniversity ? (
          <div className="grid sm:grid-cols-2 gap-2 text-sm">
            <p>
              <span className="font-medium text-gray-700">Name:</span>{" "}
              {firstUniversity.name}
            </p>
            <p>
              <span className="font-medium text-gray-700">Type:</span>{" "}
              {firstUniversity.universityType}
            </p>
            <p>
              <span className="font-medium text-gray-700">Location:</span>{" "}
              {firstUniversity.city}, {firstUniversity.country}
            </p>
            <p>
              <span className="font-medium text-gray-700">Created:</span>{" "}
              {format(new Date(firstUniversity.createdAt), "PPP")}
            </p>
            <p>
              <span className="font-medium text-gray-700">Programs:</span>{" "}
              {programs.length}
            </p>
          </div>
        ) : (
          <p className="text-gray-500 text-sm">No university data.</p>
        )}
      </section>

      {/* Metrics */}
      <section className="grid gap-6 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 animate-fade-in">
        {metrics.map((m) => (
          <div
            key={m.label}
            className="flex items-center bg-white shadow rounded-xl p-4 gap-4 transition-transform duration-300 hover:scale-[1.02]"
          >
            {m.icon}
            <div>
              <p className="text-sm text-gray-500">{m.label}</p>
              <p className="text-xl font-semibold">
                {metricsLoading ? "..." : m.value}
              </p>
            </div>
          </div>
        ))}
      </section>

      {/* Upcoming appointments */}
      <section className="bg-white rounded-xl shadow p-6 animate-fade-in">
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
      <section className="bg-white rounded-xl shadow p-6 animate-fade-in">
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

      {/* Statistics & Insights */}
      <section className="bg-white rounded-xl shadow p-6 space-y-4 animate-fade-in">
        <h2 className="text-lg font-semibold">Statistics &amp; Insights</h2>
        {stats ? (
          <div className="grid gap-6 sm:grid-cols-2">
            <div>
              <p className="flex items-center gap-2 font-medium text-gray-700 mb-2">
                <CheckCircle2 className="text-green-500" /> Completed Appointments
              </p>
              <div className="w-full bg-gray-200 rounded-full h-2">
                <div
                  className="bg-green-500 h-2 rounded-full transition-[width] duration-700"
                  style={{ width: `${completedRate}%` }}
                />
              </div>
              <p className="text-sm text-gray-500 mt-1">
                {stats.completed} of {stats.total} ({completedRate}%)
              </p>
            </div>

            <div>
              <p className="flex items-center gap-2 font-medium text-gray-700 mb-2">
                <XCircle className="text-red-500" /> Canceled Appointments
              </p>
              <div className="w-full bg-gray-200 rounded-full h-2">
                <div
                  className="bg-red-500 h-2 rounded-full transition-[width] duration-700"
                  style={{ width: `${canceledRate}%` }}
                />
              </div>
              <p className="text-sm text-gray-500 mt-1">
                {stats.canceled} of {stats.total} ({canceledRate}%)
              </p>
            </div>
          </div>
        ) : (
          <p className="text-gray-500">No statistics available.</p>
        )}
      </section>
    </div>
  );
}
