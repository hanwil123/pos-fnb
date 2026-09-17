package handler

import (
	"net/http"

	"github.com/yourorg/pos-fnb-backend/internal/utils"
)

// HealthCheck — dipakai load balancer/orchestrator untuk cek service hidup.
// GET /health
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	utils.Success(w, http.StatusOK, map[string]string{"status": "ok"})
}
