package service

import (
	"legalpermit/internal/model"
	"legalpermit/internal/repository"
)

type VendorService struct {
	vendors *repository.VendorRepository
}

func NewVendorService(vendors *repository.VendorRepository) *VendorService {
	return &VendorService{vendors: vendors}
}

type VendorInput struct {
	Category      string `json:"category"`
	Name          string `json:"name" binding:"required"`
	Address       string `json:"address"`
	KTPNumber     string `json:"ktp_number"`
	AccountNumber string `json:"account_number"`
	BankName      string `json:"bank_name"`
	AccountHolder string `json:"account_holder"`
	Notes         string `json:"notes"`
}

func (s *VendorService) Create(in VendorInput, createdBy uint) (*model.Vendor, error) {
	if in.Category == "" {
		in.Category = "Legal Permit"
	}
	v := &model.Vendor{
		Category:      in.Category,
		Name:          in.Name,
		Address:       in.Address,
		KTPNumber:     in.KTPNumber,
		AccountNumber: in.AccountNumber,
		BankName:      in.BankName,
		AccountHolder: in.AccountHolder,
		Notes:         in.Notes,
		CreatedBy:     createdBy,
	}
	if err := s.vendors.Create(v); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *VendorService) Update(id uint, in VendorInput) (*model.Vendor, error) {
	v, err := s.vendors.FindByID(id)
	if err != nil {
		return nil, err
	}
	if in.Category != "" {
		v.Category = in.Category
	}
	v.Name = in.Name
	v.Address = in.Address
	v.KTPNumber = in.KTPNumber
	v.AccountNumber = in.AccountNumber
	v.BankName = in.BankName
	v.AccountHolder = in.AccountHolder
	v.Notes = in.Notes
	if err := s.vendors.Update(v); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *VendorService) List() ([]model.Vendor, error) { return s.vendors.List() }

func (s *VendorService) Get(id uint) (*model.Vendor, error) { return s.vendors.FindByID(id) }
