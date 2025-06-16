// src/layouts/DashboardLayout.tsx
import { useState } from "react";

import { Outlet } from "react-router-dom";
import SidebarAdmin from "./Sidebaradmin";
import Topbaradmin from "./Topbaradmin";

export default function DashboardLayout() {
  const [collapsed, setCollapsed] = useState(false);

  return (
    <div className="flex h-screen">
      <SidebarAdmin collapsed={collapsed} setCollapsed={setCollapsed} />
<div className="flex flex-col flex-1 h-screen relative">
        <Topbaradmin collapsed={collapsed} setCollapsed={setCollapsed} />
        <main
          className={`flex-1 mt-20 p-6 overflow-y-auto ${
            collapsed ? "ml-16" : "ml-64"
          }`}
        >
          <Outlet />
        </main>
      </div>
    </div>
  );
}
