package service

import "legalpermit/internal/model"

// StepTemplate is the static definition of a process step. Creating a project
// instantiates a model.ProcessStep for each template. The whole A..I flow is
// data-driven from the slices below — extending the system is just editing data.
type StepTemplate struct {
	Code               string
	Category           string
	Name               string
	RequiresPrice      bool
	RequiresSPK        bool
	PriceLabel         string
	NotifyDepartments  bool
	ConfidentialOutput bool
	SLADays            int      // default SLA; KADEP can override per step
	MetadataKeys       []string // structured fields the UI should collect
	DocTypes           []string // expected document slots for upload steps
}

// CategoryLabels maps a category code to its human label (also used by the UI).
var CategoryLabels = map[string]string{
	"A": "A · Pra-Akad Lahan (Simultan)",
	"B": "B · Akad Lahan (Kesinambungan)",
	"C": "C · Permit Pasca Akad Lahan",
	"D": "D · Legal Pasca Akad Lahan",
	"F": "F · Master Data PKS Bank",
	"H": "H · Flow Bisnis PKS Bank",
}

// ProcessA — Pra Akad Lahan (Simultan).
var ProcessA = []StepTemplate{
	{Code: "A1", Category: "A", Name: "Cek SHM", RequiresPrice: true, PriceLabel: "Harga Fix", SLADays: 7},
	{Code: "A2", Category: "A", Name: "Cek Zonasi", RequiresPrice: true, PriceLabel: "Harga Fix", SLADays: 7},
	{Code: "A3", Category: "A", Name: "Ukur BPN", RequiresPrice: true, RequiresSPK: true, PriceLabel: "Harga Fix", SLADays: 14},
	{Code: "A4", Category: "A", Name: "Cek Intip Ploting Lahan", RequiresPrice: true, PriceLabel: "Harga Fix", SLADays: 7},
	{Code: "A5", Category: "A", Name: "Kesepakatan Lingkungan", SLADays: 10,
		MetadataKeys: []string{"nama_rt_rw", "alamat_lahan", "tanggal"}},
	{Code: "A6", Category: "A", Name: "Data Identitas Pemilik Lahan", SLADays: 5,
		DocTypes: []string{"KTP", "KK", "NPWP", "Buku Nikah", "Lainnya"}},
	{Code: "A7", Category: "A", Name: "Data Legalitas Lahan", SLADays: 5,
		DocTypes: []string{"SHM", "PBB", "Lainnya"}},
	{Code: "A8", Category: "A", Name: "MOU Final Developer - Pemilik Lahan", RequiresPrice: true, PriceLabel: "UTJ", SLADays: 14},
}

// ProcessB — Akad Lahan (Kesinambungan).
var ProcessB = []StepTemplate{
	{Code: "B1", Category: "B", Name: "Draf Akta Notaris", SLADays: 7},
	{Code: "B2", Category: "B", Name: "Legal Meeting dengan Notaris", SLADays: 7,
		MetadataKeys: []string{"tanggal_meeting", "notaris", "catatan"}},
	{Code: "B3", Category: "B", Name: "TTD Akad Lahan di Notaris", SLADays: 14,
		MetadataKeys: []string{"biaya_akta_notaris", "biaya_pph", "biaya_bphtb", "biaya_ditanggung_para_pihak", "pt_yang_dipakai"}},
	{Code: "B4", Category: "B", Name: "Salinan Akta Notaris", SLADays: 14,
		DocTypes: []string{"Akta KSO", "Akta Kuasa Jual", "Akta Kuasa Mengelola Lahan", "Akta PPJB Termin"}},
}

