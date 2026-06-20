// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package thirdPartyLink

import (
	"errors"
	"fmt"
	"io"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
)

var (
	errGetThirdPartyLink    = errors.New("failed to get third party link")
	errInsertThirdPartyLink = errors.New("failed to insert third party link")
	errUpdateThirdPartyLink = errors.New("failed to update third party link")
	errDeleteThirdPartyLink = errors.New("failed to delete third party link")
)

type ThirdPartyLink struct {
	common.Fields
	LinkingToken      string `yaml:"linking_token" json:"linking_token,omitempty"`
	Type              string `yaml:"type" json:"type,omitempty"`
	LinkStatus        string `yaml:"link_status" json:"link_status,omitempty"`
	ExternalChannelId string `yaml:"external_channel_id" json:"external_channel_id,omitempty"`
}

type IThirdPartyLink[T any] interface {
	Get() ([]*T, error)
	List(io.Writer) error
	Insert(io.Writer) error
	Update(io.Writer) error
	Delete(io.Writer) error
}

type Option func(*ThirdPartyLink)

func NewThirdPartyLink(opts ...Option) IThirdPartyLink[youtube.ThirdPartyLink] {
	tpl := &ThirdPartyLink{Fields: common.Fields{}}
	for _, opt := range opts {
		opt(tpl)
	}
	return tpl
}

func (tpl *ThirdPartyLink) Get() ([]*youtube.ThirdPartyLink, error) {
	if err := tpl.EnsureService(); err != nil {
		return nil, err
	}
	call := tpl.Service.ThirdPartyLinks.List(tpl.Parts)
	if tpl.LinkingToken != "" {
		call = call.LinkingToken(tpl.LinkingToken)
	}
	if tpl.Type != "" {
		call = call.Type(tpl.Type)
	}
	if tpl.ExternalChannelId != "" {
		call = call.ExternalChannelId(tpl.ExternalChannelId)
	}

	res, err := call.Do()
	if err != nil {
		return nil, errors.Join(errGetThirdPartyLink, err)
	}
	return res.Items, nil
}

func (tpl *ThirdPartyLink) List(writer io.Writer) error {
	links, err := tpl.Get()
	if err != nil {
		return err
	}

	common.PrintList(
		tpl.Output, links, writer,
		table.Row{"Linking Token", "Type", "Link Status"},
		func(link *youtube.ThirdPartyLink) table.Row {
			var linkStatus, linkType string
			if link.Status != nil {
				linkStatus = link.Status.LinkStatus
			}
			if link.Snippet != nil {
				linkType = link.Snippet.Type
			}
			return table.Row{link.LinkingToken, linkType, linkStatus}
		},
	)
	return nil
}

func (tpl *ThirdPartyLink) Insert(writer io.Writer) error {
	if err := tpl.EnsureService(); err != nil {
		return err
	}
	link := &youtube.ThirdPartyLink{
		LinkingToken: tpl.LinkingToken,
		Snippet: &youtube.ThirdPartyLinkSnippet{
			Type: tpl.Type,
		},
		Status: &youtube.ThirdPartyLinkStatus{
			LinkStatus: tpl.LinkStatus,
		},
	}

	call := tpl.Service.ThirdPartyLinks.Insert(tpl.Parts, link)
	if tpl.ExternalChannelId != "" {
		call = call.ExternalChannelId(tpl.ExternalChannelId)
	}

	res, err := call.Do()
	if err != nil {
		return errors.Join(errInsertThirdPartyLink, err)
	}

	common.PrintResult(
		tpl.Output, res, writer, "Third party link inserted: %s\n", res.LinkingToken,
	)
	return nil
}

func (tpl *ThirdPartyLink) Update(writer io.Writer) error {
	if err := tpl.EnsureService(); err != nil {
		return err
	}

	newType := tpl.Type
	tpl.Type = ""
	links, err := tpl.Get()
	tpl.Type = newType
	if err != nil {
		return errors.Join(errUpdateThirdPartyLink, err)
	}
	if len(links) == 0 {
		return errGetThirdPartyLink
	}

	link := links[0]
	if tpl.LinkStatus != "" {
		if link.Status == nil {
			link.Status = &youtube.ThirdPartyLinkStatus{}
		}
		link.Status.LinkStatus = tpl.LinkStatus
	}
	if tpl.Type != "" {
		if link.Snippet == nil {
			link.Snippet = &youtube.ThirdPartyLinkSnippet{}
		}
		link.Snippet.Type = tpl.Type
	}

	call := tpl.Service.ThirdPartyLinks.Update(tpl.Parts, link)
	if tpl.ExternalChannelId != "" {
		call = call.ExternalChannelId(tpl.ExternalChannelId)
	}

	res, err := call.Do()
	if err != nil {
		return errors.Join(errUpdateThirdPartyLink, err)
	}

	common.PrintResult(
		tpl.Output, res, writer, "Third party link updated: %s\n", res.LinkingToken,
	)
	return nil
}

func (tpl *ThirdPartyLink) Delete(writer io.Writer) error {
	if err := tpl.EnsureService(); err != nil {
		return err
	}
	call := tpl.Service.ThirdPartyLinks.Delete(tpl.LinkingToken, tpl.Type)
	if tpl.ExternalChannelId != "" {
		call = call.ExternalChannelId(tpl.ExternalChannelId)
	}
	if len(tpl.Parts) > 0 {
		call = call.Part(tpl.Parts...)
	}

	err := call.Do()
	if err != nil {
		return errors.Join(errDeleteThirdPartyLink, err)
	}

	_, _ = fmt.Fprintf(writer, "Third party link deleted: %s\n", tpl.LinkingToken)
	return nil
}

func WithLinkingToken(linkingToken string) Option {
	return func(tpl *ThirdPartyLink) {
		tpl.LinkingToken = linkingToken
	}
}

func WithType(type_ string) Option {
	return func(tpl *ThirdPartyLink) {
		tpl.Type = type_
	}
}

func WithLinkStatus(linkStatus string) Option {
	return func(tpl *ThirdPartyLink) {
		tpl.LinkStatus = linkStatus
	}
}

func WithExternalChannelId(externalChannelId string) Option {
	return func(tpl *ThirdPartyLink) {
		tpl.ExternalChannelId = externalChannelId
	}
}

var (
	WithParts   = common.WithParts[*ThirdPartyLink]
	WithOutput  = common.WithOutput[*ThirdPartyLink]
	WithService = common.WithService[*ThirdPartyLink]
)
