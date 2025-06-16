import { useState } from "react";
import { useUniversities } from "../hooks/useUniversity";
import UniversityProgramModal from "../components/UniversityProgramModal";
import AllEnrollmentsTable from "../components/AllEnrollmentsTable"; // You said this is created
import { GraduationCap } from "lucide-react";

export default function UniversityListPage() {
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const [typeFilter, setTypeFilter] = useState("");
  const [selectedUniversityId, setSelectedUniversityId] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<"universities" | "enrollments">("universities");

  const limit = 10;
  const { data, isLoading, isError } = useUniversities(page, limit);

  const totalPages = data ? Math.ceil(data.totalCount / data.limit) : 1;
  const universityTypes = Array.from(new Set(data?.items.map((u) => u.universityType)));

  const filteredItems =
    data?.items.filter((u) => {
      const matchSearch = u.name.toLowerCase().includes(search.toLowerCase());
      const matchType = typeFilter ? u.universityType === typeFilter : true;
      return matchSearch && matchType;
    }) ?? [];

  return (
    <div className="p-6 max-w-7xl mx-auto space-y-6">
      {/* Header & Tabs */}
      <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
        <h1 className="text-3xl font-bold text-gray-800 flex items-center gap-2">
          <GraduationCap className="text-blue-500" /> Dashboard
        </h1>
        <div className="flex gap-2 mt-2">
          <button
            onClick={() => setActiveTab("universities")}
            className={`px-4 py-2 rounded-md text-sm font-medium ${
              activeTab === "universities"
                ? "bg-blue-600 text-white"
                : "bg-gray-200 text-gray-700"
            }`}
          >
            Universities
          </button>
          <button
            onClick={() => setActiveTab("enrollments")}
            className={`px-4 py-2 rounded-md text-sm font-medium ${
              activeTab === "enrollments"
                ? "bg-blue-600 text-white"
                : "bg-gray-200 text-gray-700"
            }`}
          >
            Enrollments
          </button>
        </div>
      </div>

      {/* Conditional Tabs Content */}
      {activeTab === "universities" && (
        <>
          <div className="flex flex-col sm:flex-row gap-2 sm:items-center">
            <input
              type="text"
              placeholder="Search by name..."
              className="border border-gray-300 rounded-md px-3 py-2 text-sm"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
            />
            <select
              value={typeFilter}
              onChange={(e) => setTypeFilter(e.target.value)}
              className="border border-gray-300 rounded-md px-3 py-2 text-sm"
            >
              <option value="">All Types</option>
              {universityTypes.map((type) => (
                <option key={type} value={type}>
                  {type}
                </option>
              ))}
            </select>
          </div>

          {isLoading ? (
            <p className="text-gray-500">Loading universities...</p>
          ) : isError ? (
            <p className="text-red-500">Error loading universities</p>
          ) : (
            <>
              <div className="overflow-x-auto rounded shadow border border-gray-200 mt-4">
                <table className="min-w-full bg-white text-sm text-left text-gray-700">
                  <thead className="bg-gray-100 font-medium">
                    <tr>
                      <th className="py-2 px-4">Name</th>
                      <th className="py-2 px-4">Country</th>
                      <th className="py-2 px-4">City</th>
                      <th className="py-2 px-4">Type</th>
                      <th className="py-2 px-4">Actions</th>
                    </tr>
                  </thead>
                  <tbody>
                    {filteredItems.map((u) => (
                      <tr key={u.id} className="hover:bg-gray-50">
                        <td className="py-2 px-4">{u.name}</td>
                        <td className="py-2 px-4">{u.country}</td>
                        <td className="py-2 px-4">{u.city}</td>
                        <td className="py-2 px-4">{u.universityType}</td>
                        <td className="py-2 px-4">
                          <button
                            onClick={() => setSelectedUniversityId(u.id)}
                            className="text-blue-600 hover:underline"
                          >
                            View Programs
                          </button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              <div className="flex items-center justify-center gap-4 mt-4">
                <button
                  disabled={page === 1}
                  onClick={() => setPage((p) => p - 1)}
                  className="px-3 py-1 bg-gray-100 rounded hover:bg-gray-200 disabled:opacity-50"
                >
                  Prev
                </button>
                <span className="text-gray-700">
                  Page {page} of {totalPages}
                </span>
                <button
                  disabled={page >= totalPages}
                  onClick={() => setPage((p) => p + 1)}
                  className="px-3 py-1 bg-gray-100 rounded hover:bg-gray-200 disabled:opacity-50"
                >
                  Next
                </button>
              </div>
            </>
          )}
        </>
      )}

      {activeTab === "enrollments" && (
        <div className="mt-4">
          <AllEnrollmentsTable />
        </div>
      )}

      {selectedUniversityId && (
        <UniversityProgramModal
          universityId={selectedUniversityId}
          onClose={() => setSelectedUniversityId(null)}
        />
      )}
    </div>
  );
}
