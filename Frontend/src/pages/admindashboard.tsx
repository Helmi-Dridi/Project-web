import { useAuth } from "../context/AuthContext";

export default function AdminDashboard() {
  const { logout } = useAuth();

  const handleLogout = () => {
   logout();
  };

  return (
    <div className="p-6 space-y-8">
      <section>
        <h1 className="text-3xl font-bold text-gray-800">Admin Dashboard</h1>
        <p className="text-gray-600 mt-1">Overview of user progress and upcoming actions.</p>
      </section>

      <section>
        <h2 className="text-xl font-semibold text-gray-800">📊 Application Status</h2>
        <ul className="list-disc list-inside text-gray-700 mt-2">
          <li>65% of applications completed</li>
          <li>Most common incomplete step: Visa Appointment</li>
          <li>Average completion time: 3.2 weeks</li>
        </ul>
      </section>

      <section>
        <h2 className="text-xl font-semibold text-gray-800">📌 Next Student Actions</h2>
        <ol className="list-decimal list-inside text-gray-700 mt-2">
          <li>Upload IELTS Certificate</li>
          <li>Sign Study Contract</li>
          <li>Book Visa Appointment</li>
        </ol>
      </section>

      <section>
        <h2 className="text-xl font-semibold text-gray-800">📅 Upcoming Events</h2>
        <ul className="list-disc list-inside text-gray-700 mt-2">
          <li>Visa Interviews: June 10, 9:00 AM</li>
          <li>Bank statement deadline: June 5</li>
        </ul>
      </section>

      <section>
        <h2 className="text-xl font-semibold text-gray-800">🔔 Latest Activities</h2>
        <ul className="list-disc list-inside text-gray-700 mt-2">
          <li>Sarah Adams submitted an application to Concordia</li>
          <li>Visa approval received for 3 students</li>
        </ul>
      </section>

      <section>
        <h2 className="text-xl font-semibold text-gray-800">📣 Notes</h2>
        <p className="text-gray-700 mt-1">Referral program ongoing — remind users to share their referral link for bonus rewards.</p>
        <p className="text-gray-700 mt-1">22 days remaining until student batch departure.</p>
      </section>

      <div className="pt-4 text-right">
        <button
          onClick={handleLogout}
          className="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600 transition"
        >
          Logout
        </button>
      </div>
    </div>
  );
}
