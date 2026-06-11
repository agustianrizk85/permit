package service

import (
	"fmt"
	"mime/multipart"
	"os"

	"legalpermit/internal/model"
	"legalpermit/internal/repository"
	"legalpermit/internal/storage"
)

type DocumentService struct {
	docs      *repository.DocumentRepository
	uploadDir string
}

func NewDocumentService(docs *repository.DocumentRepository, uploadDir string) *DocumentService {
	return &DocumentService{docs: docs, uploadDir: uploadDir}
}

// UploadParams describes a single file upload.
type UploadParams struct {
	ProjectID     uint
	ProcessStepID *uint
	DocType       string
	Confidential  bool
	UploadedBy    uint
}

// Upload stores the file on disk and records its metadata.
func (s *DocumentService) Upload(fh *multipart.FileHeader, p UploadParams, save storage.SaveFunc) (*model.Document, error) {
	subdir := fmt.Sprintf("project_%d", p.ProjectID)
	saved, err := storage.Save(s.uploadDir, subdir, p.DocType, fh, save)
	if err != nil {
		return nil, err
	}

	doc := &model.Document{
		ProjectID:     p.ProjectID,
		ProcessStepID: p.ProcessStepID,
		DocType:       p.DocType,
		OriginalName:  fh.Filename,
		StoredName:    saved.StoredName,
		Path:          saved.Path,
		MimeType:      fh.Header.Get("Content-Type"),
		SizeBytes:     fh.Size,
		Confidential:  p.Confidential,
		UploadedBy:    p.UploadedBy,
	}
	if err := s.docs.Create(doc); err != nil {
		_ = os.Remove(saved.Path)
		return nil, err
	}
	return doc, nil
}

func (s *DocumentService) Get(id uint) (*model.Document, error) {
	return s.docs.FindByID(id)
}

// Search finds documents by free text across type and original name, optionally
// scoped to a project. Powers the dashboard "Search All Dokumen" feature.
func (s *DocumentService) Search(query string, projectID *uint) ([]model.Document, error) {
	return s.docs.Search(query, projectID)
}
