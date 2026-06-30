package sickleave

import (
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/pluginapi"
	"github.com/pkg/errors"
)

const (
	activeKeyPrefix = "sickleave:active:"
	recordKeyPrefix = "sickleave:record:"
)

type Store interface {
	GetActive(userID string) (*Record, error)
	SetActive(userID string, record *Record) error
	ClearActive(userID string) error
	SaveRecord(record *Record) error
	GetRecord(recordID string) (*Record, error)
}

type KVStore struct {
	client *pluginapi.Client
}

func NewStore(client *pluginapi.Client) Store {
	return &KVStore{client: client}
}

func (s *KVStore) GetActive(userID string) (*Record, error) {
	var record Record
	err := s.client.KV.Get(activeKey(userID), &record)
	if err != nil {
		var appErr *model.AppError
		if errors.As(err, &appErr) && appErr.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to get active sick leave record")
	}
	if record.ID == "" {
		return nil, nil
	}
	return &record, nil
}

func (s *KVStore) SetActive(userID string, record *Record) error {
	if _, err := s.client.KV.Set(activeKey(userID), record); err != nil {
		return errors.Wrap(err, "failed to set active sick leave record")
	}
	return nil
}

func (s *KVStore) ClearActive(userID string) error {
	if err := s.client.KV.Delete(activeKey(userID)); err != nil {
		return errors.Wrap(err, "failed to clear active sick leave record")
	}
	return nil
}

func (s *KVStore) SaveRecord(record *Record) error {
	if _, err := s.client.KV.Set(recordKey(record.ID), record); err != nil {
		return errors.Wrap(err, "failed to save sick leave record")
	}
	return nil
}

func (s *KVStore) GetRecord(recordID string) (*Record, error) {
	var record Record
	err := s.client.KV.Get(recordKey(recordID), &record)
	if err != nil {
		var appErr *model.AppError
		if errors.As(err, &appErr) && appErr.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to get sick leave record")
	}
	if record.ID == "" {
		return nil, nil
	}
	return &record, nil
}

func activeKey(userID string) string {
	return fmt.Sprintf("%s%s", activeKeyPrefix, userID)
}

func recordKey(recordID string) string {
	return fmt.Sprintf("%s%s", recordKeyPrefix, recordID)
}
