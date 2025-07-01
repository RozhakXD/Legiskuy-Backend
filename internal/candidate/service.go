package candidate

import (
	"errors"
	"strings"
)

type Service interface {
	CreateCandidate(input *CreateCandidateInput) (*Candidate, error)
	GetAllCandidates(name, party, sortBy, order string) ([]Candidate, error)
	GetCandidateByID(id int) (*Candidate, error)
	UpdateCandidate(id int, input *UpdateCandidateInput) (*Candidate, error)
	DeleteCandidate(id int) error
}

type service struct {
	repository Repository
}

func NewService(repo Repository) Service {
	return &service{
		repository: repo,
	}
}

type CreateCandidateInput struct {
	Name  string `json:"name"`
	Party string `json:"party"`
}

type UpdateCandidateInput struct {
	Name  string `json:"name"`
	Party string `json:"party"`
}

func (s *service) CreateCandidate(input *CreateCandidateInput) (*Candidate, error) {
	if input.Name == "" {
		return nil, errors.New("name is required")
	}
	if input.Party == "" {
		return nil, errors.New("party is required")
	}

	candidate := &Candidate{
		Name:  input.Name,
		Party: input.Party,
	}

	id, err := s.repository.Create(candidate)
	if err != nil {
		return nil, err
	}

	candidate.ID = int(id)
	return candidate, nil
}

func (s *service) GetAllCandidates(name, party, sortBy, order string) ([]Candidate, error) {
	candidates, err := s.repository.FindAll(name, party)
	if err != nil {
		return nil, err
	}

	if candidates == nil {
		candidates = []Candidate{}
	}

	if sortBy != "" {
		isDescending := strings.ToLower(order) == "desc"

		switch strings.ToLower(sortBy) {
		case "name":
			selectionSort(candidates, "name", isDescending)
		case "party":
			insertionSort(candidates, "party", isDescending)
		case "votes":
			selectionSort(candidates, "votes", isDescending)
		}
	}
	return candidates, nil
}

func (s *service) GetCandidateByID(id int) (*Candidate, error) {
	return s.repository.FindByID(id)
}

func (s *service) UpdateCandidate(id int, input *UpdateCandidateInput) (*Candidate, error) {
	if input.Name == "" {
		return nil, errors.New("name is required")
	}
	if input.Party == "" {
		return nil, errors.New("party is required")
	}

	candidateToUpdate := &Candidate{
		Name:  input.Name,
		Party: input.Party,
	}

	err := s.repository.Update(id, candidateToUpdate)
	if err != nil {
		return nil, err
	}

	candidateToUpdate.ID = id
	return candidateToUpdate, nil
}

func (s *service) DeleteCandidate(id int) error {
	return s.repository.Delete(id)
}

func selectionSort(arr []Candidate, by string, descending bool) {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		extremeIdx := i
		for j := i + 1; j < n; j++ {
			isJMoreExtreme := false
			switch by {
			case "name":
				if descending {
					isJMoreExtreme = arr[j].Name > arr[extremeIdx].Name
				} else {
					isJMoreExtreme = arr[j].Name < arr[extremeIdx].Name
				}
			}
			if isJMoreExtreme {
				extremeIdx = j
			}
		}
		arr[i], arr[extremeIdx] = arr[extremeIdx], arr[i]
	}
}

func insertionSort(arr []Candidate, by string, descending bool) {
	n := len(arr)
	for i := 1; i < n; i++ {
		key := arr[i]
		j := i - 1

		for j >= 0 {
			shouldMove := false
			switch by {
			case "party":
				if descending {
					shouldMove = arr[j].Party < key.Party
				} else {
					shouldMove = arr[j].Party > key.Party
				}
			case "votes":
				if descending {
					shouldMove = arr[j].Votes < key.Votes
				} else {
					shouldMove = arr[j].Votes > key.Votes
				}
			}
			if shouldMove {
				arr[j+1] = arr[j]
				j = j - 1
			} else {
				break
			}
		}
		arr[j+1] = key
	}
}
