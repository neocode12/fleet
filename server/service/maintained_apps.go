package service

import (
	"context"
	"errors"

	"github.com/fleetdm/fleet/v4/server/fleet"
	"github.com/fleetdm/fleet/v4/server/mdm/maintainedapps"
)

type addFleetMaintainedAppRequest struct {
	TeamID            *uint    `json:"team_id"`
	AppID             uint     `json:"fleet_maintained_app_id"`
	InstallScript     string   `json:"install_script"`
	PreInstallQuery   string   `json:"pre_install_query"`
	PostInstallScript string   `json:"post_install_script"`
	SelfService       bool     `json:"self_service"`
	UninstallScript   string   `json:"uninstall_script"`
	LabelsIncludeAny  []string `json:"labels_include_any"`
	LabelsExcludeAny  []string `json:"labels_exclude_any"`
	AutomaticInstall  bool     `json:"automatic_install"`
}

type addFleetMaintainedAppResponse struct {
	SoftwareTitleID uint  `json:"software_title_id,omitempty"`
	Err             error `json:"error,omitempty"`
}

func (r addFleetMaintainedAppResponse) Error() error { return r.Err }

func addFleetMaintainedAppEndpoint(ctx context.Context, request interface{}, svc fleet.Service) (fleet.Errorer, error) {
	req := request.(*addFleetMaintainedAppRequest)
	ctx, cancel := context.WithTimeout(ctx, maintained_apps.InstallerTimeout)
	defer cancel()
	titleId, err := svc.AddFleetMaintainedApp(
		ctx,
		req.TeamID,
		req.AppID,
		req.InstallScript,
		req.PreInstallQuery,
		req.PostInstallScript,
		req.UninstallScript,
		req.SelfService,
		req.AutomaticInstall,
		req.LabelsIncludeAny,
		req.LabelsExcludeAny,
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			err = fleet.NewGatewayTimeoutError("Couldn't add. Request timeout. Please make sure your server and load balancer timeout is long enough.", err)
		}

		return &addFleetMaintainedAppResponse{Err: err}, nil
	}
	return &addFleetMaintainedAppResponse{SoftwareTitleID: titleId}, nil
}

func (svc *Service) AddFleetMaintainedApp(ctx context.Context, _ *uint, _ uint, _, _, _, _ string, _ bool, _ bool, _, _ []string) (uint, error) {
	// skipauth: No authorization check needed due to implementation returning
	// only license error.
	svc.authz.SkipAuthorization(ctx)

	return 0, fleet.ErrMissingLicense
}

type listFleetMaintainedAppsRequest struct {
	fleet.ListOptions
	TeamID *uint `query:"team_id,optional"`
}

type listFleetMaintainedAppsResponse struct {
	FleetMaintainedApps []fleet.MaintainedApp     `json:"fleet_maintained_apps"`
	Meta                *fleet.PaginationMetadata `json:"meta"`
	Err                 error                     `json:"error,omitempty"`
}

func (r listFleetMaintainedAppsResponse) Error() error { return r.Err }

func listFleetMaintainedAppsEndpoint(ctx context.Context, request any, svc fleet.Service) (fleet.Errorer, error) {
	req := request.(*listFleetMaintainedAppsRequest)

	apps, meta, err := svc.ListFleetMaintainedApps(ctx, req.TeamID, req.ListOptions)
	if err != nil {
		return listFleetMaintainedAppsResponse{Err: err}, nil
	}

	listResp := listFleetMaintainedAppsResponse{
		FleetMaintainedApps: apps,
		Meta:                meta,
	}

	return listResp, nil
}

func (svc *Service) ListFleetMaintainedApps(ctx context.Context, teamID *uint, opts fleet.ListOptions) ([]fleet.MaintainedApp, *fleet.PaginationMetadata, error) {
	// skipauth: No authorization check needed due to implementation returning
	// only license error.
	svc.authz.SkipAuthorization(ctx)

	return nil, nil, fleet.ErrMissingLicense
}

type getFleetMaintainedAppRequest struct {
	AppID  uint  `url:"app_id"`
	TeamID *uint `query:"team_id,optional"`
}

type getFleetMaintainedAppResponse struct {
	FleetMaintainedApp *fleet.MaintainedApp `json:"fleet_maintained_app"`
	Err                error                `json:"error,omitempty"`
}

func (r getFleetMaintainedAppResponse) Error() error { return r.Err }

func getFleetMaintainedApp(ctx context.Context, request any, svc fleet.Service) (fleet.Errorer, error) {
	req := request.(*getFleetMaintainedAppRequest)

	app, err := svc.GetFleetMaintainedApp(ctx, req.AppID, req.TeamID)
	if err != nil {
		return getFleetMaintainedAppResponse{Err: err}, nil
	}

	return getFleetMaintainedAppResponse{FleetMaintainedApp: app}, nil
}

func (svc *Service) GetFleetMaintainedApp(ctx context.Context, appID uint, teamID *uint) (*fleet.MaintainedApp, error) {
	// skipauth: No authorization check needed due to implementation returning
	// only license error.
	svc.authz.SkipAuthorization(ctx)

	return nil, fleet.ErrMissingLicense
}
