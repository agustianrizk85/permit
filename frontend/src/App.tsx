import { Navigate, Route, Routes } from "react-router-dom";
import { useAuth } from "@/context/AuthContext";
import { LoginPage } from "@/pages/LoginPage";
import { DashboardPage } from "@/pages/DashboardPage";
import { ProjectDetailPage } from "@/pages/ProjectDetailPage";
import { PTPage } from "@/pages/PTPage";
import { VendorPage } from "@/pages/VendorPage";
import { SPKPage } from "@/pages/SPKPage";
import { DeadlinePage } from "@/pages/DeadlinePage";
import { SettingsPage } from "@/pages/SettingsPage";
import { Layout } from "@/components/Layout";
import type { ReactNode } from "react";

function RequireAuth({ children }: { children: ReactNode }) {
  const { user, loading } = useAuth();
  if (loading) return <div className="center muted">Memuat…</div>;
  if (!user) return <Navigate to="/login" replace />;
  return <>{children}</>;
}

export function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route
        element={
          <RequireAuth>
            <Layout />
          </RequireAuth>
        }
      >
        <Route path="/" element={<DashboardPage />} />
        <Route path="/projects/:id" element={<ProjectDetailPage />} />
        <Route path="/pt" element={<PTPage />} />
        <Route path="/vendors" element={<VendorPage />} />
        <Route path="/spk" element={<SPKPage />} />
        <Route path="/deadline" element={<DeadlinePage />} />
        <Route path="/settings" element={<SettingsPage />} />
      </Route>
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}
