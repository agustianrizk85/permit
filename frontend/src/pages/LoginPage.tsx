import { useState, type FormEvent } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "@/context/AuthContext";

export function LoginPage() {
  const { login, user } = useAuth();
  const navigate = useNavigate();
  const [email, setEmail] = useState("dirops@greenpark.id");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  if (user) navigate("/", { replace: true });

  const onSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      await login(email, password);
      navigate("/", { replace: true });
    } catch (err) {
      setError(err instanceof Error ? err.message : "Gagal masuk");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-wrap">
      <form className="login-card" onSubmit={onSubmit}>
        <div className="login-brand">🌿 Green Park</div>
        <h1>Legal Permit System</h1>
        <p className="muted">Masuk untuk melanjutkan</p>

        <label className="field">
          <span>Email</span>
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            autoFocus
          />
        </label>

        <label className="field">
          <span>Kata Sandi</span>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </label>

        {error && <div className="alert alert-error">{error}</div>}

        <button className="btn btn-primary btn-block" type="submit" disabled={loading}>
          {loading ? "Memproses…" : "Masuk"}
        </button>

        <div className="login-hint muted">
          Akun: <code>dirops@greenpark.id</code> / <code>kadep@greenpark.id</code>
        </div>
      </form>
    </div>
  );
}
