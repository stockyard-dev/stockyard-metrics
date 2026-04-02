package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-metrics/internal/server";"github.com/stockyard-dev/stockyard-metrics/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="9730"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./metrics-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("metrics: %v",err)};defer db.Close();srv:=server.New(db)
fmt.Printf("\n  Metrics — real-time metrics dashboard\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n\n",port,port)
log.Printf("metrics: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}
