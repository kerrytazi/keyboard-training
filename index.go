package main


import (
    "log"
    "os"
    "strconv"
    "time"
    "math/rand"
    "io/ioutil"
    "net/http"
    "net/url"
    "path/filepath"
)


var CURRENT string;
var FILES = make(map[string]string);
var TEXTS = make(map[string]int);
var PORT = os.Getenv("PORT");


func init() {
    cur, err := filepath.Abs(".");
    if err != nil {
        log.Fatal(err);
    }

    if len(PORT) == 0 {
        PORT = "80";
    }

    CURRENT = cur;
    FILES["/"] = filepath.Join(CURRENT, "public",  "index.html");
    FILES["/main.css"] = filepath.Join(CURRENT, "public", "main.css");
    FILES["/main.js"] = filepath.Join(CURRENT, "public", "main.js");
    FILES["/vue.min.js"] = filepath.Join(CURRENT, "public", "vue.min.js");

    locales, _ := ioutil.ReadDir(filepath.Join(CURRENT, "texts"));
    for _, locale := range locales {
        name := locale.Name();
        t, _ := ioutil.ReadDir(filepath.Join(CURRENT, "texts", name));
        TEXTS[name] = len(t);
    }

    rand.Seed(time.Now().UTC().UnixNano());
}


func getLocale(rawQuery string) string {
    urlQueryMap, err := url.ParseQuery(rawQuery);
    if err != nil {
        return "en";
    }

    if val, ok := urlQueryMap["locale"]; ok {
        return val[0];
    }

    return "en";
}


func getRandomTextFile(locale string) string {
    textsLen := TEXTS[locale];
    ind := strconv.Itoa(rand.Intn(textsLen));
    return "/texts/" + locale + "/" + ind;
}


func handleRedirect(res http.ResponseWriter, req *http.Request) {
    link := getRandomTextFile(getLocale(req.URL.RawQuery));
    http.Redirect(res, req, link, 302);
}


func handleRawPath(res http.ResponseWriter, req *http.Request) {
    path := filepath.Join(CURRENT, req.URL.Path);
    http.ServeFile(res, req, path);
}


func handleSpecialPath(res http.ResponseWriter, req *http.Request) {
    path := req.URL.Path;
    if _, ok := FILES[path]; !ok { return; }
    http.ServeFile(res, req, FILES[path]);
}


func main() {
    http.HandleFunc("/", handleSpecialPath);
    http.HandleFunc("/texts/", handleRawPath);
    http.HandleFunc("/locales/", handleRawPath);
    http.HandleFunc("/random_text/", handleRedirect);
    log.Fatal(http.ListenAndServe(":" + PORT, nil));
}
