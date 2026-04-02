package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Metric struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Value int `json:"value"`
	Unit string `json:"unit"`
	Source string `json:"source"`
	Tags string `json:"tags"`
	Status string `json:"status"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"metrics.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS metrics(id TEXT PRIMARY KEY,name TEXT NOT NULL,type TEXT DEFAULT 'gauge',value INTEGER DEFAULT 0,unit TEXT DEFAULT '',source TEXT DEFAULT '',tags TEXT DEFAULT '',status TEXT DEFAULT 'active',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Metric)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO metrics(id,name,type,value,unit,source,tags,status,created_at)VALUES(?,?,?,?,?,?,?,?,?)`,e.ID,e.Name,e.Type,e.Value,e.Unit,e.Source,e.Tags,e.Status,e.CreatedAt);return err}
func(d *DB)Get(id string)*Metric{var e Metric;if d.db.QueryRow(`SELECT id,name,type,value,unit,source,tags,status,created_at FROM metrics WHERE id=?`,id).Scan(&e.ID,&e.Name,&e.Type,&e.Value,&e.Unit,&e.Source,&e.Tags,&e.Status,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Metric{rows,_:=d.db.Query(`SELECT id,name,type,value,unit,source,tags,status,created_at FROM metrics ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Metric;for rows.Next(){var e Metric;rows.Scan(&e.ID,&e.Name,&e.Type,&e.Value,&e.Unit,&e.Source,&e.Tags,&e.Status,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Update(e *Metric)error{_,err:=d.db.Exec(`UPDATE metrics SET name=?,type=?,value=?,unit=?,source=?,tags=?,status=? WHERE id=?`,e.Name,e.Type,e.Value,e.Unit,e.Source,e.Tags,e.Status,e.ID);return err}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM metrics WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM metrics`).Scan(&n);return n}

func(d *DB)Search(q string, filters map[string]string)[]Metric{
    where:="1=1"
    args:=[]any{}
    if q!=""{
        where+=" AND (name LIKE ?)"
        args=append(args,"%"+q+"%");
    }
    if v,ok:=filters["type"];ok&&v!=""{where+=" AND type=?";args=append(args,v)}
    if v,ok:=filters["source"];ok&&v!=""{where+=" AND source=?";args=append(args,v)}
    if v,ok:=filters["status"];ok&&v!=""{where+=" AND status=?";args=append(args,v)}
    rows,_:=d.db.Query(`SELECT id,name,type,value,unit,source,tags,status,created_at FROM metrics WHERE `+where+` ORDER BY created_at DESC`,args...)
    if rows==nil{return nil};defer rows.Close()
    var o []Metric;for rows.Next(){var e Metric;rows.Scan(&e.ID,&e.Name,&e.Type,&e.Value,&e.Unit,&e.Source,&e.Tags,&e.Status,&e.CreatedAt);o=append(o,e)};return o
}

func(d *DB)Stats()map[string]any{
    m:=map[string]any{"total":d.Count()}
    rows,_:=d.db.Query(`SELECT status,COUNT(*) FROM metrics GROUP BY status`)
    if rows!=nil{defer rows.Close();by:=map[string]int{};for rows.Next(){var s string;var c int;rows.Scan(&s,&c);by[s]=c};m["by_status"]=by}
    return m
}
