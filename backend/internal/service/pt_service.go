package service

import (
	"fmt"
	"mime/multipart"
	"os"

	"legalpermit/internal/model"
	"legalpermit/internal/repository"
	"legalpermit/internal/storage"
)

type PTService struct {
	pts       *repository.PTRepository
	uploadDir string
}

func NewPTService(pts *repository.PTRepository, uploadDir string) *PTService {
	return &PTService{pts: pts, uploadDir: uploadDir}
}

type CreatePTInput struct {
	Name  string `json:"name" binding:"required"`
	NPWP  string `json:"npwp"`
	NIB   string `json:"nib"`
	Notes string `json:"notes"`
}

func (s *PTService) Create(in CreatePTInput, createdBy uint) (*model.PTMaster, error) {
	pt := &model.PTMaster{Name: in.Name, NPWP: in.NPWP, NIB: in.NIB, Notes: in.Notes, CreatedBy: createdBy}
	if err := s.pts.Create(pt); err != nil {
		return nil, err
	}
	return pt, nil
}

func (s *PTService) List() ([]model.PTMaster, error) { return s.pts.List() }

func (s *PTService) Get(id uint) (*model.PTMaster, error) { return s.pts.FindByID(id) }

func (s *PTService) GetDocument(id uint) (*model.PTDocument, error) { return s.pts.FindDocument(id) }

// UploadDocument stores a PT master document.
func (s *PTService) UploadDocument(ptID uint, docType string, fh *multipart.FileHeader, uploadedBy uint, save storage.SaveFunc) (*model.PTDocument, error) {
	if _, err := s.pts.FindByID(ptID); err != nil {
		return nil, err
	}
	saved, err := storage.Save(s.uploadDir, fmt.Sprintf("pt_%d", ptID), docType, fh, save)
	if err != nil {
		return nil, err
	}
	doc := &model.PTDocument{
		PTMasterID:   ptID,
		DocType:      docType,
		OriginalName: fh.Filename,
		StoredName:   saved.StoredName,
		Path:         saved.Path,
		MimeType:     fh.Header.Get("Content-Type"),
		SizeBytes:    fh.Size,
		UploadedBy:   uploadedBy,
	}
	if err := s.pts.AddDocument(doc); err != nil {
		_ = os.Remove(saved.Path)
		return nil, err
	}
	return doc, nil
}
