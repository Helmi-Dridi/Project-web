import React from "react";
import { motion } from "framer-motion";

type SectionCardProps = {
  title: string;
  children: React.ReactNode;
};

export const SectionCardBack = ({ title, children }: SectionCardProps) => {
  return (
    <motion.section
      initial={{ opacity: 0, y: 16 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3 }}
      className="bg-gray-800 border border-gray-700 rounded-2xl shadow-lg p-6 w-full space-y-6 transition-all"
    >
      <div className="flex justify-between items-center border-b border-gray-700 pb-4">
        <h2 className="text-2xl font-semibold text-white">{title}</h2>
      </div>

      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.1 }}
        className="space-y-4 text-gray-200"
      >
        {children}
      </motion.div>
    </motion.section>
  );
};
