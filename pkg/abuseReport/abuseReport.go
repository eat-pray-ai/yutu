// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package abuseReport

import (
	"errors"
	"io"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"google.golang.org/api/youtube/v3"
)

var errInsertAbuseReport = errors.New("failed to insert abuse report")

type AbuseReport struct {
	common.Fields
	AbuseTypes      []string `yaml:"abuse_types" json:"abuse_types,omitempty"`
	Description     string   `yaml:"description" json:"description,omitempty"`
	SubjectId       string   `yaml:"subject_id" json:"subject_id,omitempty"`
	SubjectTypeId   string   `yaml:"subject_type_id" json:"subject_type_id,omitempty"`
	SubjectUrl      string   `yaml:"subject_url" json:"subject_url,omitempty"`
	RelatedEntityId string   `yaml:"related_entity_id" json:"related_entity_id,omitempty"`
}

type IAbuseReport interface {
	Insert(io.Writer) error
}

type Option func(*AbuseReport)

func NewAbuseReport(opts ...Option) IAbuseReport {
	r := &AbuseReport{Fields: common.Fields{}}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *AbuseReport) Insert(writer io.Writer) error {
	if err := r.EnsureService(); err != nil {
		return err
	}

	var abuseTypes []*youtube.AbuseType
	for _, id := range r.AbuseTypes {
		abuseTypes = append(abuseTypes, &youtube.AbuseType{Id: id})
	}

	report := &youtube.AbuseReport{
		AbuseTypes:  abuseTypes,
		Description: r.Description,
		Subject: &youtube.Entity{
			Id:     r.SubjectId,
			TypeId: r.SubjectTypeId,
			Url:    r.SubjectUrl,
		},
	}

	if r.RelatedEntityId != "" {
		report.RelatedEntities = []*youtube.RelatedEntity{
			{Entity: &youtube.Entity{Id: r.RelatedEntityId}},
		}
	}

	call := r.Service.AbuseReports.Insert(r.Parts, report)
	res, err := call.Do()
	if err != nil {
		return errors.Join(errInsertAbuseReport, err)
	}

	common.PrintResult(
		r.Output, res, writer, "Abuse report submitted\n",
	)
	return nil
}

func WithAbuseTypes(abuseTypes []string) Option {
	return func(r *AbuseReport) {
		r.AbuseTypes = abuseTypes
	}
}

func WithDescription(description string) Option {
	return func(r *AbuseReport) {
		r.Description = description
	}
}

func WithSubjectId(subjectId string) Option {
	return func(r *AbuseReport) {
		r.SubjectId = subjectId
	}
}

func WithSubjectTypeId(subjectTypeId string) Option {
	return func(r *AbuseReport) {
		r.SubjectTypeId = subjectTypeId
	}
}

func WithSubjectUrl(subjectUrl string) Option {
	return func(r *AbuseReport) {
		r.SubjectUrl = subjectUrl
	}
}

func WithRelatedEntityId(relatedEntityId string) Option {
	return func(r *AbuseReport) {
		r.RelatedEntityId = relatedEntityId
	}
}

var (
	WithParts   = common.WithParts[*AbuseReport]
	WithOutput  = common.WithOutput[*AbuseReport]
	WithService = common.WithService[*AbuseReport]
)
