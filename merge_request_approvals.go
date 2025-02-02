package gitlab

import (
	"fmt"
	"time"
)

// MergeRequestApprovalsService handles communication with the merge request
// approvals related methods of the GitLab API. This includes reading/updating
// approval settings and approve/unapproving merge requests
//
// GitLab API docs: https://docs.gitlab.com/ee/api/merge_request_approvals.html
type MergeRequestApprovalsService struct {
	client *Client
}

// MergeRequestApprovals represents GitLab merge request approvals.
//
// GitLab API docs:
// https://docs.gitlab.com/ee/api/merge_request_approvals.html#merge-request-level-mr-approvals
type MergeRequestApprovals struct {
	ID                   int                          `json:"id"`
	ProjectID            int                          `json:"project_id"`
	Title                string                       `json:"title"`
	Description          string                       `json:"description"`
	State                string                       `json:"state"`
	CreatedAt            *time.Time                   `json:"created_at"`
	UpdatedAt            *time.Time                   `json:"updated_at"`
	MergeStatus          string                       `json:"merge_status"`
	ApprovalsBeforeMerge int                          `json:"approvals_before_merge"`
	ApprovalsRequired    int                          `json:"approvals_required"`
	ApprovalsLeft        int                          `json:"approvals_left"`
	ApprovedBy           []*MergeRequestApproverUser  `json:"approved_by"`
	Approvers            []*MergeRequestApproverUser  `json:"approvers"`
	ApproverGroups       []*MergeRequestApproverGroup `json:"approver_groups"`
	SuggestedApprovers   []*BasicUser                 `json:"suggested_approvers"`
}

func (m MergeRequestApprovals) String() string {
	return Stringify(m)
}

// MergeRequestApproverGroup  represents GitLab project level merge request approver group.
//
// GitLab API docs:
// https://docs.gitlab.com/ee/api/merge_request_approvals.html#project-level-mr-approvals
type MergeRequestApproverGroup struct {
	Group struct {
		ID                   int    `json:"id"`
		Name                 string `json:"name"`
		Path                 string `json:"path"`
		Description          string `json:"description"`
		Visibility           string `json:"visibility"`
		AvatarURL            string `json:"avatar_url"`
		WebURL               string `json:"web_url"`
		FullName             string `json:"full_name"`
		FullPath             string `json:"full_path"`
		LFSEnabled           bool   `json:"lfs_enabled"`
		RequestAccessEnabled bool   `json:"request_access_enabled"`
	}
}

// MergeRequestApproverUser  represents GitLab project level merge request approver user.
//
// GitLab API docs:
// https://docs.gitlab.com/ee/api/merge_request_approvals.html#project-level-mr-approvals
type MergeRequestApproverUser struct {
	User *BasicUser
}

// ApproveMergeRequestOptions represents the available ApproveMergeRequest() options.
//
// GitLab API docs:
// https://docs.gitlab.com/ee/api/merge_request_approvals.html#approve-merge-request
type ApproveMergeRequestOptions struct {
	SHA *string `url:"sha,omitempty" json:"sha,omitempty"`
}

