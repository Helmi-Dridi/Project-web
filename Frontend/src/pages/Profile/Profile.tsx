import { useState } from "react";
import {
  DndContext,
  closestCenter,
  useSensor,
  useSensors,
  PointerSensor,
  type DragEndEvent,
} from "@dnd-kit/core";
import {
  arrayMove,
  SortableContext,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";

import {
  useStudentProfile,
  useStudentContact,
  useStudentAcademic,
  useStudentSettings,
} from "../../hooks/useProfile";

import { SectionCard } from "../../components/SectionCard";
import { SortableCard } from "../../components/SortableCard";
import UserProfileCard from "../../components/UserProfileCard";

import PersonalInfoForm from "../../components/PersonalInfoForm";
import PersonalInfoDisplay from "../../components/PersonalInfoDisplay";
import ContactDetailsForm from "../../components/ContactDetailsForm";
import ContactDetailsDisplay from "../../components/ContactDetailsDisplay";
import AcademicBackgroundForm from "../../components/AcademicBackgroundForm";
import AcademicBackgroundDisplay from "../../components/AcademicBackgroundDisplay";
import SettingsForm from "../../components/SettingsForm";
import SettingsDisplay from "../../components/SettingsDisplay";

import { uploadProfileImage } from "../../services/admin.service";
import { useAuth } from "../../context/AuthContext";
import type { StudentAcademicBackground } from "../../services/profile.service";

const initialSections = [
  { id: "personal", title: "Personal Information" },
  { id: "contact", title: "Contact Details" },
  { id: "academic", title: "Academic Background" },
  { id: "settings", title: "Preferences & Emergency" },
];

type SectionId = typeof initialSections[number]["id"];

export default function ProfilePage() {
  const [editingSection, setEditingSection] = useState<SectionId | null>(null);
  const [sections, setSections] = useState(initialSections);

  const { user, setUser } = useAuth();
  const sensors = useSensors(useSensor(PointerSensor));

  const profileQuery = useStudentProfile();
  const contactQuery = useStudentContact();
  const academicQuery = useStudentAcademic();
  const settingsQuery = useStudentSettings();

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;
    if (active.id !== over?.id) {
      const oldIndex = sections.findIndex((s) => s.id === active.id);
      const newIndex = sections.findIndex((s) => s.id === over?.id);
      setSections(arrayMove(sections, oldIndex, newIndex));
    }
  };

  const handleImageUpload = async (file: File): Promise<string> => {
    try {
      const relativePath = await uploadProfileImage(file);
      const fullUrl = `http://localhost:8080/${relativePath.replace(/\\/g, "/")}`;
      const updated = { ...user!, profilePicture: fullUrl };
      localStorage.setItem("user", JSON.stringify(updated));
      setUser(updated);
      return fullUrl;
    } catch (error) {
      console.error("Error uploading image:", error);
      alert("Failed to upload profile image.");
      return user?.profilePicture || "";
    }
  };

  const renderSection = (section: typeof initialSections[number]) => {
    const { id, title } = section;

    const sharedProps = {
      isEditing: editingSection === id,
      onEdit: () => setEditingSection(id),
      onCancel: () => setEditingSection(null),
      title,
    };

    return (
      <SortableCard key={id} id={id}>
        {(dragHandleProps) => (
          <SectionCard {...sharedProps} dragHandleProps={dragHandleProps}>
            {editingSection === id ? renderForm(id) : renderDisplay(id)}
          </SectionCard>
        )}
      </SortableCard>
    );
  };

  const renderForm = (id: SectionId) => {
    switch (id) {
      case "personal":
        return profileQuery.data ? (
          <PersonalInfoForm
            initialValues={profileQuery.data.data}
            onSuccess={() => setEditingSection(null)}
          />
        ) : (
          <p className="text-gray-500 italic">Loading...</p>
        );
      case "contact":
        return contactQuery.data ? (
          <ContactDetailsForm
            initialValues={contactQuery.data.data}
            onSuccess={() => setEditingSection(null)}
          />
        ) : (
          <p className="text-gray-500 italic">Loading...</p>
        );
      case "academic":
        return academicQuery.data ? (
          <AcademicBackgroundForm
            initialValues={academicQuery.data?.data ?? {} as StudentAcademicBackground}
            onSuccess={() => setEditingSection(null)}
          />
        ) : (
          <p className="text-gray-500 italic">Loading...</p>
        );
      case "settings":
        return settingsQuery.data ? (
          <SettingsForm
            initialValues={settingsQuery.data.data}
            onSuccess={() => setEditingSection(null)}
          />
        ) : (
          <p className="text-gray-500 italic">Loading...</p>
        );
    }
  };

  const renderDisplay = (id: SectionId) => {
    switch (id) {
      case "personal":
        return <PersonalInfoDisplay data={profileQuery.data?.data} />;
      case "contact":
        return <ContactDetailsDisplay data={contactQuery.data?.data} />;
      case "academic":
        return <AcademicBackgroundDisplay data={academicQuery.data?.data} />;
      case "settings":
        return <SettingsDisplay data={settingsQuery.data?.data} />;
    }
  };

 return (
  <div className="max-w-6xl mx-auto px-4 py-10">
    <div className="bg-gray-900 border border-gray-700 rounded-3xl p-8 shadow-lg space-y-10">
      <UserProfileCard onImageUpload={handleImageUpload} />

      <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
        <SortableContext items={sections.map((s) => s.id)} strategy={verticalListSortingStrategy}>
          <div className="space-y-6">{sections.map((section) => renderSection(section))}</div>
        </SortableContext>
      </DndContext>
    </div>
  </div>
);

}
