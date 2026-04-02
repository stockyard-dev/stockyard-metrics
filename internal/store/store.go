package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Metric struct{
	ID string `json:"id"`
	Name string `json:"name"`
	Value float64 `json:"value"`
	Unit string `json:"unit"`
	Source string `json:"source"`
	Tags string `json:"tags"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"metrics.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS metrics(id TEXT PRIMARY KEY,name TEXT NOT NULL,value REAL DEFAULT 0,unit TEXT DEFAULT '',source TEXT DEFAULT '',tags TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Metric)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO metrics(id,name,value,unit,source,tags,created_at)VALUES(?,?,?,?,?,?,?)`,e.ID,e.Name,e.Value,e.Unit,e.Source,e.Tags,e.CreatedAt);return err}
func(d *DB)Get(id string)*Metric{var e Metric;if d.db.QueryRow(`SELECT id,name,value,unit,source,tags,created_at FROM metrics WHERE id=?`,id).Scan(&e.ID,&e.Name,&e.Value,&e.Unit,&e.Source,&e.Tags,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Metric{rows,_:=d.db.Query(`SELECT id,name,value,unit,source,tags,created_at FROM metrics ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Metric;for rows.Next(){var e Metric;rows.Scan(&e.ID,&e.Name,&e.Value,&e.Unit,&e.Source,&e.Tags,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM metrics WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM metrics`).Scan(&n);return n}