// ProcessC — Permit Pasca Akad Lahan.
var ProcessC = []StepTemplate{
	{Code: "C1", Category: "C", Name: "Izin Lingkungan (Form Pemda)", RequiresPrice: true, PriceLabel: "Surat Kompensasi", SLADays: 21,
		DocTypes: []string{"KTP Warga & RT/RW", "Izin Tetangga", "Surat Kompensasi", "Kwitansi"}},
	{Code: "C2", Category: "C", Name: "Rekom Lurah/Camat (Simultan)", RequiresPrice: true, PriceLabel: "Harga Fix", SLADays: 14},
	{Code: "C3", Category: "C", Name: "ITR (Informasi Tata Ruang)", RequiresPrice: true, RequiresSPK: true, PriceLabel: "Harga Fix", SLADays: 21},
	{Code: "C4", Category: "C", Name: "Aspek TGT / Pertek BPN", RequiresPrice: true, RequiresSPK: true, PriceLabel: "Harga Fix", SLADays: 30},
	{Code: "C5", Category: "C", Name: "KKPR", RequiresPrice: true, RequiresSPK: true, PriceLabel: "Harga Fix", SLADays: 30},
	{Code: "C6", Category: "C", Name: "Pertek Limbah / Amdal", RequiresPrice: true, RequiresSPK: true, PriceLabel: "Harga Fix", SLADays: 45},
	{Code: "C7", Category: "C", Name: "Peil Banjir", RequiresPrice: true, RequiresSPK: true, PriceLabel: "Harga Fix", SLADays: 30},
	{Code: "C8", Category: "C", Name: "UKL-UPL", RequiresPrice: true, RequiresSPK: true, PriceLabel: "Harga Fix", SLADays: 45},
	{Code: "C9", Category: "C", Name: "Andalalin", RequiresPrice: true, RequiresSPK: true, PriceLabel: "Harga Fix", SLADays: 45},
	{Code: "C10", Category: "C", Name: "Siteplan", RequiresPrice: true, RequiresSPK: true, PriceLabel: "Harga Fix", NotifyDepartments: true, SLADays: 30},
	{Code: "C11", Category: "C", Name: "TPU-PSU", RequiresPrice: true, RequiresSPK: true, PriceLabel: "Harga Fix", SLADays: 30},
	{Code: "C12", Category: "C", Name: "IMB Induk-Rincik / PBG Induk-Rincik", RequiresPrice: true, RequiresSPK: true, PriceLabel: "Harga Fix", ConfidentialOutput: true, SLADays: 60},
	{Code: "C13", Category: "C", Name: "Izin Lainnya (Opsional)", SLADays: 30,
		MetadataKeys: []string{"nama_izin", "catatan"}},
}

// ProcessD — Legal Pasca Akad Lahan (berdasarkan Siteplan C10).
var ProcessD = []StepTemplate{
	{Code: "D1", Category: "D", Name: "Pemecahan SHM", RequiresPrice: true, RequiresSPK: true, PriceLabel: "Harga Fix", ConfidentialOutput: true, SLADays: 60},
	{Code: "D2", Category: "D", Name: "Pemecahan PBB", RequiresPrice: true, RequiresSPK: true, PriceLabel: "Harga Fix", ConfidentialOutput: true, SLADays: 45},
}

// ProcessF — Master Data PKS Bank (acuan data + formulir).
var ProcessF = []StepTemplate{
	{Code: "F1", Category: "F", Name: "Berkas Acuan PKS (Identitas, Legalitas, Akta, Siteplan, PT, Marketing)", SLADays: 7,
		MetadataKeys: []string{"gdrive_folder_link"}},
	{Code: "F2", Category: "F", Name: "Formulir PKS Bank (Upload)", SLADays: 7,
		DocTypes: []string{"Formulir PKS Bank"}},
}

// ProcessH — Flow Bisnis PKS Bank (Data Warehouse).
var ProcessH = []StepTemplate{
	{Code: "H1", Category: "H", Name: "Penyerahan Berkas ke Bank (Share Link Gdrive)", SLADays: 5,
		MetadataKeys: []string{"gdrive_link", "tanggal_serah", "nama_bank"}},
	{Code: "H2", Category: "H", Name: "Survei Bank ke Proyek/Lahan", SLADays: 14,
		MetadataKeys: []string{"catatan_survei", "tanggal_survei"}, DocTypes: []string{"Foto Survei"}},
	{Code: "H3", Category: "H", Name: "Draft Akta PKS Bank (Review Internal)", SLADays: 10,
		MetadataKeys: []string{"catatan_revisi"}},
	{Code: "H4", Category: "H", Name: "Tanda Tangan Para Pihak (PT & Bank)", SLADays: 14,
		MetadataKeys: []string{"tanggal_ttd"}},
	{Code: "H5", Category: "H", Name: "Salinan Akta PKS Bank", SLADays: 7,
		DocTypes: []string{"Salinan Akta PKS Bank"}},
}

// Catalog returns every step template instantiated for a new project, in order.
// SPK Legal Permit (Proses J) and Vendor master (Proses I) are standalone
// modules, not per-project steps, so they are not part of the catalog.
func Catalog() []StepTemplate {
	var all []StepTemplate
	all = append(all, ProcessA...)
	all = append(all, ProcessB...)
	all = append(all, ProcessC...)
	all = append(all, ProcessD...)
	all = append(all, ProcessF...)
	all = append(all, ProcessH...)
	return all
}

// toModel converts a template into a persistable ProcessStep. The deadline /
// due date is applied separately from the Master Deadline configuration, so a
// step only gets a due date when KADEP has enabled it.
func (t StepTemplate) toModel(projectID uint, sequence int) model.ProcessStep {
	return model.ProcessStep{
		ProjectID:          projectID,
		Code:               t.Code,
		Category:           t.Category,
		Name:               t.Name,
		Sequence:           sequence,
		Status:             model.StatusPending,
		RequiresPrice:      t.RequiresPrice,
		RequiresSPK:        t.RequiresSPK,
		PriceLabel:         t.PriceLabel,
		NotifyDepartments:  t.NotifyDepartments,
		ConfidentialOutput: t.ConfidentialOutput,
		SLADays:            t.SLADays,
	}
}
