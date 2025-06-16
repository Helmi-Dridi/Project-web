import { useState } from "react";
import { GraduationCap, Users, UserPlus } from "lucide-react";
import { useUniversityPrograms } from "../hooks/useUniversity";
import { useStudents } from "../hooks/useAdmin";
import {
  useAssignStudent,
  useEnrollmentsByProgram,
  useUnassignStudent,
} from "../hooks/useEnrollment";
import Modal from "./Modal";

interface TabbedModalProps {
  universityId: string;
  onClose: () => void;
}

const TabbedUniversityModal = ({ universityId, onClose }: TabbedModalProps) => {
  const [activeTab, setActiveTab] = useState("programs");
  const [selectedProgramId, setSelectedProgramId] = useState<string | null>(null);
  const [selectedStudentId, setSelectedStudentId] = useState<string | null>(null);
  const [programTypeFilter, setProgramTypeFilter] = useState("");

  const { data: programs = [] } = useUniversityPrograms(universityId);
  const { data: students = [] } = useStudents();
  const { data: enrollments = [] } = useEnrollmentsByProgram(selectedProgramId || "");
  const assignMutation = useAssignStudent();
  const unassignMutation =  useUnassignStudent(selectedProgramId || "");

  const handleAssign = () => {
    if (selectedStudentId && selectedProgramId) {
      assignMutation.mutate({ userId: selectedStudentId, programId: selectedProgramId });
    }
  };

  const handleUnassign = (userId: string) => {
    unassignMutation.mutate(userId);
  };

  const filteredPrograms = programTypeFilter
    ? programs.filter((p) => p.programType === programTypeFilter)
    : programs;

  return (
    <Modal isOpen={true} onClose={onClose}>
      <div className="space-y-6">
        <div className="flex gap-4 border-b border-gray-700 pb-2">
          <button
            onClick={() => setActiveTab("programs")}
            className={`px-4 py-2 rounded ${
              activeTab === "programs" ? "bg-gray-700 text-white" : "text-gray-400"
            }`}
          >
            <GraduationCap className="inline mr-2" size={16} /> Programs
          </button>
          <button
            onClick={() => setActiveTab("assign")}
            className={`px-4 py-2 rounded ${
              activeTab === "assign" ? "bg-gray-700 text-white" : "text-gray-400"
            }`}
          >
            <UserPlus className="inline mr-2" size={16} /> Assign
          </button>
          <button
            onClick={() => setActiveTab("enrolled")}
            className={`px-4 py-2 rounded ${
              activeTab === "enrolled" ? "bg-gray-700 text-white" : "text-gray-400"
            }`}
          >
            <Users className="inline mr-2" size={16} /> Enrolled
          </button>
        </div>

        {activeTab === "programs" && (
          <div className="space-y-4">
            <div className="flex justify-end">
              <select
                value={programTypeFilter}
                onChange={(e) => setProgramTypeFilter(e.target.value)}
                className="px-3 py-2 rounded bg-gray-800 text-white border border-gray-600 text-sm"
              >
                <option value="">All Types</option>
                {[...new Set(programs.map((p) => p.programType))].map((type) => (
                  <option key={type} value={type}>
                    {type}
                  </option>
                ))}
              </select>
            </div>

            <ul className="space-y-4">
              {filteredPrograms.map((p) => (
                <li
                  key={p.id}
                  className={`border p-4 rounded-xl cursor-pointer ${
                    selectedProgramId === p.id ? "border-yellow-500" : "border-gray-700"
                  }`}
                  onClick={() => setSelectedProgramId(p.id)}
                >
                  <p>
                    <strong>Name:</strong> {p.programName}
                  </p>
                  <p>
                    <strong>Type:</strong> {p.programType}
                  </p>
                  <p>
                    <strong>Language:</strong> {p.languageOfInstruction}
                  </p>
                </li>
              ))}
            </ul>
          </div>
        )}

        {activeTab === "assign" && (
          <div className="space-y-4">
            <select
              onChange={(e) => setSelectedStudentId(e.target.value)}
              value={selectedStudentId || ""}
              className="w-full px-4 py-2 rounded bg-gray-800 border border-gray-600 text-white"
            >
              <option value="">Select Student</option>
              {students.map((s) => (
                <option key={s.id} value={s.id}>
                  {s.firstName} || {s.lastName}
                </option>
              ))}
            </select>
            <button
              onClick={handleAssign}
              className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded"
              disabled={!selectedStudentId || !selectedProgramId}
            >
              Assign Student to Program
            </button>
          </div>
        )}

        {activeTab === "enrolled" && (
  <ul className="space-y-3">
    {enrollments?.map((en) => (
      <li
        key={en.userId}
        className="bg-gray-800 p-4 rounded text-white flex justify-between items-center"
      >
        <div>
          <p>
            <strong>Name:</strong> {en.student.firstName} {en.student.lastName}
          </p>
          <p>
            <strong>Email:</strong> {en.student.email}
          </p>
          <p>
            <strong>Program:</strong> {en.program.programName}
          </p>
        </div>
        <button
          onClick={() => handleUnassign(en.userId)}
          className="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded"
        >
          Unassign
        </button>
      </li>
    ))}
  </ul>
)}

      </div>
    </Modal>
  );
};

export default TabbedUniversityModal;
