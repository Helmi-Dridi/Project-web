// =======================
// Imports
// =======================
import { useAuth } from "../context/AuthContext";
import { useStudentEnrollment } from "../hooks/useEnrollment";
import {
  useStudentProfile,
  useStudentContact,
  useStudentAcademic,
  useStudentSettings,
} from "../hooks/useProfile";
import { useDocuments } from "../hooks/hookdocument";
import { useMyAppointments } from "../hooks/useBooking";
import { FaCheckCircle, FaCalendarAlt } from "react-icons/fa";
import { motion } from "framer-motion";

// =======================
// Component
// =======================
const Dashboard = () => {
  const { user } = useAuth();

  // 🔐 Block rendering until user.ID is ready (for dependent hooks like useMyAppointments)
  if (!user?.ID) {
    return (
      <div className="text-center py-12 text-gray-500 font-medium text-sm">
        Loading your dashboard...
      </div>
    );
  }

  // =======================
  // Data Hooks
  // =======================
  const { data: enrollment, isLoading, isError } = useStudentEnrollment();
  const profile = useStudentProfile();
  const contact = useStudentContact();
  const academic = useStudentAcademic();
  const settings = useStudentSettings();
  const documents = useDocuments(1, 20);
  const { data: appointmentsData, isLoading: loadingAppointments } = useMyAppointments();

  // =======================
  // Stepper Logic
  // =======================
  const isStepComplete = (entry: any) => entry?.data?.data?.step === true;

  const allDocsApproved =
    documents.data?.items?.length &&
    documents.data.items.every((doc) => doc.verification === "approved");

  const steps = [
    {
      id: "personal",
      label: "Personal Info",
      description: "Basic personal details like name, birth date, and nationality.",
      completed: isStepComplete(profile),
    },
    {
      id: "contact",
      label: "Contact Info",
      description: "Your phone, email, and physical address.",
      completed: isStepComplete(contact),
    },
    {
      id: "academic",
      label: "Academic Background",
      description: "GPA, institution, and qualifications.",
      completed: isStepComplete(academic),
    },
    {
      id: "settings",
      label: "Preferences & Emergency",
      description: "Language preferences and emergency contact info.",
      completed: isStepComplete(settings),
    },
    {
      id: "documents",
      label: "Document Uploads",
      description: "Upload your official documents like passport, CV, and certificates.",
      completed: allDocsApproved,
    },
    {
      id: "enrollment",
      label: "Enrollment",
      description: "Enrollment confirmation in your program and university.",
      completed: !!enrollment,
    },
    {
      id: "visa",
      label: "Visa Booking",
      description: "Book your visa appointment and upload your confirmation.",
      completed: false,
    },
  ];

  const currentStep = steps.findIndex((s) => !s.completed);
  const allAppointments = appointmentsData?.items || [];

  // =======================
  // Render
  // =======================
  return (
    <div className="max-w-7xl mx-auto px-6 py-10 space-y-12">

      {/* ======================= */}
      {/* Header */}
      {/* ======================= */}
      <div className="relative bg-gradient-to-r from-blue-600 to-blue-500 text-white px-8 py-10 rounded-3xl shadow-xl overflow-hidden">
        <div className="flex flex-col sm:flex-row sm:items-center gap-6">
          <img
            src={user?.profilePicture || "/default-avatar.png"}
            alt="Profile"
            className="w-20 h-20 sm:w-24 sm:h-24 rounded-full object-cover border-4 border-white shadow-lg"
          />
          <div>
            <h1 className="text-3xl sm:text-4xl font-bold mb-1">
              Welcome, {user?.name?.split(" ")[0]}! 👋
            </h1>
            <p className="text-sm sm:text-base opacity-90 mb-1">{user?.email}</p>
            <p className="text-sm sm:text-base opacity-80 mt-1 max-w-xl">
              We're glad to have you at <span className="font-semibold">ScholarRev</span> — your personalized platform to guide you from admission to visa booking. Let's achieve your academic goals together!
            </p>
          </div>
        </div>
        <div className="absolute bottom-0 right-0 w-48 h-48 bg-white/10 rounded-full blur-2xl opacity-20 translate-x-12 translate-y-12 pointer-events-none" />
      </div>

      {/* ======================= */}
      {/* Stepper */}
      {/* ======================= */}
      <div className="overflow-x-auto scrollbar-thin scrollbar-thumb-blue-300 scrollbar-track-transparent">
        <div className="min-w-[1000px] flex items-center justify-start gap-12 py-6 px-4 relative">
          {steps.map((step, index) => {
            const isCurrent = index === currentStep;
            return (
              <div key={step.id} className="relative flex flex-col items-center min-w-[120px]">
                {index !== 0 && (
                  <div className="absolute top-6 -left-[60px] w-[60px] h-1 bg-gray-300 z-0" />
                )}
                <motion.div
                  initial={{ scale: 0.9, opacity: 0.6 }}
                  animate={{ scale: 1, opacity: 1 }}
                  transition={{ duration: 0.3, delay: index * 0.05 }}
                  className={`w-14 h-14 rounded-full border-4 flex items-center justify-center text-lg font-semibold shadow-lg z-10
                    ${step.completed
                      ? "bg-green-500 border-green-600 text-white"
                      : isCurrent
                      ? "bg-blue-500 border-blue-600 text-white animate-pulse"
                      : "bg-white border-gray-300 text-gray-400"}
                  `}
                >
                  {step.completed ? <FaCheckCircle className="w-6 h-6" /> : index + 1}
                </motion.div>
                <div className="text-sm text-center mt-3 font-medium text-gray-700 w-32">{step.label}</div>
                <div className="mt-2 text-xs text-center text-gray-500 w-36 leading-snug">{step.description}</div>
              </div>
            );
          })}
        </div>
      </div>

      {/* ======================= */}
      {/* Enrollment Section */}
      {/* ======================= */}
      <div className="bg-white rounded-3xl p-6 shadow-lg border border-gray-200 space-y-4">
        <h2 className="text-xl font-bold text-gray-800">🎓 Enrollment Status</h2>
        {isLoading ? (
          <div className="bg-gray-100 p-4 rounded text-center text-sm text-gray-600">
            Checking enrollment status...
          </div>
        ) : isError || !enrollment ? (
          <div className="bg-yellow-100 border-l-4 border-yellow-500 text-yellow-800 p-4 rounded shadow-md">
            <h3 className="text-lg font-semibold mb-1">⏳ Under Review</h3>
            <p>Your enrollment is being reviewed by your advisor. We'll notify you once you're officially accepted. Stay tuned!</p>
          </div>
        ) : (
          <div className="bg-green-100 border-l-4 border-green-500 text-green-800 p-4 rounded shadow-md">
            <h3 className="text-lg font-semibold mb-1">✅ Congratulations!</h3>
            <p>
              You've been successfully enrolled in <strong>{enrollment.program.programName}</strong> at{" "}
              <strong>{enrollment.university.name}</strong>. We’re thrilled to have you onboard! 🎉
            </p>
            <p className="mt-2 text-sm text-green-700">Keep an eye on your email for orientation details and next steps.</p>
          </div>
        )}
      </div>

      {/* ======================= */}
      {/* Appointments Section */}
      {/* ======================= */}
      <div className="bg-white rounded-3xl p-6 shadow-lg border border-gray-200 space-y-4">
        <h2 className="text-xl font-bold text-gray-800">📅 Confirmed Appointments</h2>
        {loadingAppointments ? (
          <div className="text-gray-500 text-sm">Loading your appointments...</div>
        ) : allAppointments.length ? (
          <ul className="divide-y divide-gray-100">
            {allAppointments.map((appt) => (
              <li key={appt.id} className="flex items-center justify-between py-3">
                <div className="flex items-center gap-3">
                  <FaCalendarAlt className="text-blue-500 w-5 h-5" />
                  <span className="text-gray-700 font-medium">
                    {appt.date} at {appt.timeSlot}
                  </span>
                </div>
                <span
                  className={`text-sm font-semibold ${
                    appt.status === "confirmed"
                      ? "text-green-600"
                      : appt.status === "booked"
                      ? "text-blue-500"
                      : appt.status === "canceled"
                      ? "text-red-500"
                      : "text-gray-500"
                  }`}
                >
                  {appt.status.charAt(0).toUpperCase() + appt.status.slice(1)}
                </span>
              </li>
            ))}
          </ul>
        ) : (
          <div className="text-sm text-gray-600 bg-gray-50 p-4 rounded border border-dashed">
            You haven't booked any appointments yet. Once you do, they will appear here.
          </div>
        )}
      </div>
    </div>
  );
};

export default Dashboard;
