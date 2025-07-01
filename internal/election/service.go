package election

import (
	"errors"
	"legiskuy-backend/internal/candidate"
	"legiskuy-backend/internal/voter"
	"strconv"
	"time"
)

type Service interface {
	CastVote(input *CastVoteInput) error
	SetElectionTime(input *SetTimeInput) error
	GetResults(qualifiedOnly bool) ([]candidate.Candidate, error)
	SetThreshold(input *SetThresholdInput) error
}

type service struct {
	electionRepo  Repository
	voterRepo     voter.Repository
	candidateRepo candidate.Repository
}

func NewService(electionRepo Repository, voterRepo voter.Repository, candidateRepo candidate.Repository) Service {
	return &service{
		electionRepo:  electionRepo,
		voterRepo:     voterRepo,
		candidateRepo: candidateRepo,
	}
}

type CastVoteInput struct {
	VoterID     int `json:"voter_id"`
	CandidateID int `json:"candidate_id"`
}

type SetTimeInput struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type SetThresholdInput struct {
	Threshold *int `json:"threshold"`
}

func (s *service) CastVote(input *CastVoteInput) error {
	startTimeStr, _ := s.electionRepo.GetSetting("start_time")
	endTimeStr, _ := s.electionRepo.GetSetting("end_time")

	if startTimeStr != "" && endTimeStr != "" {
		startTime, err1 := time.Parse(time.RFC3339, startTimeStr)
		endTime, err2 := time.Parse(time.RFC3339, endTimeStr)
		if err1 == nil && err2 == nil {
			now := time.Now().UTC()
			startTime = startTime.UTC()
			endTime = endTime.UTC()

			if now.Before(startTime) || now.After(endTime) {
				return errors.New("election is not currently active")
			}
		}
	}

	if input.VoterID == 0 || input.CandidateID == 0 {
		return errors.New("voter_id and candidate_id are required")
	}

	voter, err := s.voterRepo.FindByID(input.VoterID)
	if err != nil || voter == nil {
		return errors.New("voter not found")
	}

	candidate, err := s.candidateRepo.FindByID(input.CandidateID)
	if err != nil || candidate == nil {
		return errors.New("candidate not found")
	}

	if voter.HasVoted {
		return errors.New("voter has already voted")
	}

	tx, err := s.electionRepo.BeginTransaction()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err := s.voterRepo.MarkAsVoted(tx, input.VoterID); err != nil {
		return err
	}

	if err := s.candidateRepo.IncrementVoteCount(tx, input.CandidateID); err != nil {
		return err
	}

	if err := s.electionRepo.CreateVote(tx, input.VoterID, input.CandidateID); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *service) SetElectionTime(input *SetTimeInput) error {
	_, err1 := time.Parse(time.RFC3339, input.StartTime)
	_, err2 := time.Parse(time.RFC3339, input.EndTime)
	if err1 != nil || err2 != nil {
		return errors.New("invalid time format, use RFC3339 format (e.g., 2025-06-13T00:00:00Z)")
	}

	if err := s.electionRepo.SetSetting("start_time", input.StartTime); err != nil {
		return err
	}
	if err := s.electionRepo.SetSetting("end_time", input.EndTime); err != nil {
		return err
	}
	return nil
}

func (s *service) GetResults(qualifiedOnly bool) ([]candidate.Candidate, error) {
	candidates, err := s.candidateRepo.FindAll("", "")
	if err != nil {
		return nil, err
	}

	if qualifiedOnly {
		thresholdStr, _ := s.electionRepo.GetSetting("threshold")
		threshold, _ := strconv.Atoi(thresholdStr)

		qualifiedCandidates := make([]candidate.Candidate, 0)
		for _, c := range candidates {
			if c.Votes >= threshold {
				qualifiedCandidates = append(qualifiedCandidates, c)
			}
		}
		candidates = qualifiedCandidates
	}

	n := len(candidates)
	for i := 1; i < n; i++ {
		key := candidates[i]
		j := i - 1
		for j >= 0 && candidates[j].Votes < key.Votes {
			candidates[j+1] = candidates[j]
			j = j - 1
		}
		candidates[j+1] = key
	}

	return candidates, nil
}

func (s *service) SetThreshold(input *SetThresholdInput) error {
	if input.Threshold == nil {
		return errors.New("threshold is required")
	}
	if *input.Threshold < 0 {
		return errors.New("threshold must be a non-negative number")
	}
	thresholdStr := strconv.Itoa(*input.Threshold)
	return s.electionRepo.SetSetting("threshold", thresholdStr)
}
