// Domain models — typed mirror of the Go backend (internal/model).

export type Role = "ceo" | "dirops" | "kadep" | "legal_permit";

export interface User {
  id: number;
  name: string;
  email: string;
  role: Role;
  created_at: string;
  updated_at: string;
}

export type ProjectStage = "pra_akad" | "akad" | "permit" | "legal" | "done";

export interface Project {
  id: number;
  name: string;
  location: string;
  owner_name: string;
  pt_name: string;
  stage: ProjectStage;
  created_by: number;
  created_at: string;
  updated_at: string;
  steps?: ProcessStep[];
}

export type StepStatus = "pending" | "in_progress" | "done";

export interface ProcessStep {
  id: number;
  project_id: number;
  code: string;
  category: string;
  name: string;
  sequence: number;
  status: StepStatus;
  requires_price: boolean;
  requires_spk: boolean;
  price_label: string;
  notify_departments: boolean;
  confidential_output: boolean;
  sla_days: number;
  due_date: string | null;
  price_fix: number;
  spk_number: string;
  notes: string;
  metadata: Record<string, unknown> | null;
  completed_by: number | null;
  completed_at: string | null;
  documents?: DocumentFile[];
}

export interface DocumentFile {
  id: number;
  project_id: number;
  process_step_id: number | null;
  doc_type: string;
  original_name: string;
  stored_name: string;
  mime_type: string;
  size_bytes: number;
  confidential: boolean;
  ocr_data: Record<string, unknown> | null;
  uploaded_by: number;
  created_at: string;
}

export interface ProjectProgress {
  project_id: number;
  total: number;
  done: number;
  percentage: number;
  by_status: Partial<Record<StepStatus, number>>;
}

export interface LoginResponse {
  token: string;
  expires_at: string;
  user: User;
}

// Request payloads.
export interface CreateProjectInput {
  name: string;
  location?: string;
  owner_name?: string;
  pt_name?: string;
}

export interface UpdateStepInput {
  status?: StepStatus;
  price_fix?: number;
  spk_number?: string;
  notes?: string;
  metadata?: Record<string, unknown>;
  sla_days?: number;
  due_date?: string;
}

// --- Master Data PT (Proses E) ---

export interface PTDocument {
  id: number;
  pt_master_id: number;
  doc_type: string;
  original_name: string;
  stored_name: string;
  mime_type: string;
  size_bytes: number;
  uploaded_by: number;
  created_at: string;
}

export interface PTMaster {
  id: number;
  name: string;
  npwp: string;
  nib: string;
  notes: string;
  created_by: number;
  created_at: string;
  updated_at: string;
  documents?: PTDocument[];
}

export interface CreatePTInput {
  name: string;
  npwp?: string;
  nib?: string;
  notes?: string;
}

// --- Dashboard ---

export type WarningSeverity = "critical" | "warning" | "info";

export interface Warning {
  project_id: number;
  project_name: string;
  step_id: number;
  step_code: string;
  step_name: string;
  severity: WarningSeverity;
  message: string;
  due_date: string | null;
}

// --- Settings (DACI / Notification) ---

export interface DACIDriver {
  code: string;
  name: string;
}

export interface DACIConfig {
  drivers: DACIDriver[];
  approver: string[];
  consulting: string[];
  informed: string[];
}

export interface NotificationConfig {
  whatsapp_enabled: boolean;
  audit_ai_enabled: boolean;
  reminder_hour: number;
  whatsapp_api_url: string;
}

// --- Master Deadline ---

export interface DeadlineRule {
  code: string;
  name: string;
  category: string;
  deadline_enabled: boolean;
  alert_enabled: boolean;
  sla_days: number;
  updated_at?: string;
}

// --- OCR ---

export interface OCRResult {
  doc_type: string;
  fields: Record<string, string>;
  confidence: number;
  provider: string;
}