// ApproveMergeRequest approves a merge request on GitLab. If a non-empty sha
// is provided then it must match the sha at the HEAD of the MR.
//
// GitLab API docs:
// https://docs.gitlab.com/ee/api/merge_request_approvals.html#approve-merge-request
func (s *MergeRequestApprovalsService) ApproveMergeRequest(pid interface{}, mr int, opt *ApproveMergeRequestOptions, options ...OptionFunc) (*MergeRequestApprovals, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/merge_requests/%d/approve", pathEscape(project), mr)

	req, err := s.client.NewRequest("POST", u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	m := new(MergeRequestApprovals)
	resp, err := s.client.Do(req, m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, err
}

// UnapproveMergeRequest unapproves a previously approved merge request on GitLab.
//
// GitLab API docs:
// https://docs.gitlab.com/ee/api/merge_request_approvals.html#unapprove-merge-request
func (s *MergeRequestApprovalsService) UnapproveMergeRequest(pid interface{}, mr int, options ...OptionFunc) (*Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("projects/%s/merge_requests/%d/unapprove", pathEscape(project), mr)

	req, err := s.client.NewRequest("POST", u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// ChangeMergeRequestApprovalConfigurationOptions represents the available
// ChangeMergeRequestApprovalConfiguration() options.
//
// GitLab API docs:
// https://docs.gitlab.com/ee/api/merge_request_approvals.html#change-approval-configuration
type ChangeMergeRequestApprovalConfigurationOptions struct {
	ApprovalsRequired *int `url:"approvals_required,omitempty" json:"approvals_required,omitempty"`
}

// ChangeApprovalConfiguration updates the approval configuration of a merge request.
//
// GitLab API docs:
// https://docs.gitlab.com/ee/api/merge_request_approvals.html#change-approval-configuration
func (s *MergeRequestApprovalsService) ChangeApprovalConfiguration(pid interface{}, mergeRequestIID int, opt *ChangeMergeRequestApprovalConfigurationOptions, options ...OptionFunc) (*MergeRequest, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/merge_requests/%d/approvals", pathEscape(project), mergeRequestIID)

	req, err := s.client.NewRequest("POST", u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	m := new(MergeRequest)
	resp, err := s.client.Do(req, m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, err
}

// ChangeMergeRequestAllowedApproversOptions represents the available
// ChangeMergeRequestAllowedApprovers() options.
//
// GitLab API docs:
// https://docs.gitlab.com/ee/api/merge_request_approvals.html#change-allowed-approvers-for-merge-request
type ChangeMergeRequestAllowedApproversOptions struct {
	ApproverIDs      []int `url:"approver_ids,omitempty" json:"approver_ids,omitempty"`
	ApproverGroupIDs []int `url:"approver_group_ids,omitempty" json:"approver_group_ids,omitempty"`
}

// ChangeAllowedApprovers updates the approvers for a merge request.
//
// GitLab API docs:
// https://docs.gitlab.com/ee/api/merge_request_approvals.html#change-allowed-approvers-for-merge-request
func (s *MergeRequestApprovalsService) ChangeAllowedApprovers(pid interface{}, mergeRequestIID int, opt *ChangeMergeRequestAllowedApproversOptions, options ...OptionFunc) (*MergeRequest, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/merge_requests/%d/approvers", pathEscape(project), mergeRequestIID)

	req, err := s.client.NewRequest("PUT", u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	m := new(MergeRequest)
	resp, err := s.client.Do(req, m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, err
}

// ProjectMergeRequestApprovalOptions represents the available
// ChangeProjectMergeRequestApprovalConfiguration() options
type ProjectMergeRequestApprovalSettings struct {
  ApprovalsBeforeMerge                      *int  `json:"approvals_before_merge,omitempty"`
  DisableOverridingApproversPerMergeRequest *bool `json:"disable_overriding_approvers_per_merge_request,omitempty"`
  MergeRequestsAuthorApproval               *bool `json:"merge_requests_author_approval,omitempty"`
  MergeRequestsDisableCommittersApproval    *bool `json:"merge_requests_disable_committers_approval,omitempty"`
  ResetApprovalsOnPush                      *bool `json:"reset_approvals_on_push,omitempty"`
}

// ChangeProjectMergeRequestApprovalConfoguration updates the approval configuration of a project.
//
// GitLab API docs:
// https://docs.gitlab.com/ee/api/merge_request_approvals.html#change-configuration
func (s *MergeRequestApprovalsService) ChangeProjectMergeRequestApprovalConfiguration(pid interface{}, opt *ProjectMergeRequestApprovalSettings, options ...OptionFunc) (*MergeRequestApprovals, *Response, error) {
  project, err := parseID(pid)
  if err != nil {
    return nil, nil, err
  }
  u := fmt.Sprintf("projects/%s/approvals", pathEscape(project))

  req, err := s.client.NewRequest("POST", u, opt, options)
  if err != nil {
    return nil, nil, err
  }

  m := new(MergeRequestApprovals)
  resp, err := s.client.Do(req, m)
  if err != nil {
    return nil, resp, err
  }

  return m, resp, err
}
