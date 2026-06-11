import { Link, NavLink, Outlet } from "react-router-dom";
import { useAuth } from "@/context/AuthContext";

const roleLabel: Record<string, string> = {
  ceo: "CEO",
  dirops: "Direktur Operasional",
  kadep: "Kepala Departemen",
  legal_permit: "Legal Permit",
};

export function Layout() {
  const { user, logout } = useAuth();
  return (
    <div className="app">
      <header className="topbar">
        <Link to="/" className="brand">
          🌿 Green Park <span className="brand-sub">Legal Permit</span>
        </Link>
        <nav className="topnav">
          <NavLink to="/" end className={({ isActive }) => (isActive ? "navlink active" : "navlink")}>
            Dashboard
          </NavLink>
          <NavLink to="/pt" className={({ isActive }) => (isActive ? "navlink active" : "navlink")}>
            Master PT
          </NavLink>
          <NavLink to="/deadline" className={({ isActive }) => (isActive ? "navlink active" : "navlink")}>
            Deadline
          </NavLink>
          <NavLink to="/settings" className={({ isActive }) => (isActive ? "navlink active" : "navlink")}>
            Setting
          </NavLink>
        </nav>
        <div className="topbar-right">
          <div className="user-chip">
            <span className="user-name">{user?.name}</span>
            <span className="user-role">{user ? roleLabel[user.role] : ""}</span>
          </div>
          <button className="btn btn-ghost" onClick={logout}>
            Keluar
          </button>
        </div>
      </header>
      <main className="content">
        <Outlet />
      </main>
    </div>
  );
}
