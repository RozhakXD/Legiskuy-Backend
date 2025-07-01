package voter

import (
	"errors"
)

type Service interface {
	CreateVoter(input *CreateVoterInput) (*Voter, error)
	GetAllVoters(name string) ([]Voter, error)
	GetVoterByID(id int) (*Voter, error)
	UpdateVoter(id int, input *UpdateVoterInput) (*Voter, error)
	DeleteVoter(id int) error
}

type service struct {
	repository Repository
}

func NewService(repo Repository) Service {
	return &service{
		repository: repo,
	}
}

type CreateVoterInput struct {
	Name string `json:"name"`
}

type UpdateVoterInput struct {
	Name string `json:"name"`
}

func (s *service) CreateVoter(input *CreateVoterInput) (*Voter, error) {
	if input.Name == "" {
		return nil, errors.New("name is required")
	}

	voter := &Voter{
		Name: input.Name,
	}
	id, err := s.repository.Create(voter)
	if err != nil {
		return nil, err
	}
	voter.ID = int(id)
	return voter, nil
}

func (s *service) GetAllVoters(name string) ([]Voter, error) {
	return s.repository.FindAll(name)
}

func (s *service) GetVoterByID(id int) (*Voter, error) {
	return s.repository.FindByID(id)
}

func (s *service) UpdateVoter(id int, input *UpdateVoterInput) (*Voter, error) {
	if input.Name == "" {
		return nil, errors.New("name is required")
	}

	voterToUpdate := &Voter{
		Name: input.Name,
	}
	err := s.repository.Update(id, voterToUpdate)
	if err != nil {
		return nil, err
	}
	voterToUpdate.ID = id
	return voterToUpdate, nil
}

func (s *service) DeleteVoter(id int) error {
	return s.repository.Delete(id)
}
