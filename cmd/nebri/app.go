package main 

import ( 
    "encoding/json" 
    "os"
    "log" 
    "net/http" 
	"io/ioutil"
    "github.com/codegangsta/negroni" 
    "github.com/gorilla/mux" 
) 

type App struct { 
    Router           *mux.Router 
}

func (a *App) Initialize() {  
    a.Router = mux.NewRouter() 
    a.initializeRoutes() 
} 

func (a *App) initializeRoutes() { 
    a.Router.HandleFunc("/smartestBreeds", a.getTop5SmartestCats).Methods("GET")
//    a.Router.HandleFunc("/countryBreeds/{id:[a-z]+}", a.getcountryBreeds).Methods("GET")
} 

func (a *App) getTop5SmartestCats(w http.ResponseWriter, r *http.Request) { 
    breeds, err := httpGetBreeds() 
    if err != nil { 
        respondWithError(w, http.StatusInternalServerError, err.Error()) 
        return 
    } 
    respondWithJSON(w, http.StatusOK, breeds) 
}

func (a *App) Run(addr string) { 
    n := negroni.Classic() 
    n.UseHandler(a.Router) 
    log.Fatal(http.ListenAndServe(addr, n)) 
} 

/* Helper Funcitions */

func respondWithError(w http.ResponseWriter, code int, message string) { 
    respondWithJSON(w, code, map[string]string{"error": message}) 
} 

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) { 
    response, _ := json.Marshal(payload) 

    w.Header().Set("Content-Type", "application/json") 
    w.WriteHeader(code) 
    w.Write(response) 
}

func httpGetBreeds() (string, error) {
    resp, err := http.Get(os.Getenv("HAWK_URI_BREEDS"))
    if err != nil {
        return "", err
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
        return "", err
		log.Fatalln(err)
    }
    return string(body), nil
}