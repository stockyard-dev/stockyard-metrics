package server
import("encoding/json";"net/http";"github.com/stockyard-dev/stockyard-metrics/internal/store")
func(s *Server)handleRecord(w http.ResponseWriter,r *http.Request){var m store.Metric;json.NewDecoder(r.Body).Decode(&m);if m.Name==""{writeError(w,400,"name required");return};if m.Type==""{m.Type="gauge"};s.db.Record(&m);writeJSON(w,201,m)}
func(s *Server)handleLatest(w http.ResponseWriter,r *http.Request){list,_:=s.db.Latest();if list==nil{list=[]map[string]interface{}{}};writeJSON(w,200,list)}
func(s *Server)handleHistory(w http.ResponseWriter,r *http.Request){name:=r.URL.Query().Get("name");if name==""{writeError(w,400,"name required");return};list,_:=s.db.History(name);if list==nil{list=[]store.Metric{}};writeJSON(w,200,list)}
func(s *Server)handleOverview(w http.ResponseWriter,r *http.Request){m,_:=s.db.Stats();writeJSON(w,200,m)}
