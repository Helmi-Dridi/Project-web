import { useAllStudentEnrollments } from "../hooks/useEnrollment";

export default function AllEnrollmentsTable() {
  const { data, isLoading, isError } = useAllStudentEnrollments();

  if (isLoading) return <p className="text-gray-500">Loading enrollments...</p>;
  if (isError) return <p className="text-red-500">Failed to load enrollments.</p>;

  return (
    <div className="overflow-x-auto rounded shadow border border-gray-200">
      <table className="min-w-full bg-white text-sm text-left text-gray-700">
        <thead className="bg-gray-100 font-medium">
          <tr>
            <th className="py-2 px-4">Student </th>
            <th className="py-2 px-4">Email </th>
            <th className="py-2 px-4">Program</th>
            <th className="py-2 px-4">University</th>
            <th className="py-2 px-4">Enrolled At</th>
          </tr>
        </thead>
        <tbody>
          {data?.map((e) => (
            <tr key={e.userId + e.programId} className="hover:bg-gray-50">
              <td className="py-2 px-4">{e.student.firstName} {e.student.lastName}</td>
              <td className="py-2 px-4">{e.student.email}</td>
              <td className="py-2 px-4">{e.program.programName}</td>
              <td className="py-2 px-4">{e.university.name}</td>
              <td className="py-2 px-4">{new Date(e.enrolledAt).toLocaleDateString()}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
