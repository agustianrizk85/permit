import { useEffect, useMemo, useState } from "react";
import { Link, useParams } from "react-router-dom";
import type { ProcessStep, Project } from "@/models";
import { projectService } from "@/services/project.service";
import { StepCard } from "@/components/StepCard";
import { categoryLabels } from "@/lib/processCatalog";

export function ProjectDetailPage() {
  const { id } = useParams<{ id: string }>();
  const projectId = Number(id);
  const [project, setProject] = useState<Project | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    setLoading(true);
    projectService
      .get(projectId)
      .then(setProject)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, [projectId]);

  const steps = project?.steps ?? [];

  const progress = useMemo(() => {
    if (steps.length === 0) return 0;
    const done = steps.filter((s) => s.status === "done").length;
    return Math.round((done / steps.length) * 100);
  }, [steps]);

  // Group steps by macro category (only A in milestone 1).
  const grouped = useMemo(() => {
    const map = new Map<string, ProcessStep[]>();
    for (const s of steps) {
      const list = map.get(s.category) ?? [];
      list.push(s);
      map.set(s.category, list);
    }
    return Array.from(map.entries());
  }, [steps]);

  const onStepChange = (updated: ProcessStep) => {
    setProject((prev) =>
      prev
        ? { ...prev, steps: prev.steps?.map((s) => (s.id === updated.id ? updated : s)) }
        : prev,
    );
  };

  if (loading) return <div className="muted">Memuat…</div>;
  if (error) return <div className="alert alert-error">{error}</div>;
  if (!project) return <div className="muted">Proyek tidak ditemukan.</div>;

  return (
    <div className="page">
      <Link to="/" className="back-link">
        ← Kembali ke Dashboard
      </Link>

      <div className="page-head">
        <div>
          <h1>{project.name}</h1>
          <p className="muted">
            {project.location || "—"} · Pemilik: {project.owner_name || "—"} · PT:{" "}
            {project.pt_name || "—"}
          </p>
        </div>
      </div>

      <div className="progress-bar-wrap">
        <div className="progress-bar">
          <div className="progress-fill" style={{ width: `${progress}%` }} />
        </div>
        <span className="progress-label">{progress}% selesai</span>
      </div>

      {grouped.map(([category, list]) => (
        <section key={category} className="category-section">
          <h2 className="category-title">{categoryLabels[category] ?? `Kategori ${category}`}</h2>
          <div className="steps">
            {list.map((s) => (
              <StepCard key={s.id} step={s} onChange={onStepChange} />
            ))}
          </div>
        </section>
      ))}
    </div>
  );
}
