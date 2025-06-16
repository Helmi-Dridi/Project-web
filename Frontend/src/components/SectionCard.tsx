import React from "react";
import { Pencil, X, GripVertical } from "lucide-react";
import { motion } from "framer-motion";

type SectionCardProps = {
  title: string;
  isEditing: boolean;
  onEdit: () => void;
  onCancel: () => void;
  children: React.ReactNode;
  dragHandleProps?: {
    listeners?: React.HTMLAttributes<HTMLElement>;
    attributes?: React.HTMLAttributes<HTMLElement>;
  };
};

export const SectionCard = ({
  title,
  isEditing,
  onEdit,
  onCancel,
  children,
  dragHandleProps,
}: SectionCardProps) => {
  return (
    <motion.section
      initial={{ opacity: 0, y: 16 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3 }}
      className="relative bg-gray-800 border border-gray-700 rounded-2xl shadow-lg p-6 w-full space-y-6 transition-all"
    >
      {/* Drag Handle (optional) */}
      {dragHandleProps && (
        <div
          {...dragHandleProps.attributes}
          {...dragHandleProps.listeners}
          className="absolute top-4 right-4 cursor-grab text-gray-400 hover:text-gray-200"
          title="Drag section"
        >
          <GripVertical size={18} />
        </div>
      )}

      <div className="flex justify-between items-center border-b border-gray-700 pb-4">
        <h2 className="text-2xl font-semibold text-white">{title}</h2>
        <button
          onClick={isEditing ? onCancel : onEdit}
          className={`flex items-center gap-2 text-sm font-medium px-4 py-2 rounded-lg transition-colors duration-200 shadow-sm ${
            isEditing
              ? "bg-gray-700 text-gray-100 hover:bg-gray-600"
              : "bg-blue-600 text-white hover:bg-blue-700"
          }`}
        >
          {isEditing ? (
            <>
              <X size={16} />
              Cancel
            </>
          ) : (
            <>
              <Pencil size={16} />
              Edit
            </>
          )}
        </button>
      </div>

      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.1 }}
        className="space-y-4"
      >
        {children}
      </motion.div>
    </motion.section>
  );
};
