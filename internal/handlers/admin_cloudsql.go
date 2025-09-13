package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"google.golang.org/api/option"
	"google.golang.org/api/sqladmin/v1"
)

type cloudSQLImpl struct {
	project			string
	svc					*sqladmin.Service
}

func NewCloudSQL(cfgProject struct{ ProjectID string }) CloudSQL {
	if cfgProject.ProjectID == "" { return nil }
	svc, _ := sqladmin.NewService(context.Background(), option.WithScopes(sqladmin.SqlserviceAdminScope))
	return &cloudSQLImpl{project: cfgProject.ProjectID, svc: svc}
}

func (c *cloudSQLImpl) ListInstances(ctx context.Context) (any, error) {
	return c.svc.Instances.List(c.project).Context(ctx).Do()
}

func (c *cloudSQLImpl) ListDatabases(ctx context.Context, instance, project string) (any, error) {
	p := project; if p == "" { p = c.project }
	return c.svc.Databases.List(p, instance).Context(ctx).Do()
}

func (c *cloudSQLImpl) CreateDatabase(ctx context.Context, project, instance, name string) (any, error) {
	p := project; if p == "" { p = c.project }
	db := &sqladmin.Database{Name: name}
	return c.svc.Databases.Insert(p, instance, db).Context(ctx).Do()
}

func (h *Handlers) listInstances(w http.ResponseWriter, r *http.Request) {
	if h.cloudSQL == nil { http.Error(w, "disabled", 404); return }
	out, err := h.cloudSQL.ListInstances(r.Context()); if err != nil { http.Error(w, err.Error(), 500); return }
	json.NewEncoder(w).Encode(out)
}

func (h *Handlers) listDatabases(w http.ResponseWriter, r *http.Request) {
	if h.cloudSQL == nil { http.Error(w, "disabled", 404); return }
	inst := r.URL.Query().Get("instance"); if inst == "" { http.Error(w, "instance required", 400); return }
	out, err := h.cloudSQL.ListDatabases(r.Context(), inst, ""); if err != nil { http.Error(w, err.Error(), 500); return }
	json.NewEncoder(w).Encode(out)
}

func (h *Handlers) createDatabase(w http.ResponseWriter, r *http.Request) {
	if h.cloudSQL == nil { http.Error(w, "disabled", 404); return }
	var req struct { Instance string `json:"instance"`; Database string `json:"database"`; Project string `json:"project,omitempty"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { http.Error(w, err.Error(), 400); return }
	if req.Instance == "" || req.Database == "" { http.Error(w, "instance and database required", 400); return }
	out, err := h.cloudSQL.CreateDatabase(r.Context(), req.Project, req.Instance, req.Database)
	if err != nil { http.Error(w, err.Error(), 500); return }
	json.NewEncoder(w).Encode(out)
}
