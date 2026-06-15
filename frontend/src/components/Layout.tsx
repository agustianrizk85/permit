import { useEffect, useState } from "react";
import { NavLink, Outlet } from "react-router-dom";
import { useAuth } from "@/context/AuthContext";

const roleLabel: Record<string, string> = {
  ceo: "CEO",
  dirops: "Direktur Operasional",
  kadep: "Kepala Departemen",
  legal_permit: "Legal Permit",
};

const NAV = [
  { to: "/", label: "Dashboard", end: true },
  { to: "/pt", label: "Master PT" },
  { to: "/vendors", label: "Vendor" },
  { to: "/spk", label: "SPK" },
  { to: "/deadline", label: "Deadline" },
  { to: "/settings", label: "Setting" },
];

function Clock() {
  const [now, setNow] = useState(() => new Date());
  useEffect(() => {
    const i = setInterval(() => setNow(new Date()), 1000);
    return () => clearInterval(i);
  }, []);
  return (
    <div className="clock">
      <div className="t">{now.toLocaleTimeString("id-ID", { hour: "2-digit", minute: "2-digit", second: "2-digit" })}</div>
      <div className="d">
        {now.toLocaleDateString("id-ID", { weekday: "short", day: "numeric", month: "short", year: "numeric" })}
      </div>
    </div>
  );
}

export function Layout() {
  const { user, logout } = useAuth();
  return (
    <div className="app">
      <header className="hdr">
        <div className="hdr-logo">
          <span>GP</span>
        </div>
        <div className="hdr-titles">
          <h1>Legal Permit System</h1>
          <div className="sub">Greenpark Group · Departemen Legal &amp; Perizinan</div>
          <div className="tag">PRA-AKAD · AKAD · PERMIT · LEGAL</div>
        </div>
        <div className="hdr-spacer" />
        <div className="hdr-meta">
          <Clock />
          <div className="hdr-user">
            <div className="hu-name">{user?.name}</div>
            <div className="hu-role">{user ? roleLabel[user.role] ?? user.role : ""}</div>
          </div>
          <button className="logout-btn" onClick={logout} title="Keluar">
            ✕
          </button>
        </div>
      </header>

      <nav className="topnav">
        {NAV.map((n) => (
          <NavLink
            key={n.to}
            to={n.to}
            end={n.end}
            className={({ isActive }) => (isActive ? "navlink active" : "navlink")}
          >
            {n.label}
          </NavLink>
        ))}
      </nav>

      <main className="content">
        <Outlet />
      </main>
    </div>
  );
}
